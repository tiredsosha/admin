package protocols

import (
	"github.com/LightInstruments/pjlink"
	"github.com/tiredsosha/admin/tools/logger"
)

func SendPjlink(ip, command string) {
	proj := pjlink.NewProjector(ip, "")

	switch command {
	case "on":
		if err := proj.TurnOn(); err != nil {
			logger.Error.Println(err)
		}
	case "off":
		if err := proj.TurnOff(); err != nil {
			logger.Error.Println(err)
		}
	}

}
