package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
)

func restartPlayer(c *gin.Context) {
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

func restartContent(c *gin.Context) {
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

func restartApp(c *gin.Context) {
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

func restartController(c *gin.Context) {
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
