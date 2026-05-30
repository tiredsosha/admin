package main

import (
	"time"

	"github.com/tiredsosha/admin/mosquitto"
	config "github.com/tiredsosha/admin/tools/configurator"
	"github.com/tiredsosha/admin/tools/logger"
	"github.com/tiredsosha/admin/tray"
	"github.com/tiredsosha/admin/web"
)

const (
	version = "1.0.0.beta.1"
)

func main() {
	logger.LogInit(true)

	cfg := config.ConfInit()
	config.ConfSubInit()

	web.StatusInit()

	hostname := "admin"
	topicPrefix := "executor/"

	logger.DebugLog(version, true, cfg.MqttOn, cfg.StatusOn, hostname, cfg.Broker, cfg.Username, cfg.Password, cfg.Port)

	// вот это нужно, но мы тоже будем с ним разбираться в конце
	go tray.TrayStart()

	if cfg.MqttOn {
		// так sub топик нам не нужен вообще, нам нужен только пуб топик. Мы наш стстум отправлять не будет. Так же надо убрать apps
		mqttData := mosquitto.MqttConf{
			ID:       hostname,
			Broker:   cfg.Broker,
			Username: cfg.Username,
			Password: cfg.Password,
			SubTopic: topicPrefix + hostname + "/commands/",
			PubTopic: topicPrefix,
			Icon:     &tray.Conn,
		}
		go mosquitto.StartBroker(mqttData)
	}

	if cfg.StatusOn {
		go func() {
			web.UpdateStatues() // первый запуск сразу

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				web.UpdateStatues()
			}
		}()
	}

	web.StartServer(cfg.Port)
}
