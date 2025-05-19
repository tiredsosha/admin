package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
	"github.com/tiredsosha/admin/tools/formater"
	"github.com/tiredsosha/admin/tools/logger"

	config "github.com/tiredsosha/admin/tools/configurator"
)

func powerPc(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	if data.Command == "on" {
		protocols.SendWOL(config.FindPC(data.Zone, "mac"))
	} else {
		protocols.SendGet(formater.CustomStr(
			"http://{ip}:3001/off",
			map[string]any{"ip": config.FindPC(data.Zone, "ip")}),
		)
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerProjector(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	protocols.SendPjlink(config.FindPJ(data.Zone, data.ID), data.Command)

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerRelay(c *gin.Context) {
	var data JsonCommand
	command := "0"

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	if data.Command == "on" {
		command = "f"
	} else {
		command = "n"
	}

	protocols.SendGet(formater.CustomStr(
		"http://admin:admin@{ip}/protect/rb0{command}.cgi",
		map[string]any{"ip": config.FindRelay(data.Zone), "command": command},
	),
	)

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerZone(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	if data.Command == "on" {
		for range 4 {
			protocols.SendWOL(config.FindPC(data.Zone, "mac"))
		}
	} else {
		protocols.SendGet(formater.CustomStr(
			"http://{ip}:3001/off",
			map[string]any{"ip": config.FindPC(data.Zone, "ip")}),
		)
	}

	zonePJ := config.FindZonePJ(data.Zone)
	for _, ip := range zonePJ {
		protocols.SendPjlink(ip, data.Command)
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerPark(c *gin.Context) {
	var data JsonNoID

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	if data.Command == "on" {
		go func() {
			for _, mac := range config.ALLMAC {
				for range 4 {
					protocols.SendWOL(mac)
				}
			}

			for _, pj := range config.ALLPJ {
				protocols.SendPjlink(pj, data.Command)
			}
		}()

	} else if data.Command == "off" {

		go func() {
			for _, pc := range config.ALLPC {
				protocols.SendGet(formater.CustomStr(
					"http://{ip}:3001/off",
					map[string]any{"ip": pc}),
				)
			}

			for _, ip := range config.ALLPJ {
				protocols.SendPjlink(ip, data.Command)
			}
		}()

	} else if data.Command == "restart" {
		go func() {
			for _, pc := range config.ALLPC {
				protocols.SendGet(formater.CustomStr(
					"http://{ip}:3001/restart",
					map[string]any{"ip": pc}),
				)
			}
		}()
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
