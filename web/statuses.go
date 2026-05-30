// страрый статус не удалять

// package web
// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"regexp"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/tiredsosha/admin/protocols"
// 	config "github.com/tiredsosha/admin/tools/configurator"
// 	"github.com/tiredsosha/admin/tools/formater"
// 	"github.com/tiredsosha/admin/tools/logger"
// 	"gopkg.in/yaml.v3"
// )

// type Status struct {
// 	Zones []Zone `yaml:"zones"`
// }

// type Zone struct {
// 	ID         string         `yaml:"id"`
// 	InnerZones []InnerZone    `yaml:"innerZones,omitempty"`
// 	Status     map[string]int `yaml:"status,omitempty"`
// }

// type InnerZone struct {
// 	ID     string         `yaml:"id"`
// 	Status map[string]int `yaml:"status"`
// }

// var (
// 	// control concurrent goroutine launches
// 	mu        sync.Mutex
// 	pcRunning bool
// 	pjRunning bool

// 	// shared status data
// 	statusData Status

// 	// protects statusData (writers use Lock, readers use RLock)
// 	statusMu sync.RWMutex
// )

// // HTTP handler: returns the latest JSON snapshot
// func statusPark(c *gin.Context) {
// 	jsonBytes, err := os.ReadFile("./configs/status.json")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read status data"})
// 		return
// 	}
// 	c.Data(http.StatusOK, "application/json", jsonBytes)
// }

// // Load initial YAML into statusData (protected)
// func StatusInit() {
// 	data, err := os.ReadFile("./configs/status.yaml")
// 	if err != nil {
// 		fmt.Println("Error reading YAML:", err)
// 		return
// 	}

// 	var loaded Status
// 	if err := yaml.Unmarshal(data, &loaded); err != nil {
// 		fmt.Println("Error unmarshalling YAML:", err)
// 		return
// 	}

// 	statusMu.Lock()
// 	statusData = loaded
// 	statusMu.Unlock()
// }

// // --- Updaters ---------------------------------------------------------------

// func updatePC() {
// 	statusMu.Lock()
// 	defer statusMu.Unlock()

// 	for i, zone := range statusData.Zones {
// 		if zone.ID == "relay" || zone.ID == "light" {
// 			continue
// 		}

// 		if zone.Status != nil {
// 			statusData.Zones[i].Status["pc_1"] = protocols.GetPC(
// 				formater.CustomStr("http://{ip}:3001/status",
// 					map[string]any{"ip": config.FindPC(zone.ID, "ip")}),
// 				1,
// 			)
// 		}

// 		for j, inner := range zone.InnerZones {
// 			statusData.Zones[i].InnerZones[j].Status["pc_1"] = protocols.GetPC(
// 				formater.CustomStr("http://{ip}:3001/status",
// 					map[string]any{"ip": config.FindPC(inner.ID, "ip")}),
// 				1,
// 			)
// 		}
// 	}
// 	logger.Debug.Println("pc statuses updated")
// }

// func countPjKeys(statusMap map[string]int) int {
// 	count := 0
// 	re := regexp.MustCompile(`^pj_\d+$`)
// 	for k := range statusMap {
// 		if re.MatchString(k) {
// 			count++
// 		}
// 	}
// 	return count
// }

// func updatePJ() {
// 	statusMu.Lock()
// 	defer statusMu.Unlock()

// 	for i := range statusData.Zones {
// 		zone := &statusData.Zones[i]

// 		maxKeys := countPjKeys(zone.Status)
// 		for n := 1; n <= maxKeys; n++ {
// 			key := fmt.Sprintf("pj_%d", n)
// 			if _, exists := zone.Status[key]; exists {
// 				zone.Status[key] = protocols.GetPjlink(config.FindPJ(zone.ID, key))
// 			}
// 		}

// 		for j := range zone.InnerZones {
// 			inner := &zone.InnerZones[j]
// 			maxKeysInner := countPjKeys(inner.Status)
// 			for n := 1; n <= maxKeysInner; n++ {
// 				key := fmt.Sprintf("pj_%d", n)
// 				if _, exists := inner.Status[key]; exists {
// 					logger.Debug.Println(config.FindPJ(inner.ID, key))
// 					inner.Status[key] = protocols.GetPjlink(config.FindPJ(inner.ID, key))
// 				}
// 			}
// 		}
// 	}
// 	logger.Debug.Println("pj statuses updated")
// }

// // --- Serialization helpers --------------------------------------------------

// func convertYamlToJson(yamlData []byte) ([]byte, error) {
// 	var parsedData any
// 	if err := yaml.Unmarshal(yamlData, &parsedData); err != nil {
// 		return nil, err
// 	}
// 	jsonData, err := json.MarshalIndent(parsedData, "", "  ")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsonData, nil
// }

// // --- Public orchestrator ----------------------------------------------------

// func UpdateStatues() {
// 	mu.Lock()
// 	if !pcRunning {
// 		pcRunning = true
// 		go func() {
// 			defer func() {
// 				mu.Lock()
// 				pcRunning = false
// 				mu.Unlock()
// 			}()
// 			updatePC()
// 		}()
// 	}
// 	if !pjRunning {
// 		pjRunning = true
// 		go func() {
// 			defer func() {
// 				mu.Lock()
// 				pjRunning = false
// 				mu.Unlock()
// 			}()
// 			updatePJ()
// 		}()
// 	}
// 	mu.Unlock()

// 	update()
// }

// // Persist current statusData to YAML and JSON (read-only path)
// func update() {
// 	statusMu.RLock()
// 	updatedYAML, err := yaml.Marshal(statusData)
// 	statusMu.RUnlock()
// 	if err != nil {
// 		logger.Warn.Println(err)
// 		logger.Error.Fatal("EXITING")
// 	}

// 	if err := os.WriteFile("./configs/status.yaml", updatedYAML, 0644); err != nil {
// 		logger.Warn.Println(err)
// 		logger.Error.Fatal("EXITING")
// 	}
// 	logger.Debug.Println("yaml updated, check status.yaml")

// 	jsonBytes, err := convertYamlToJson(updatedYAML)
// 	if err != nil {
// 		logger.Warn.Println(err)
// 		logger.Error.Fatal("EXITING")
// 	}

// 	if err := os.WriteFile("./configs/status.json", jsonBytes, 0644); err != nil {
// 		logger.Warn.Println(err)
// 		logger.Error.Fatal("EXITING")
// 	}
// 	logger.Debug.Println("json updated, check status.json")
// }

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

// -------------------- Data model --------------------

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

// -------------------- Shared state --------------------

var (
	// shared status data
	statusData Status

	// protects statusData
	statusMu  sync.RWMutex
	persistMu sync.Mutex

	// prevents re-entry for same updater
	pcRunning atomic.Bool
	pjRunning atomic.Bool
)

// -------------------- Paths --------------------

const (
	statusYAMLPath = "./configs/status.yaml"
	statusJSONPath = "./configs/status.json"
)

// -------------------- HTTP handler --------------------

// Returns latest JSON snapshot from file (because a separate program reads files too)
func statusPark(c *gin.Context) {
	b, err := os.ReadFile(statusJSONPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read status data"})
		return
	}
	c.Data(http.StatusOK, "application/json", b)
}

// -------------------- Init --------------------

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

	statusMu.Lock()
	statusData = loaded
	statusMu.Unlock()
}

// -------------------- Update orchestration --------------------

// Keeps your original name for compatibility.
func UpdateStatues() {
	UpdateStatuses()
}

// запускает обновления независимо:
// - pc обновление не стартует повторно, пока предыдущее не завершилось
// - pj обновление не стартует повторно, пока предыдущее не завершилось
// - pc и pj не блокируют друг друга (нет lock во время сети)
func UpdateStatuses() {
	if pcRunning.CompareAndSwap(false, true) {
		go func() {
			defer pcRunning.Store(false)
			updatePC()
			_ = persistSnapshots() // ошибки уже логируются внутри
		}()
	}

	if pjRunning.CompareAndSwap(false, true) {
		go func() {
			defer pjRunning.Store(false)
			updatePJ()
			_ = persistSnapshots()
		}()
	}

	// Можно оставить мгновенную запись тоже, но обычно не нужно.
	// Здесь сознательно НЕ вызываем persistSnapshots(), чтобы не писать старые данные.
}

// -------------------- Jobs (no locks while network) --------------------

type pcJob struct {
	zoneIdx  int
	innerIdx int // -1 => zone.Status
	id       string
}

type pjJob struct {
	zoneIdx  int
	innerIdx int // -1 => zone.Status
	id       string
	key      string // "pj_1", "pj_2"...
}

var pjKeyRe = regexp.MustCompile(`^pj_\d+$`)

func updatePC() {
	// 1) Build jobs under RLock (no network)
	jobs := buildPCJobs()

	// 2) Network phase without locks
	results := make([]int, len(jobs))
	for i, job := range jobs {
		ip := config.FindPC(job.id, "ip")
		url := formater.CustomStr("http://{ip}:3001/status", map[string]any{"ip": ip})
		results[i] = protocols.GetPC(url, 1)
	}

	// 3) Apply results under Lock (short)
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

func updatePJ() {
	// 1) Build jobs under RLock (no network)
	jobs := buildPJJobs()

	// 2) Network phase without locks
	results := make([]int, len(jobs))
	for i, job := range jobs {
		addr := config.FindPJ(job.id, job.key) // твоя функция ожидает zoneID + "pj_1" и т.п.
		results[i] = protocols.GetPjlink(addr)
	}

	// 3) Apply results under Lock (short)
	statusMu.Lock()
	for i, job := range jobs {
		v := results[i]

		if job.innerIdx == -1 {
			ensureZoneStatusMap(job.zoneIdx)
			// пишем только если ключ реально существует/нужен
			statusData.Zones[job.zoneIdx].Status[job.key] = v
			continue
		}

		ensureInnerStatusMap(job.zoneIdx, job.innerIdx)
		statusData.Zones[job.zoneIdx].InnerZones[job.innerIdx].Status[job.key] = v
	}
	statusMu.Unlock()

	logger.Debug.Println("pj statuses updated")
}

func buildPCJobs() []pcJob {
	statusMu.RLock()
	defer statusMu.RUnlock()

	jobs := make([]pcJob, 0, len(statusData.Zones))

	for i, z := range statusData.Zones {
		// твой хардкод пропуска
		if z.ID == "relay" || z.ID == "light" {
			continue
		}

		// Обновляем pc_1 для зоны, если у неё вообще есть status-map (как у тебя было)
		if z.Status != nil {
			jobs = append(jobs, pcJob{zoneIdx: i, innerIdx: -1, id: z.ID})
		}

		// InnerZones — всегда пытаемся обновлять pc_1 (как было раньше)
		for j := range z.InnerZones {
			jobs = append(jobs, pcJob{zoneIdx: i, innerIdx: j, id: z.InnerZones[j].ID})
		}
	}

	return jobs
}

func buildPJJobs() []pjJob {
	statusMu.RLock()
	defer statusMu.RUnlock()

	var jobs []pjJob

	for i := range statusData.Zones {
		z := statusData.Zones[i]

		// zone level pj_*
		for key := range z.Status {
			if pjKeyRe.MatchString(key) {
				jobs = append(jobs, pjJob{zoneIdx: i, innerIdx: -1, id: z.ID, key: key})
			}
		}

		// inner level pj_*
		for j := range z.InnerZones {
			in := z.InnerZones[j]
			for key := range in.Status {
				if pjKeyRe.MatchString(key) {
					jobs = append(jobs, pjJob{zoneIdx: i, innerIdx: j, id: in.ID, key: key})
				}
			}
		}
	}

	// стабильный порядок (приятно для логов/отладки)
	sort.Slice(jobs, func(a, b int) bool {
		if jobs[a].id == jobs[b].id {
			return jobs[a].key < jobs[b].key
		}
		return jobs[a].id < jobs[b].id
	})

	return jobs
}

func ensureZoneStatusMap(zoneIdx int) {
	if statusData.Zones[zoneIdx].Status == nil {
		statusData.Zones[zoneIdx].Status = map[string]int{}
	}
}

func ensureInnerStatusMap(zoneIdx, innerIdx int) {
	if statusData.Zones[zoneIdx].InnerZones[innerIdx].Status == nil {
		statusData.Zones[zoneIdx].InnerZones[innerIdx].Status = map[string]int{}
	}
}

// -------------------- Persistence (atomic files) --------------------
func persistSnapshots() error {
	persistMu.Lock()
	defer persistMu.Unlock()

	statusMu.RLock()
	defer statusMu.RUnlock()

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

	if err := os.WriteFile(statusYAMLPath, yamlBytes, 0644); err != nil {
		logger.Warn.Println(err)
		return err
	}

	if err := os.WriteFile(statusJSONPath, jsonBytes, 0644); err != nil {
		logger.Warn.Println(err)
		return err
	}

	return nil
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	tmp := filepath.Join(dir, "."+base+".tmp")
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// // -------------------- Optional: utility --------------------

// // If you still want to count pj keys somewhere else.
// func countPjKeys(statusMap map[string]int) int {
// 	if statusMap == nil {
// 		return 0
// 	}
// 	count := 0
// 	for k := range statusMap {
// 		if pjKeyRe.MatchString(k) {
// 			count++
// 		}
// 	}
// 	return count
// }
