package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
	"github.com/tiredsosha/admin/tools/logger"

	config "github.com/tiredsosha/admin/tools/configurator"
)

func brightnessChange(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	out, lamps, _ := config.FindDali(data.Zone)
	protocols.ArlightControl(out, lamps, data.Command)

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func brightnessDefault(c *gin.Context) {

	logger.Info.Println("request data all light to default")

	go func() {
		for _, dali := range config.ALLDALI {
			out, lamps, defaults := config.FindDali(dali)
			protocols.ArlightDefault(out, lamps, defaults)
		}
	}()

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func brightnessOFF(c *gin.Context) {

	logger.Info.Println("request data all light to default")

	// Turn off lights
	go func() {
		for _, dali := range config.ALLDALI {
			out, lamps, _ := config.FindDali(dali)
			protocols.ArlightControl(out, lamps, "0")
		}
	}()

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
