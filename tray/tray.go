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
	menuOffPark := systray.AddMenuItem("PARK OFF", "Turn off whole park")
	menuOnPark := systray.AddMenuItem("PARK ON", "Turn on whole park")
	menuQuit := systray.AddMenuItem("QUIT", "Quit the whole app")

	go func() {
		<-menuQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		<-menuOnPark.ClickedCh
		protocols.SendPost("http://127.0.0.1:8080/power/park", "off")
	}()
	go func() {
		<-menuOffPark.ClickedCh
		protocols.SendPost("http://127.0.0.1:8080/power/park", "on")
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
