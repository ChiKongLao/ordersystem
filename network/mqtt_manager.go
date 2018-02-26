package network

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"strconv"
	"time"
	"github.com/sirupsen/logrus"
	"sync"
)

type MqttCallback func(message MQTT.Message)

type MqttManager interface {
	Publish(msg MQTT.Message)
	Subscribe(topic string)
	RegisterCallback(callback MqttCallback)
}

type mqttManager struct {
	mqttClient MQTT.Client
	// 回调列
	callbackList []MqttCallback
}

var instance *mqttManager
var once sync.Once


func GetMqttInstance() MqttManager {
	once.Do(func() {
		instance = &mqttManager{}
		instance.initClient()
	})

	return instance
}



// 连接mqtt
func (m *mqttManager) initClient() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MqttUrl)
	opts.SetClientID(fmt.Sprintf("system_%s", strconv.FormatInt(time.Now().Unix(), 10)))
	//opts.SetUsername(*user)
	//opts.SetPassword(*password)
	//opts.SetCleanSession(*cleansess)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		logrus.Info(string(msg.Payload()))
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("mqtt连接失败")
		panic(token.Error())
	}
	logrus.Infof("连接mqtt成功. %s", MqttUrl)
	m.mqttClient = client

}

// 发送消息
func (m *mqttManager) Publish(msg MQTT.Message) {
	if m.mqttClient == nil || !m.mqttClient.IsConnected() {
		go m.initClient()
		return
	}
	m.mqttClient.Publish(msg.Topic(), msg.Qos(), false, msg.Payload())
}

// 订阅topic
func (m *mqttManager) Subscribe(topic string) {
	if m.mqttClient == nil || !m.mqttClient.IsConnected() {
		go m.initClient()
		return
	}
	if token := m.mqttClient.Subscribe(topic, byte(1), nil); token.Wait() && token.Error() != nil {
		logrus.Infof("订阅topic失败. %s", topic)
	}
}

func (m *mqttManager) RegisterCallback(callback MqttCallback) {
	m.callbackList = append(m.callbackList, callback)
}
