package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
	config "github.com/tiredsosha/admin/tools/configurator"
	"github.com/tiredsosha/admin/tools/formater"
	"github.com/tiredsosha/admin/tools/logger"
	"gopkg.in/yaml.v3"
)

type Status struct {
	Zones []Zone `yaml:"zones"`
}

type Zone struct {
	ID         string         `yaml:"id"`
	InnerZones []InnerZone    `yaml:"innerZones,omitempty"`
	Status     map[string]int `yaml:"status,omitempty"`
}

type InnerZone struct {
	ID     string         `yaml:"id"`
	Status map[string]int `yaml:"status"`
}

var (
	// control concurrent goroutine launches
	mu        sync.Mutex
	pcRunning bool
	pjRunning bool

	// shared status data
	statusData Status

	// protects statusData (writers use Lock, readers use RLock)
	statusMu sync.RWMutex
)

// HTTP handler: returns the latest JSON snapshot
func statusPark(c *gin.Context) {
	jsonBytes, err := os.ReadFile("./configs/status.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read status data"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonBytes)
}

// Load initial YAML into statusData (protected)
func StatusInit() {
	data, err := os.ReadFile("./configs/status.yaml")
	if err != nil {
		fmt.Println("Error reading YAML:", err)
		return
	}

	var loaded Status
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return
	}

	statusMu.Lock()
	statusData = loaded
	statusMu.Unlock()
}

// --- Updaters ---------------------------------------------------------------

func updatePC() {
	statusMu.Lock()
	defer statusMu.Unlock()

	for i, zone := range statusData.Zones {
		if zone.ID == "relay" || zone.ID == "light" {
			continue
		}

		if zone.Status != nil {
			statusData.Zones[i].Status["pc_1"] = protocols.GetPC(
				formater.CustomStr("http://{ip}:3001/status",
					map[string]any{"ip": config.FindPC(zone.ID, "ip")}),
				1,
			)
		}

		for j, inner := range zone.InnerZones {
			statusData.Zones[i].InnerZones[j].Status["pc_1"] = protocols.GetPC(
				formater.CustomStr("http://{ip}:3001/status",
					map[string]any{"ip": config.FindPC(inner.ID, "ip")}),
				1,
			)
		}
	}
	logger.Debug.Println("pc statuses updated")
}

func countPjKeys(statusMap map[string]int) int {
	count := 0
	re := regexp.MustCompile(`^pj_\d+$`)
	for k := range statusMap {
		if re.MatchString(k) {
			count++
		}
	}
	return count
}

func updatePJ() {
	statusMu.Lock()
	defer statusMu.Unlock()

	for i := range statusData.Zones {
		zone := &statusData.Zones[i]

		maxKeys := countPjKeys(zone.Status)
		for n := 1; n <= maxKeys; n++ {
			key := fmt.Sprintf("pj_%d", n)
			if _, exists := zone.Status[key]; exists {
				zone.Status[key] = protocols.GetPjlink(config.FindPJ(zone.ID, key))
			}
		}

		for j := range zone.InnerZones {
			inner := &zone.InnerZones[j]
			maxKeysInner := countPjKeys(inner.Status)
			for n := 1; n <= maxKeysInner; n++ {
				key := fmt.Sprintf("pj_%d", n)
				if _, exists := inner.Status[key]; exists {
					logger.Debug.Println(config.FindPJ(inner.ID, key))
					inner.Status[key] = protocols.GetPjlink(config.FindPJ(inner.ID, key))
				}
			}
		}
	}
	logger.Debug.Println("pj statuses updated")
}

// --- Serialization helpers --------------------------------------------------

func convertYamlToJson(yamlData []byte) ([]byte, error) {
	var parsedData any
	if err := yaml.Unmarshal(yamlData, &parsedData); err != nil {
		return nil, err
	}
	jsonData, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// --- Public orchestrator ----------------------------------------------------

func UpdateStatues() {
	mu.Lock()
	if !pcRunning {
		pcRunning = true
		go func() {
			defer func() {
				mu.Lock()
				pcRunning = false
				mu.Unlock()
			}()
			updatePC()
		}()
	}
	// if !pjRunning {
	// 	pjRunning = true
	// 	go func() {
	// 		defer func() {
	// 			mu.Lock()
	// 			pjRunning = false
	// 			mu.Unlock()
	// 		}()
	// 		updatePJ()
	// 	}()
	// }
	mu.Unlock()

	update()
}

// Persist current statusData to YAML and JSON (read-only path)
func update() {
	statusMu.RLock()
	updatedYAML, err := yaml.Marshal(statusData)
	statusMu.RUnlock()
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	if err := os.WriteFile("./configs/status.yaml", updatedYAML, 0644); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}
	logger.Debug.Println("yaml updated, check status.yaml")

	jsonBytes, err := convertYamlToJson(updatedYAML)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	if err := os.WriteFile("./configs/status.json", jsonBytes, 0644); err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}
	logger.Debug.Println("json updated, check status.json")
}
