package protocols

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/LightInstruments/pjlink"
	"github.com/tiredsosha/admin/tools/logger"
	"github.com/tiredsosha/gopjlink"
)

// func SendPjlink(ip, command string) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			// Log the panic or handle it as needed
// 			logger.Error.Printf("Recovered from panic in GetPjlink: %v", r)
// 		}
// 	}()

// 	proj := pjlink.NewProjector(ip, "")
// 	// fmt.Println(ip, command)

// 	switch command {
// 	case "on":
// 		if err := proj.TurnOn(); err != nil {
// 			logger.Error.Printf("couldn't send execute pjlink %v", err)
// 		}
// 	case "off":
// 		if err := proj.TurnOff(); err != nil {
// 			logger.Error.Printf("couldn't send execute pjlink %v", err)
// 		}
// 	}

// }

func SendPjlink(ip, command string) {
	var boolCommand bool

	if command == "on" {
		boolCommand = true
	} else {
		boolCommand = false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Initialize your projector object with connection details
	proj := gopjlink.NewProjector(ip) // or however the constructor is defined

	// Send power off command
	err := proj.SetPower(ctx, boolCommand)
	if err != nil {
		fmt.Println("Failed to power off the projector:", err)
	} else {
		fmt.Println("Power off command sent successfully.")
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
