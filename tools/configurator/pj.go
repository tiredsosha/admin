package configurator

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiredsosha/admin/tools/logger"
)

var PJ map[string]map[string]string

func configPJ() error {
	// Read the YAML file
	data, err := os.ReadFile("./configs/configPJ.yaml")
	if err != nil {
		logger.Error.Printf("error reading file: %v", err)
		return err
	}

	// Unmarshal YAML into the map
	err = yaml.Unmarshal(data, &PJ)
	if err != nil {
		logger.Error.Printf("error unmarshaling YAML: %v", err)
		return err
	}
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
