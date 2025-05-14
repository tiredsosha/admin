package protocols

import (
	"strconv"

	"github.com/LightInstruments/pjlink"
	"github.com/tiredsosha/admin/tools/logger"
)

func SendPjlink(ip, command string) {
	proj := pjlink.NewProjector(ip, "")

	switch command {
	case "on":
		if err := proj.TurnOn(); err != nil {
			logger.Error.Printf("couldn't send execute pjlink %v", err)
		}
	case "off":
		if err := proj.TurnOff(); err != nil {
			logger.Error.Printf("couldn't send execute pjlink %v", err)
		}
	}

}

func GetPjlink(ip string) int {
	response := 520
	proj := pjlink.NewProjector(ip, "")
	status, err := proj.GetPowerStatus()
	if err != nil {
		logger.Error.Printf("couldn't send execute pjlink %v", err)
	} else {
		boolStatus, _ := strconv.ParseBool(status.Response[0])
		if boolStatus {
			response = 200
		} else {
			response = 521
		}
	}
	return response
}
