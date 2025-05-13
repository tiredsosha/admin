package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
)

func restartPc(c *gin.Context) {
	var data JsonID
	c.Bind(&data)

	protocols.SendUdp("127.0.0.1", 8010, "restart")
	// c.JSON(200, gin.H{
	// 	"a": b.NestedStruct,
	// 	"b": b.FieldB,
	// })

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
