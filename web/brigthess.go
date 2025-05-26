package web

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/tools/logger"

	config "github.com/tiredsosha/admin/tools/configurator"

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

func brightnessChange(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	out, lamps := config.FindDali(data.Zone)
	command := mapTo254(data.Command)
	if command == 500 {
		c.JSON(400, gin.H{"error": "Invalid brightness value"})
		return
	} else {
		if out != 0 || len(lamps) != 0 {
			handler := modbus.NewTCPClientHandler("10.1.31.251:502")
			handler.SlaveId = out
			handler.Timeout = 4 * time.Second
			err := handler.Connect()
			if err != nil {
				logger.Error.Printf("Connection error: %v", err)
			}
			defer handler.Close()

			client := modbus.NewClient(handler)

			for _, addr := range lamps {
				reg := 256 + addr
				_, err := client.WriteSingleRegister(reg, command)
				if err != nil {
					logger.Error.Printf("Failed to write lamp %d: %v", addr, err)
				} else {
					logger.Info.Printf("Set lamp %d to %s\n", addr, data.Command)
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
