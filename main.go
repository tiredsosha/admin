package main

import (
	"github.com/tiredsosha/admin/mosquitto"
	config "github.com/tiredsosha/admin/tools/configurator"
	"github.com/tiredsosha/admin/tools/logger"
	"github.com/tiredsosha/admin/tray"
	"github.com/tiredsosha/admin/web"
)

const (
	version = "0.0.1.beta"
)

func main() {
	logger.LogInit(true)

	cfg := config.ConfInit()

	hostname := "admin"
	topicPrefix := "warden/"

	logger.DebugLog(version, true, hostname, cfg.Broker, cfg.Username, cfg.Password, cfg.Port)

	// вот это нужно, но мы тоже будем с ним разбираться в конце
	go tray.TrayStart()

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

	web.StartServer(cfg.Port)

}
