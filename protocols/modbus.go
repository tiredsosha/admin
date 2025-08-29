package protocols

import (
	"strconv"
	"time"

	"github.com/tiredsosha/admin/tools/logger"

	"github.com/goburrow/modbus"
)

func mapTo254(reqStr string) uint16 {
	reqInt, err := strconv.Atoi(reqStr)
	if err != nil {
		logger.Error.Println("invalid brightness value:", reqStr)
		return 500
	} else {
		if reqInt < 0 {
			reqInt = 0
		} else if reqInt > 100 {
			reqInt = 100
		}
		return uint16((reqInt * 254) / 100)
	}
}

func ArlightControl(out byte, lamps []uint16, commandStr string) {
	command := mapTo254(commandStr)

	if out != 0 || len(lamps) != 0 || command != 500 {
		handler := modbus.NewTCPClientHandler("10.1.0.121:502")
		handler.SlaveId = out
		handler.Timeout = 4 * time.Second
		err := handler.Connect()
		if err != nil {
			logger.Error.Printf("connection error: %v", err)
		}
		defer handler.Close()

		client := modbus.NewClient(handler)

		for _, addr := range lamps {
			reg := 256 + addr
			_, err := client.WriteSingleRegister(reg, command)
			if err != nil {
				logger.Error.Printf("failed to write lamp %d: %v", addr, err)
			} else {
				logger.Info.Printf("set lamp %d to %d\n", addr, command)
			}
		}
	}
}

func ArlightDefault(out byte, lamps, defaults []uint16) {
	if out != 0 || len(lamps) != 0 {
		handler := modbus.NewTCPClientHandler("10.1.31.251:502")
		handler.SlaveId = out
		handler.Timeout = 15 * time.Second
		err := handler.Connect()
		if err != nil {
			logger.Error.Printf("connection error: %v", err)
		}
		defer handler.Close()

		client := modbus.NewClient(handler)

		for i, addr := range lamps {
			reg := 256 + addr
			_, err := client.WriteSingleRegister(reg, defaults[i])
			if err != nil {
				logger.Error.Printf("failed to write lamp %d: %v", addr, err)
			} else {
				logger.Info.Printf("set lamp %d to %d\n", addr, defaults[i])
			}
		}
	}
}
