package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tiredsosha/admin/protocols"
	"github.com/tiredsosha/admin/tools/formater"
	"github.com/tiredsosha/admin/tools/logger"

	config "github.com/tiredsosha/admin/tools/configurator"

	"net"
)

func powerPc(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if data.Zone == "faces" || data.Zone == "art" || data.Zone == "city" {
		conn, err := net.Dial("udp", "255.255.255.255:5000")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		conn.Write([]byte(data.Command))

	}

	logger.Info.Println("request data -", data)

	if data.Command == "on" {
		protocols.SendWOL(config.FindPC(data.Zone, "mac"))
	} else {
		protocols.SendGet(formater.CustomStr(
			"http://{ip}:3001/off",
			map[string]any{"ip": config.FindPC(data.Zone, "ip")}), 2,
		)
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerProjector(c *gin.Context) {
	var data JsonCommand

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	protocols.SendPjlink(config.FindPJ(data.Zone, data.ID), data.Command)

	// fmt.Println(data.Zone)

	// вот это тока для рязани исключение, в других проектах надо его убирать

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerRelay(c *gin.Context) {
	var data JsonCommand
	command := "0"

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	zoneRelay := config.FindRelay(data.Zone)
	for _, ip := range zoneRelay {
		protocols.SendGet(formater.CustomStr(
			"http://admin:admin@{ip}/protect/rb0{command}.cgi",
			map[string]any{"ip": ip, "command": command},
		), 2,
		)
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func powerZone(c *gin.Context) {
	var data JsonCommand
	command := "0"

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	if data.Command == "on" {
		command = "n"
		for range 4 {
			protocols.SendWOL(config.FindPC(data.Zone, "mac"))
		}
	} else {
		command = "f"
		protocols.SendGet(formater.CustomStr(
			"http://{ip}:3001/off",
			map[string]any{"ip": config.FindPC(data.Zone, "ip")}), 2,
		)
	}

	zonePJ := config.FindZonePJ(data.Zone)
	for _, ip := range zonePJ {
		protocols.SendPjlink(ip, data.Command)
	}

	zoneRelay := config.FindRelay(data.Zone)
	for _, ip := range zoneRelay {
		protocols.SendGet(formater.CustomStr(
			"http://admin:admin@{ip}/protect/rb0{command}.cgi",
			map[string]any{"ip": ip, "command": command},
		), 2,
		)
	}

	if data.Zone == "faces" {
		conn, err := net.Dial("udp", "192.168.10.56:5000")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		conn.Write([]byte(data.Command))

		conn2, err := net.Dial("udp", "192.168.10.57:5000")
		if err != nil {
			panic(err)
		}
		defer conn2.Close()

		conn2.Write([]byte(data.Command))

	}

	if data.Zone == "art" {
		conn, err := net.Dial("udp", "192.168.10.51:5000")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		conn.Write([]byte(data.Command))

		conn2, err := net.Dial("udp", "192.168.10.52:5000")
		if err != nil {
			panic(err)
		}
		defer conn2.Close()

		conn2.Write([]byte(data.Command))

		conn3, err := net.Dial("udp", "192.168.10.53:5000")
		if err != nil {
			panic(err)
		}
		defer conn3.Close()

		conn3.Write([]byte(data.Command))

		conn4, err := net.Dial("udp", "192.168.10.54:5000")
		if err != nil {
			panic(err)
		}
		defer conn4.Close()

		conn4.Write([]byte(data.Command))

		conn5, err := net.Dial("udp", "192.168.10.55:5000")
		if err != nil {
			panic(err)
		}
		defer conn5.Close()

		conn5.Write([]byte(data.Command))

	}

	if data.Zone == "city" {
		conn, err := net.Dial("udp", "192.168.10.50:5000")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		conn.Write([]byte(data.Command))
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}

// func powerPark(c *gin.Context) {
// 	var data JsonNoID

// 	// Bind JSON and validate
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		logger.Error.Println("Invalid input:", err)
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	logger.Info.Println("request data -", data)

// 	if data.Command == "on" {
// 		go func() {
// 			for _, mac := range config.ALLMAC {
// 				for range 4 {
// 					protocols.SendWOL(mac)
// 				}
// 			}

// 			for _, pj := range config.ALLPJ {
// 				protocols.SendPjlink(pj, data.Command)
// 			}
// 		}()

// 	} else if data.Command == "off" {

// 		go func() {
// 			for _, pc := range config.ALLPC {
// 				protocols.SendGet(formater.CustomStr(
// 					"http://{ip}:3001/off",
// 					map[string]any{"ip": pc}), 2,
// 				)
// 			}

// 			for _, ip := range config.ALLPJ {
// 				protocols.SendPjlink(ip, data.Command)
// 			}
// 		}()

// 	} else if data.Command == "restart" {
// 		go func() {
// 			for _, pc := range config.ALLPC {
// 				protocols.SendGet(formater.CustomStr(
// 					"http://{ip}:3001/restart",
// 					map[string]any{"ip": pc}), 2,
// 				)
// 			}
// 		}()
// 	}

// 	c.JSON(200, gin.H{
// 		"message": "OK",
// 	})
// }

func powerPark(c *gin.Context) {
	var data JsonNoID

	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error.Println("Invalid input:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Println("request data -", data)

	switch data.Command {
	case "on":

		command := "n"

		// Wake MACs
		go func() {
			for _, mac := range config.ALLMAC {
				for i := 0; i < 4; i++ {
					protocols.SendWOL(mac)
				}
			}
		}()

		// Turn on projectors
		go func() {
			for _, pj := range config.ALLPJ {

				// вот это тока для рязани исключение, в других проектах надо его убирать
				if pj != "172.16.3.73" {
					protocols.SendPjlink(pj, "on")
				}
			}
		}()

		// Turn on relay
		go func() {
			protocols.SendGet(formater.CustomStr(
				"http://admin:admin@{ip}/protect/rb0{command}.cgi",
				map[string]any{"ip": "172.16.3.76", "command": command},
			), 2,
			)
		}()

		// // Turn defaults on lights
		// go func() {
		// 	for _, dali := range config.ALLDALI {
		// 		out, lamps, defaults := config.FindDali(dali)
		// 		protocols.ArlightDefault(out, lamps, defaults)
		// 	}
		// }()

	case "off":

		command := "f"

		// Power off PCs by IP
		go func() {
			for _, ip := range config.ALLPC {
				url := formater.CustomStr("http://{ip}:3001/off", map[string]any{"ip": ip})
				protocols.SendGet(url, 2)
			}
		}()

		// Turn off projectors
		go func() {
			for _, pj := range config.ALLPJ {
				protocols.SendPjlink(pj, "off")
			}
		}()

		// Turn off relay
		go func() {
			protocols.SendGet(formater.CustomStr(
				"http://admin:admin@{ip}/protect/rb0{command}.cgi",
				map[string]any{"ip": "172.16.3.76", "command": command},
			), 2,
			)
		}()

	case "restart":
		// Restart PCs by IP
		go func() {
			for _, ip := range config.ALLPC {
				url := formater.CustomStr("http://{ip}:3001/restart", map[string]any{"ip": ip})
				protocols.SendGet(url, 2)
			}
		}()

	default:
		logger.Warn.Printf("unknown command for park power: %s", data.Command)
	}

	c.JSON(200, gin.H{"message": "OK"})
}
