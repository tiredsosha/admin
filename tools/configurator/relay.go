package configurator

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiredsosha/admin/tools/logger"
)

var Relays map[string]map[string]string

func configRelay() error {
	// Read the YAML file
	data, err := os.ReadFile("./configs/configRelay.yaml")
	if err != nil {
		logger.Error.Printf("error reading configRelay.yaml: %v", err)
		return err
	}

	// Unmarshal YAML into the map
	err = yaml.Unmarshal(data, &Relays)
	if err != nil {
		logger.Error.Printf("error unmarshaling configRelay.yaml: %v", err)
		return err
	}
	return err
}

// func FindRelay(id string) string {

// 	if value, exists := Relays[id]; exists {
// 		logger.Info.Printf("value for key '%s' - '%s'\n", id, value)
// 		return value
// 	} else {
// 		logger.Error.Printf("key '%s' not found\n", id)
// 		return "none"
// 	}

// }

func FindRelay(main string) []string {

	// Get the map for the zone
	zoneMap, exists := Relays[main]
	if !exists {
		logger.Debug.Printf("zone unfound '%s'\n", main)
		return []string{}
	}

	// Collect all IPs into a slice
	var ipList []string
	for _, ip := range zoneMap {
		ipList = append(ipList, ip)
	}
	logger.Debug.Printf("re;lay ips in zone '%s': %v\n", main, ipList)
	return ipList
}
