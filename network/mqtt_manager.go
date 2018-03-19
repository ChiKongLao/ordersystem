package network

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"strconv"
	"time"
	"github.com/sirupsen/logrus"
	"sync"
	"github.com/chikong/ordersystem/model"
	"github.com/chikong/ordersystem/manager"
	"github.com/garyburd/redigo/redis"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/configs"
)


type MqttManager interface {
	Publish(msg *model.Message)
	Subscribe(topic string)
	RegisterCallback(func(message MQTT.Message))
}

type mqttManager struct {
	mqttClient MQTT.Client
	// 回调列
	callbackList []func(message MQTT.Message)
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

// 订阅通用业务topic
func (m *mqttManager) subscribeCommon(){
	m.Subscribe(MqttProject+"/+/+"+TopicChat)
}


// 连接mqtt
func (m *mqttManager) initClient() {
	if !configs.GetConfig().Mqtt.Open {
		return
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(MqttUrl)
	opts.SetConnectTimeout(5*time.Second)
	opts.SetAutoReconnect(true)
	opts.SetClientID(fmt.Sprintf("system_%s", strconv.FormatInt(time.Now().Unix(), 10)))
	//opts.SetUsername(*user)
	//opts.SetPassword(*password)
	//opts.SetCleanSession(*cleansess)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		key := fmt.Sprintf("%s_%s",msg.Topic(),string(msg.Payload()))
		exist, _ := redis.Bool(manager.GetRedisConn().Do(manager.RedisExists,key))
		if exist {
			logrus.Warnf("重复消息. topic= %s", msg.Topic())
			return
		}
		manager.GetRedisConn().Do(manager.RedisSet,key,"",manager.RedisEx,constant.TimeCacheMsgDuplicate)
		for _, value := range m.callbackList {
			value(msg)
		}
	})

	client := MQTT.NewClient(opts)
	logrus.Infof("连接mqtt. %s", MqttUrl)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("mqtt连接失败: %s",token.Error())
		return
	}
	logrus.Info("连接mqtt成功")
	m.mqttClient = client
	m.subscribeCommon()
}

// 发送消息
func (m *mqttManager) Publish(msg *model.Message) {
	if m.mqttClient == nil || !m.mqttClient.IsConnected() {
		logrus.Warn("发送消息失败.mqtt未连接")
		go m.initClient()
		return
	}
	m.mqttClient.Publish(msg.Topic, byte(msg.Qos), false, msg.Payload)
	logrus.Infof("发送 %s 消息: %s, %s",msg.Desc,msg.Topic,msg.Payload)
}

// 订阅topic
func (m *mqttManager) Subscribe(topic string) {
	if m.mqttClient == nil || !m.mqttClient.IsConnected() {
		logrus.Warn("发送消息失败.mqtt未连接")
		go m.initClient()
		return
	}
	if token := m.mqttClient.Subscribe(topic, byte(1), nil); token.Wait() && token.Error() != nil {
		logrus.Infof("订阅topic失败. %s", topic)
	}
	logrus.Infof("订阅topic成功. %s", topic)
}

func (m *mqttManager) RegisterCallback(fun func(message MQTT.Message)) {
	m.callbackList = append(m.callbackList, fun)
}
