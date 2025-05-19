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

// func GetPjlink(ip string) int {
// 	response := 520
// 	proj := pjlink.NewProjector(ip, "")
// 	status, err := proj.GetPowerStatus()
// 	if err != nil {
// 		logger.Error.Printf("couldn't send execute pjlink %v", err)
// 	} else {
// 		if len(status.Response) > 0 {
// 			//logger.Info.Println("pjlink status -", status.Response)
// 			boolStatus, _ := strconv.ParseBool(status.Response[0])
// 			if boolStatus {
// 				response = 200
// 			} else {
// 				response = 521
// 			}
// 		}
// 	}
// 	logger.Info.Println("pjlink status -", response)
// 	return response
// }

func GetPjlink(ip string) int {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic or handle it as needed
			logger.Error.Printf("Recovered from panic in GetPjlink: %v", r)
		}
	}()

	response := 520
	proj := pjlink.NewProjector(ip, "")
	status, err := proj.GetPowerStatus()
	if err != nil {
		logger.Error.Printf("couldn't send execute pjlink %v", err)
	} else {
		if len(status.Response) > 0 {
			boolStatus, err := strconv.ParseBool(status.Response[0])
			if err != nil {
				logger.Error.Printf("Error parsing bool from response: %v", err)
			} else {
				if boolStatus {
					response = 200
				} else {
					response = 521
				}
			}
		} else {
			logger.Error.Println("Empty response received")
		}
	}
	// logger.Info.Println("pjlink status -", response)
	return response
}
