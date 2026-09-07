package web

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
	config "github.com/tiredsosha/admin/tools/configurator"
	"github.com/tiredsosha/admin/tools/formater"
	"github.com/tiredsosha/admin/tools/logger"
	"gopkg.in/yaml.v3"
)

// -------------------- Модели данных --------------------

type Status struct {
	Zones []Zone `yaml:"zones" json:"zones"`
}

type Zone struct {
	ID         string         `yaml:"id" json:"id"`
	InnerZones []InnerZone    `yaml:"innerZones,omitempty" json:"innerZones,omitempty"`
	Status     map[string]int `yaml:"status,omitempty" json:"status,omitempty"`
}

type InnerZone struct {
	ID     string         `yaml:"id" json:"id"`
	Status map[string]int `yaml:"status" json:"status"`
}

// -------------------- Общее состояние --------------------

var (
	// Общие данные о статусах.
	statusData Status

	// Защищает statusData от одновременного чтения и изменения.
	statusMu  sync.RWMutex
	persistMu sync.Mutex

	// Не позволяет запустить один и тот же тип обновления повторно.
	pcRunning atomic.Bool
	pjRunning atomic.Bool
)

type ipRecord struct {
	Zone      string `yaml:"zone" json:"zone"`
	Equipment string `yaml:"equipment" json:"equipment"`
	IP        string `yaml:"ip" json:"ip"`
}

// -------------------- Пути к файлам --------------------

const (
	statusYAMLPath = "./configs/status.yaml"
	statusJSONPath = "./configs/status.json"
)

// -------------------- HTTP-обработчики --------------------

// statusPark возвращает последний JSON-снимок статусов из файла.
func statusPark(c *gin.Context) {
	b, err := os.ReadFile(statusJSONPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read status data"})
		return
	}
	c.Data(http.StatusOK, "application/json", b)
}

// StatusInit загружает начальные данные о статусах из YAML-файла.
func StatusInit() {
	data, err := os.ReadFile(statusYAMLPath)
	if err != nil {
		logger.Warn.Println("Error reading YAML:", err)
		return
	}

	var loaded Status
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		logger.Warn.Println("Error unmarshalling YAML:", err)
		return
	}

	if len(loaded.Zones) == 0 {
		logger.Error.Println("StatusInit blocked: loaded zero zones, keeping existing statusData")
		return
	}

	statusMu.Lock()
	statusData = loaded
	statusMu.Unlock()

	logger.Info.Printf("StatusInit: loaded %d zones", len(loaded.Zones))
}

// UpdateStatues запускает независимое обновление статусов ПК и проекторов.
func UpdateStatues() {
	if pcRunning.CompareAndSwap(false, true) {
		go func() {
			defer pcRunning.Store(false)
			updatePC()
			_ = persistSnapshots() // Ошибки уже логируются внутри.
		}()
	}

	if pjRunning.CompareAndSwap(false, true) {
		go func() {
			defer pjRunning.Store(false)
			updatePJ()
			_ = persistSnapshots()
		}()
	}

	// Мгновенную запись здесь не выполняем, чтобы не сохранять устаревшие данные.
}

// -------------------- Задачи (без блокировок во время сетевых запросов) --------------------

type pcJob struct {
	zoneIdx  int
	innerIdx int // -1 означает использование zone.Status.
	id       string
}

type pjJob struct {
	zoneIdx  int
	innerIdx int // -1 означает использование zone.Status.
	id       string
	key      string // Например: "pj_1", "pj_2".
}

var pjKeyRe = regexp.MustCompile(`^pj_\d+$`)

// updatePC запрашивает и сохраняет статусы всех настроенных ПК.
func updatePC() {
	// 1. Формируем задачи под RLock, сетевых запросов здесь нет.
	jobs := buildPCJobs()

	// 2. Выполняем сетевые запросы без блокировок.
	results := make([]int, len(jobs))
	for i, job := range jobs {
		ip := config.FindPC(job.id, "ip")
		url := formater.CustomStr("http://{ip}:3001/status", map[string]any{"ip": ip})
		results[i] = protocols.GetPC(url, 1)
	}

	// 3. Короткое применение результатов под Lock.
	statusMu.Lock()
	for i, job := range jobs {
		v := results[i]

		if job.innerIdx == -1 {
			ensureZoneStatusMap(job.zoneIdx)
			statusData.Zones[job.zoneIdx].Status["pc_1"] = v
			continue
		}

		ensureInnerStatusMap(job.zoneIdx, job.innerIdx)
		statusData.Zones[job.zoneIdx].InnerZones[job.innerIdx].Status["pc_1"] = v
	}
	statusMu.Unlock()

	logger.Debug.Println("pc statuses updated")
}

// updatePJ запрашивает и сохраняет статусы всех настроенных проекторов.
func updatePJ() {
	// 1. Формируем задачи под RLock, сетевых запросов здесь нет.
	jobs := buildPJJobs()

	// 2. Выполняем сетевые запросы без блокировок.
	results := make([]int, len(jobs))
	for i, job := range jobs {
		addr := config.FindPJ(job.id, job.key) // Функция получает ID зоны и ключ "pj_1" и т.п.
		results[i] = protocols.GetPjlink(addr)
	}

	// 3. Короткое применение результатов под Lock.
	statusMu.Lock()
	for i, job := range jobs {
		v := results[i]

		if job.innerIdx == -1 {
			ensureZoneStatusMap(job.zoneIdx)
			// Записываем только нужный ключ оборудования.
			statusData.Zones[job.zoneIdx].Status[job.key] = v
			continue
		}

		ensureInnerStatusMap(job.zoneIdx, job.innerIdx)
		statusData.Zones[job.zoneIdx].InnerZones[job.innerIdx].Status[job.key] = v
	}
	statusMu.Unlock()

	logger.Debug.Println("pj statuses updated")
}

// buildPCJobs создаёт задачи обновления ПК по текущей структуре статусов.
func buildPCJobs() []pcJob {
	statusMu.RLock()
	defer statusMu.RUnlock()

	jobs := make([]pcJob, 0, len(statusData.Zones))

	for i, z := range statusData.Zones {
		// Эти зоны не относятся к ПК.
		if z.ID == "relay" || z.ID == "light" {
			continue
		}

		// Добавляем ПК зоны, если у неё есть карта статусов.
		if z.Status != nil {
			jobs = append(jobs, pcJob{zoneIdx: i, innerIdx: -1, id: z.ID})
		}

		// Для вложенных зон также добавляем задачу обновления ПК.
		for j := range z.InnerZones {
			jobs = append(jobs, pcJob{zoneIdx: i, innerIdx: j, id: z.InnerZones[j].ID})
		}
	}

	return jobs
}

// buildPJJobs создаёт задачи обновления проекторов по текущей структуре статусов.
func buildPJJobs() []pjJob {
	statusMu.RLock()
	defer statusMu.RUnlock()

	var jobs []pjJob

	for i := range statusData.Zones {
		z := statusData.Zones[i]

		// Проекторы на уровне основной зоны.
		for key := range z.Status {
			if pjKeyRe.MatchString(key) {
				jobs = append(jobs, pjJob{zoneIdx: i, innerIdx: -1, id: z.ID, key: key})
			}
		}

		// Проекторы на уровне вложенной зоны.
		for j := range z.InnerZones {
			in := z.InnerZones[j]
			for key := range in.Status {
				if pjKeyRe.MatchString(key) {
					jobs = append(jobs, pjJob{zoneIdx: i, innerIdx: j, id: in.ID, key: key})
				}
			}
		}
	}

	// Фиксируем порядок для удобства логирования и отладки.
	sort.Slice(jobs, func(a, b int) bool {
		if jobs[a].id == jobs[b].id {
			return jobs[a].key < jobs[b].key
		}
		return jobs[a].id < jobs[b].id
	})

	return jobs
}

// ensureZoneStatusMap создаёт карту статусов зоны, если она отсутствует.
func ensureZoneStatusMap(zoneIdx int) {
	if statusData.Zones[zoneIdx].Status == nil {
		statusData.Zones[zoneIdx].Status = map[string]int{}
	}
}

// ensureInnerStatusMap создаёт карту статусов вложенной зоны, если она отсутствует.
func ensureInnerStatusMap(zoneIdx, innerIdx int) {
	if statusData.Zones[zoneIdx].InnerZones[innerIdx].Status == nil {
		statusData.Zones[zoneIdx].InnerZones[innerIdx].Status = map[string]int{}
	}
}

// -------------------- Сохранение (атомарная запись файлов) --------------------

// persistSnapshots сериализует текущие статусы и сохраняет оба снимка.
func persistSnapshots() error {
	persistMu.Lock()
	defer persistMu.Unlock()

	statusMu.RLock()
	defer statusMu.RUnlock()

	if len(statusData.Zones) == 0 {
		logger.Error.Println("persistSnapshots blocked: statusData.Zones is empty")
		return nil // Не сохраняем пустой статус и не создаём лишнюю запись в логе.
	}

	yamlBytes, err := yaml.Marshal(statusData)
	if err != nil {
		logger.Warn.Println(err)
		return err
	}

	jsonBytes, err := json.MarshalIndent(statusData, "", "  ")
	if err != nil {
		logger.Warn.Println(err)
		return err
	}

	if err := writeFileAtomic(statusYAMLPath, yamlBytes, 0644); err != nil {
		logger.Warn.Println(err)
		return err
	}

	if err := writeFileAtomic(statusJSONPath, jsonBytes, 0644); err != nil {
		logger.Warn.Println(err)
		return err
	}

	return nil
}

// writeFileAtomic полностью записывает временный файл и заменяет им целевой.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	tmpFile, err := os.CreateTemp(dir, "."+base+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

// statusIp читает YAML с IP-адресами и возвращает его в формате JSON.
func statusIp(c *gin.Context) {
	data, err := os.ReadFile("./configs/ip.yaml")
	if err != nil {
		logger.Error.Printf("read ip.yaml: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read ip.yaml"})
		return
	}

	var records []ipRecord
	if err := yaml.Unmarshal(data, &records); err != nil {
		logger.Error.Printf("unmarshal ip.yaml: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse ip.yaml"})
		return
	}

	c.JSON(http.StatusOK, records)
}
