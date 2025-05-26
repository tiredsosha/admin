package tray

import (
	"time"

	"github.com/tiredsosha/admin/protocols"
	"github.com/tiredsosha/admin/tools/logger"

	"github.com/getlantern/systray"
)

var Conn bool = false
var icon bool

func onReady() {

	systray.SetIcon(grayIcon)
	systray.SetTitle("Executor Server")

	systray.SetTooltip("Executor")
	menuOnLights := systray.AddMenuItem("LIGHT ON", "Turn on all the dali lights")
	menuOffLights := systray.AddMenuItem("LIGHT OFF", "Turn off all the dali lights")
	menuOffPark := systray.AddMenuItem("PARK OFF", "Turn off whole park")
	menuOnPark := systray.AddMenuItem("PARK ON", "Turn on whole park")

	menuQuit := systray.AddMenuItem("QUIT", "Quit the whole app")

	go func() {
		<-menuQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		<-menuOnPark.ClickedCh
		protocols.SendPost("http://127.0.0.1:8080/power/park", "off", 2)
	}()
	go func() {
		<-menuOffPark.ClickedCh
		protocols.SendPost("http://127.0.0.1:8080/power/park", "on", 2)
	}()

	go func() {
		<-menuOffLights.ClickedCh
		protocols.SendGet("http://127.0.0.1:8080/light/off", 2)
	}()

	go func() {
		<-menuOnLights.ClickedCh
		protocols.SendGet("http://127.0.0.1:8080/light/default", 2)
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)

			if Conn == icon {
				continue
			}
			if Conn {
				systray.SetIcon(blueIcon)
			} else {
				systray.SetIcon(grayIcon)
			}

			icon = Conn
		}
	}()
}

func onExit() {
	logger.Error.Fatal("EXITING MANUALLY")
}

func TrayStart() {
	systray.Run(onReady, onExit)
}
