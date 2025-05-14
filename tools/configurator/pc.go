package configurator

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiredsosha/admin/tools/logger"
)

var PC map[string]map[string]string

func configPC() error {
	// Read the YAML file
	data, err := os.ReadFile("./configs/configPC.yaml")
	if err != nil {
		logger.Error.Printf("error reading file: %v", err)
		return err
	}

	// Unmarshal YAML into the map
	err = yaml.Unmarshal(data, &PC)
	if err != nil {
		logger.Error.Printf("error unmarshaling YAML: %v", err)
		return err
	}
	return err
}

func FindPC(id, command string) string {

	if value, exists := PC[id][command]; exists {
		logger.Info.Printf("value for key '%s' - '%s'\n", id, value)
		return value
	} else {
		logger.Error.Printf("key '%s' not found\n", id)
		return "none"
	}
}
