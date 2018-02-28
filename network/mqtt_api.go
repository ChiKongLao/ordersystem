package network

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"github.com/chikong/ordersystem/model"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

type MqttApi interface {
	Publish(msg MQTT.Message)
	Subscribe(topic string)
	RegisterCallback(callback MqttCallback)
}

type mqttApi struct {
	mqttClient MQTT.Client
	// 回调列
	callbackList []MqttCallback
}

const project  = "orderSystem"

// 发送聊天消息
func SendChatMessage(content string, user *model.User, businessId, tableId int){
	topic := fmt.Sprintf("%s/%v/%v/chat",project,businessId,tableId)
	data, _ := json.Marshal(model.NewChatMsg(user,content))
	payload := string(data)
	logrus.Debugf("发送mqtt聊天消息: %s , %s",topic,payload)
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload))

}

// 发送购物车变化消息
func SendShoppingCartMessage(businessId, tableId int){
	topic := fmt.Sprintf("%s/%v/%v/shopping_cart",project,businessId,tableId)
	payload := "update shopping cart"
	logrus.Debugf("发送mqtt购物车变化消息: %s , %s",topic,payload)
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload))

}
