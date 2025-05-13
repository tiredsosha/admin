package mosquitto

import (
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tiredsosha/admin/tools/logger"
)

const (
	port          = 1883
	keyLifeTime   = 10 // minute
	reconnTime    = 20 // sec
	pubVolumeTime = 5  // sec
	pubMuteTime   = 7  // sec
)

var conn mqtt.Client

type MqttConf struct {
	ID       string
	Broker   string
	Username string
	Password string
	SubTopic string
	PubTopic string
	Icon     *bool
}

func (data *MqttConf) messageHandler(client mqtt.Client, msgHand mqtt.Message) {
	topic := msgHand.Topic()
	msg := strings.TrimSpace(string(msgHand.Payload()))
	logger.Debug.Printf("got %s on %q\n", msg, topic)
}

func (data *MqttConf) connectHandler(client mqtt.Client) {
	topic := data.SubTopic + "#"
	client.Unsubscribe(topic)

	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	logger.Debug.Printf("subscribed to %q\n", topic)
	logger.Debug.Printf("publishing to topic prefix %q\n", data.PubTopic+"#")
	logger.Info.Println("connection to mqtt broker is successful")
	tokenPub := client.Publish(data.PubTopic+data.ID+"/online", 0, true, "true")
	tokenPub.Wait()
	*data.Icon = true
}

func (data *MqttConf) lostHandler(client mqtt.Client, err error) {
	logger.Warn.Printf("mqtt: connection to mqtt broker is lost - %s\n", err)
	*data.Icon = false
}

func (data *MqttConf) SendMqtt(hostname, topic, command string) {
	finalTopic := data.PubTopic + hostname + "/commands/" + topic
	token := conn.Publish(finalTopic, 0, false, command)
	token.Wait()
}

func StartBroker(data MqttConf) {
	messagePubHandler := data.messageHandler
	connectHandler := data.connectHandler
	connectLostHandler := data.lostHandler

	// MQTT INIT //
	mqttHandler := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", data.Broker, port)).
		SetClientID(data.ID).
		SetUsername(data.Username).
		SetPassword(data.Password).
		SetAutoReconnect(true).
		SetDefaultPublishHandler(messagePubHandler).
		SetConnectionLostHandler(connectLostHandler).
		SetOnConnectHandler(connectHandler).
		SetKeepAlive(keyLifeTime*time.Minute).
		SetWill(data.PubTopic+data.ID+"/online", "false", 2, true)

	//conn := mqtt.NewClient(mqttHandler)
	conn = mqtt.NewClient(mqttHandler)

	for {
		status := true
		if token := conn.Connect(); token.Wait() && token.Error() != nil {
			logger.Warn.Printf("mqtt: can't connect to mqtt broker - %s\n", token.Error())
			status = false
		}

		if status {
			break
		}
		time.Sleep(reconnTime * time.Second)
	}
}

func (data *MqttConf) publisher(topic, command string) {
	token := conn.Publish(data.PubTopic+topic, 0, false, command)
	token.Wait()
}
