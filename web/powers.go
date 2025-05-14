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

	logger.Info.Println(data)

	if data.Command == "on" {
		go protocols.SendWOL(config.FindPC(data.Zone, "mac"))
	} else {
		go protocols.SendGet(formater.CustomStr(
			"http://{ip}/power/off",
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

	logger.Info.Println(data)

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

	logger.Info.Println(data)

	if data.Command == "on" {
		command = "n"
	} else {
		command = "f"
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

	logger.Info.Println(data)

	if data.Command == "on" {
		go protocols.SendWOL(config.FindPC(data.Zone, "mac"))
	} else {
		go protocols.SendGet(formater.CustomStr(
			"http://{ip}/power/off",
			map[string]any{"ip": config.FindPC(data.Zone, "ip")}),
		)
	}
	// go protocols.SendPjlink(config.FindPJ(data.Zone, data.ID), data.Command)

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func powerPark(c *gin.Context) {
	var data JsonID
	c.Bind(&data)

	protocols.SendUdp("127.0.0.1", 8090, "restart")
	// c.JSON(200, gin.H{
	// 	"a": b.NestedStruct,
	// 	"b": b.FieldB,
	// })

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
