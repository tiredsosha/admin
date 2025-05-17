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
		go protocols.SendWOL(config.FindPC(data.Zone, "mac"))
	} else {
		go protocols.SendGet(formater.CustomStr(
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

	go protocols.SendPjlink(config.FindPJ(data.Zone, data.ID), data.Command)

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

	go protocols.SendGet(formater.CustomStr(
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
		go protocols.SendWOL(config.FindPC(data.Zone, "mac"))
	} else {
		go protocols.SendGet(formater.CustomStr(
			"http://{ip}:3001/off",
			map[string]any{"ip": config.FindPC(data.Zone, "ip")}),
		)
	}

	zonePJ := config.FindZonePJ(data.Zone)
	for _, ip := range zonePJ {
		go protocols.SendPjlink(ip, data.Command)
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

		for _, mac := range config.ALLMAC {
			go protocols.SendWOL(mac)
		}

		for _, pj := range config.ALLPJ {
			go protocols.SendPjlink(pj, data.Command)
		}
	} else if data.Command == "off" {

		for _, pc := range config.ALLPC {
			go protocols.SendGet(formater.CustomStr(
				"http://{ip}:3001/off",
				map[string]any{"ip": pc}),
			)
		}

		for _, ip := range config.ALLPJ {
			go protocols.SendPjlink(ip, data.Command)
		}
	} else if data.Command == "restart" {
		for _, pc := range config.ALLPC {
			go protocols.SendGet(formater.CustomStr(
				"http://{ip}:3001/restart",
				map[string]any{"ip": pc}),
			)
		}
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
