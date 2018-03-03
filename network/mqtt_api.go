package network

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"github.com/chikong/ordersystem/model"
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

const(
	project  = "orderSystem"
	TopicChat = project + "/%v/%v/chat"
	TopicShoppingCart = project + "/%v/%v/shoppingCart"
	TopicOrder = project + "/%v/%v/order"
)

// 发送聊天消息
func SendChatMessage(content string, user *model.User, businessId, tableId int){
	topic := fmt.Sprintf(TopicChat,businessId,tableId)
	data, _ := json.Marshal(model.NewChatMsg(user,content))
	payload := string(data)
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload,"聊天"))

}

// 发送购物车变化消息
func SendShoppingCartMessage(businessId, tableId int){
	topic := fmt.Sprintf(TopicShoppingCart,businessId,tableId)
	payload := "update shopping cart"
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload,"购物车变化"))

}

// 发送订单状态变化消息
func SendOrderMessage(businessId, orderId,status int){
	topic := fmt.Sprintf(TopicOrder,businessId,orderId)
	payload := fmt.Sprintf(`"status":%v`,status)
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload,"订单状态变化"))

}
