package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func router(server *gin.Engine) {
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

	status := server.Group("/status")
	{
		status.GET("/park", statusPark)
	}

	debug := server.Group("/debug")
	{
		debug.GET("/getTest", testGet)
		debug.POST("/postTest", testPost)
	}
}

func StartServer(port int) {
	portSrt := ":" + strconv.Itoa(port)
	route := gin.Default()
	router(route)

	server := &http.Server{
		Addr:           portSrt,
		Handler:        route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
