package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
	config "github.com/tiredsosha/admin/tools/configurator"
	"github.com/tiredsosha/admin/tools/formater"
	"github.com/tiredsosha/admin/tools/logger"
)

func soundZone(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	go protocols.SendPost(formater.CustomStr(
		"http://{ip}:3001/sound",
		map[string]any{"ip": config.FindPC(data.Zone, "ip")}), data.Command,
	)

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
