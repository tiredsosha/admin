package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/tray"
)

func router(server *gin.Engine) {
	power := server.Group("/power")
	{
		power.POST("/park", powerPark)
		power.POST("/zone", powerZone)
		power.POST("/pc", powerPc)
		power.POST("/projector", powerProjector)
		power.POST("/relay", powerRelay)
	}

	restart := server.Group("/restart")
	{
		restart.POST("/pc", restartPc)
	}

	sound := server.Group("/sound")
	{
		sound.POST("/zone", soundZone)
	}
	light := server.Group("/light")
	{
		light.POST("/brightness", brightnessChange)
		light.GET("/off", brightnessOFF)
		light.GET("/default", brightnessDefault)
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
	tray.Conn = true

	server.ListenAndServe()
}
