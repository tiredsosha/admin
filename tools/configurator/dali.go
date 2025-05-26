package configurator

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/tiredsosha/admin/tools/logger"
)

type Light struct {
	Out    string   `yaml:"out"`
	Lights []uint16 `yaml:"lights"`
}

type DALIConfig map[string]Light

var Dali DALIConfig

// Load YAML config once
func configDALI() error {
	// Read the YAML file
	data, err := os.ReadFile("./configs/configDali.yaml")
	if err != nil {
		logger.Error.Printf("error reading configRelay.yaml: %v", err)
		return err
	}

	// Unmarshal YAML into the global `Dali` map
	err = yaml.Unmarshal(data, &Dali)
	if err != nil {
		logger.Error.Printf("error unmarshaling configRelay.yaml: %v", err)
		return err
	}
	return nil
}

func getBusID(out string) byte {
	switch strings.ToUpper(out) {
	case "A":
		return 1
	case "B":
		return 2
	case "C":
		return 3
	case "D":
		return 4
	default:
		return 1
	}
}

func FindDali(zone string) (byte, []uint16) {
	// Get the data for the zone
	data, exists := Dali[zone]
	if !exists {
		logger.Debug.Printf("light zone not found '%s'\n", zone)
		return 0, []uint16{}
	}

	logger.Debug.Printf("zone '%s' => out: %s, lamps: %v\n", zone, data.Out, data.Lights)
	return getBusID(data.Out), data.Lights
}
