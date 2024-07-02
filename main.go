package main

import (
	"github.com/tiredsosha/warden/mosquitto"
	config "github.com/tiredsosha/warden/tools/configurator"
	"github.com/tiredsosha/warden/tools/logger"
	"github.com/tiredsosha/warden/tray"
)

const (
	version = "0.0.1.beta"
)

// func getHostname() (hostname string) {
// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		logger.Warn.Printf("can't get hostname - %s\n", err)
// 		hostname = "default"
// 	}
// 	return
// }

func main() {

	// Логгер нужен, в нем надо разобраться в самом конце
	logger.LogInit(true)

	// инициализация конфига
	cfg := config.ConfInit()

	// хостнейм нам не нужен, мы будем использовать всегда только admin
	hostname := "admin"
	topicPrefix := "warden/"

	// так на это пока надо забить, логгер буду делать в конце
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

	// дальше нам надо инициализировать http сервак на 8080 порту.

}

// func main() {
// 	// Это отвечает за то, чтобы получать аргументы из командной строки. Нам не надо. Можно удалять
// 	cfg, state, debug := config.ArgsInit()

// 	// Логгер нужен, в нем надо разобраться в самом конце
// 	logger.LogInit(debug)

// 	// инициализация конфига, надо ее делать всегда, потому что у нас не будет cli
// 	if !state {
// 		cfg = config.ConfInit()
// 	}

// 	// хостнейм нам не нужен, мы будем использовать всегда только admin
// 	hostname := getHostname()
// 	topicPrefix := "warden/" + hostname + "/"

// 	// так на это пока надо забить, логгер буду делать в конце
// 	logger.DebugLog(version, debug, state, hostname, cfg.Broker, cfg.Username, cfg.Password, cfg.Apps)

// 	// вот это нужно, но мы тоже будем с ним разбираться в конце
// 	go tray.TrayStart()

// 	// так sub топик нам не нужен вообще, нам нужен только пуб топик. Мы наш стстум отправлять не будет. Так же надо убрать apps
// 	mqttData := mosquitto.MqttConf{
// 		ID:       hostname,
// 		Broker:   cfg.Broker,
// 		Username: cfg.Username,
// 		Password: cfg.Password,
// 		SubTopic: topicPrefix + "commands/",
// 		PubTopic: topicPrefix + "status/",
// 		Icon:     &tray.Conn,
// 		Apps:     cfg.Apps,
// 	}

// 	mosquitto.StartBroker(mqttData)
// }
