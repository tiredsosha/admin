package configurator

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiredsosha/admin/tools/logger"
)

var PJ map[string]map[string]string
var ALLPJ []string

func configPJ() error {
	// Read the YAML file
	data, err := os.ReadFile("./configs/configPJ.yaml")
	if err != nil {
		logger.Error.Printf("error reading configPJ.yaml: %v", err)
		return err
	}

	// Unmarshal YAML into the map
	err = yaml.Unmarshal(data, &PJ)
	if err != nil {
		logger.Error.Printf("error unmarshaling configPJ.yaml: %v", err)
		return err
	}
	FindAllPJ()
	return err
}

func FindPJ(main, sub string) string {

	if mainID, ok := PJ[main]; ok {
		if ip, exists := mainID[sub]; exists {
			logger.Info.Printf("value for key '%s' in zone '%s' - '%s'\n", sub, main, ip)
			return ip
		} else {
			logger.Debug.Printf("key '%s' unfound in zone '%s'\n", sub, main)
			return "none"
		}
	} else {
		logger.Debug.Printf("zone unfound '%s'\n", main)
		return "none"
	}
}

func FindZonePJ(main string) []string {

	// Get the map for the zone
	zoneMap, exists := PJ[main]
	if !exists {
		logger.Debug.Printf("zone unfound '%s'\n", main)
		return []string{}
	}

	// Collect all IPs into a slice
	var ipList []string
	for _, ip := range zoneMap {
		ipList = append(ipList, ip)
	}
	logger.Debug.Printf("ips in zone '%s': %v\n", main, ipList)
	return ipList
}

func FindAllPJ() {

	// Collect all IPs
	ipSet := make(map[string]struct{}) // use a set to avoid duplicates
	for _, zone := range PJ {
		for _, ip := range zone {
			ipSet[ip] = struct{}{}
		}
	}

	// Convert set to list
	for ip := range ipSet {
		ALLPJ = append(ALLPJ, ip)
	}
	logger.Debug.Println("all pj ip listed")
	// logger.Debug.Printf("all pj ip - '%v'\n", ALLPJ)
}
