package web

import (
	"github.com/gin-gonic/gin"
)

func Router(server *gin.Engine) {
	power := server.Group("/power")
	{
		power.POST("/zone", powerZone)
		power.POST("/pc", powerPc)
		power.POST("/projector", powerProjector)
		power.POST("/lidar", powerLidar)
		power.POST("/camera", powerCamera)
		power.POST("/phone", powerPhone)
		power.POST("/park", powerPark)
	}

	restart := server.Group("/restart")
	{
		restart.POST("/player", restartPlayer)
		restart.POST("/content", restartContent)
		restart.POST("/app", restartApp)
		restart.POST("/controller", restartController)
	}

	light := server.Group("/light")
	{
		light.POST("/mode", lightMode)
		light.POST("/lamp", lightLamp)
		light.POST("/brightness", lightBrightness)
	}

	sound := server.Group("/sound")
	{
		sound.POST("/zone", soundZone)
	}

	content := server.Group("/content")
	{
		content.POST("/playmode", contentPlaymode)
	}
}
