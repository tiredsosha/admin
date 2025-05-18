package configurator

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiredsosha/admin/tools/logger"
)

var PC map[string]map[string]string
var ALLPC []string
var ALLMAC []string

func configPC() error {
	// Read the YAML file
	data, err := os.ReadFile("./configs/configPC.yaml")
	if err != nil {
		logger.Error.Printf("error reading configPC.yaml: %v", err)
		return err
	}

	// Unmarshal YAML into the map
	err = yaml.Unmarshal(data, &PC)
	if err != nil {
		logger.Error.Printf("error unmarshaling configPC.yaml: %v", err)
		return err
	}
	FindAllPC()
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

func FindAllPC() {
	// Collect MACs and IPs
	for _, node := range PC {
		if mac, ok := node["mac"]; ok {
			ALLMAC = append(ALLMAC, mac)
		}
		if ip, ok := node["ip"]; ok {
			ALLPC = append(ALLPC, ip)
		}
	}
	logger.Info.Println("all pc ip listed")
	logger.Info.Println("all pc mac listed")

	// logger.Info.Printf("all pc ip - '%v'\n", ALLPC)
	// logger.Info.Printf("all pc mac - '%v'\n", ALLMAC)
}
