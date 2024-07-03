package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
)

func lightMode(c *gin.Context) {
	var data JsonID
	c.Bind(&data)

	protocols.SendUdp("127.0.0.1", 8090, "lights")
	// c.JSON(200, gin.H{
	// 	"a": b.NestedStruct,
	// 	"b": b.FieldB,
	// })

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func lightLamp(c *gin.Context) {
	var data JsonID
	c.Bind(&data)

	protocols.SendUdp("127.0.0.1", 8090, "lights")
	// c.JSON(200, gin.H{
	// 	"a": b.NestedStruct,
	// 	"b": b.FieldB,
	// })

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func lightBrightness(c *gin.Context) {
	var data JsonID
	c.Bind(&data)

	protocols.SendUdp("127.0.0.1", 8090, "lights")
	// c.JSON(200, gin.H{
	// 	"a": b.NestedStruct,
	// 	"b": b.FieldB,
	// })

	c.JSON(200, gin.H{
		"message": "pong",
	})
}