package configurator

import (
	"strings"

	"github.com/tiredsosha/admin/tools/logger"
)

type LightEntry struct {
	Address uint16 `yaml:"address"`
	Default uint16 `yaml:"default"`
}

type Light struct {
	Out    string       `yaml:"out"`
	Lights []LightEntry `yaml:"lights"`
}

type DALIConfig map[string]Light

var Dali DALIConfig
var ALLDALI []string

// Load YAML config once
// func configDALI() error {
// 	// Read the YAML file
// 	data, err := os.ReadFile("./configs/configDali.yaml")
// 	if err != nil {
// 		logger.Error.Printf("error reading configRelay.yaml: %v", err)
// 		return err
// 	}

// 	// Unmarshal YAML into the global `Dali` map
// 	err = yaml.Unmarshal(data, &Dali)
// 	if err != nil {
// 		logger.Error.Printf("error unmarshaling configRelay.yaml: %v", err)
// 		return err
// 	}

// 	// Reset ALLDALI in case configDALI is called more than once
// 	ALLDALI = make([]string, 0, len(Dali))

// 	for zone := range Dali {
// 		ALLDALI = append(ALLDALI, zone)
// 	}

// 	logger.Debug.Printf("loaded zones: %v", ALLDALI)
// 	return nil
// }

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

func FindDali(zone string) (busID byte, addresses, defaults []uint16) {
	data, exists := Dali[zone]
	if !exists {
		logger.Debug.Printf("light zone not found '%s'\n", zone)
		return 0, []uint16{}, []uint16{}
	}

	for _, light := range data.Lights {
		addresses = append(addresses, light.Address)
		defaults = append(defaults, light.Default)
	}

	logger.Debug.Printf("zone '%s' => out: %s, addresses: %v, defaults: %v\n", zone, data.Out, addresses, defaults)
	return getBusID(data.Out), addresses, defaults
}
