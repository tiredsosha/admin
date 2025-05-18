package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

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
	Status     map[string]int `yaml:"status,omitempty"` // For zones with direct status
}

type InnerZone struct {
	ID     string         `yaml:"id"`
	Status map[string]int `yaml:"status"`
}

var statusData Status

func statusPark(c *gin.Context) {
	// Read the JSON data from the file
	jsonBytes, err := os.ReadFile("./configs/status.json")
	if err != nil {
		// Handle error: file not found or read error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read status data"})
		return
	}

	// Respond with the JSON data
	c.Data(http.StatusOK, "application/json", jsonBytes)

}

func StatusInit() {
	// Read your YAML file
	data, err := os.ReadFile("./configs/status.yaml")
	if err != nil {
		fmt.Println("Error reading YAML:", err)
		return
	}

	err = yaml.Unmarshal(data, &statusData)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return
	}

}

func updatePC() {
	// Loop through all zones and innerZones
	for i, zone := range statusData.Zones {
		if zone.ID != "relay" {
			// Check for direct status map
			if zone.Status != nil {
				// Update 'pc_1' status for the zone
				statusData.Zones[i].Status["pc_1"] = protocols.GetPC(formater.CustomStr(
					"http://{ip}:3001/status",
					map[string]any{"ip": config.FindPC(zone.ID, "ip")}),
				)
			}
			// Loop through innerZones
			for j, inner := range zone.InnerZones {
				// Update 'pc_1' status for each inner zone
				statusData.Zones[i].InnerZones[j].Status["pc_1"] = protocols.GetPC(formater.CustomStr(
					"http://{ip}:3001/status",
					map[string]any{"ip": config.FindPC(inner.ID, "ip")}),
				)
			}
		}
	}
	logger.Debug.Println("pc statuses updated")
}

func updateRelay() {
	for i, zone := range statusData.Zones {
		if zone.ID == "relay" {
			for j, inner := range zone.InnerZones {
				// Update status
				statusData.Zones[i].InnerZones[j].Status["controller_1"] = protocols.GetRelay(formater.CustomStr(
					"http://{ip}/pstat.xml",
					map[string]any{"ip": config.FindRelay(inner.ID)}),
				)
			}
			break
		}
	}
	logger.Debug.Println("relay statuses updated")
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
	// Loop through all zones
	for i, zone := range statusData.Zones {
		// Determine maxKeys for this zone based on existing keys
		maxKeys := countPjKeys(statusData.Zones[i].Status)

		// Loop through possible pj_* keys
		for n := 1; n <= maxKeys; n++ {
			key := fmt.Sprintf("pj_%d", n)
			if _, exists := statusData.Zones[i].Status[key]; !exists {
				statusData.Zones[i].Status[key] = protocols.GetPjlink(config.FindPJ(zone.ID, key))
			}
		}

		// Loop through inner zones
		for j := range zone.InnerZones {
			inner := &zone.InnerZones[j]
			// Count existing pj_* keys in inner zone
			maxKeysInner := countPjKeys(inner.Status)

			for n := 1; n <= maxKeysInner; n++ {
				key := fmt.Sprintf("pj_%d", n)
				if _, exists := inner.Status[key]; !exists {
					inner.Status[key] = protocols.GetPjlink(config.FindPJ(inner.ID, key))
				}
			}
		}
	}
	logger.Debug.Println("pj statuses updated")
}

func convertYamlToJson(yamlData []byte) ([]byte, error) {
	var parsedData any
	// Unmarshal YAML into an interface{}
	if err := yaml.Unmarshal(yamlData, &parsedData); err != nil {
		return nil, err
	}
	// Marshal the data into JSON with indentation
	jsonData, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func UpdateStatues() {
	updatePC()
	updateRelay()
	updatePJ()

	// Marshal updated data back to YAML
	updatedYAML, err := yaml.Marshal(&statusData)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	// Save the updated YAML to file
	err = os.WriteFile("./configs/status.yaml", updatedYAML, 0644)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	logger.Debug.Println("yaml updated, check status.yaml")

	// Convert YAML content to JSON
	jsonBytes, err := convertYamlToJson(updatedYAML)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	// Write JSON to a file
	err = os.WriteFile("./configs/status.json", jsonBytes, 0644)
	if err != nil {
		logger.Warn.Println(err)
		logger.Error.Fatal("EXITING")
	}

	logger.Debug.Println("json updated, check status.json")
}
