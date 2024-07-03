package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
)

func powerPc(c *gin.Context) {
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

func powerProjector(c *gin.Context) {
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

func powerLidar(c *gin.Context) {
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

func powerCamera(c *gin.Context) {
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

func powerPhone(c *gin.Context) {
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

func powerZone(c *gin.Context) {
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
