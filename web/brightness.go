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
