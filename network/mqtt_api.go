package network

import (
	"fmt"
	"github.com/chikong/ordersystem/model"
	"encoding/json"
)

type MqttApi interface {

}

type mqttApi struct {
}

const(
	MqttProject  = "orderSystem"
	TopicChat = "/chat"
	TopicShoppingCart = MqttProject + "/shoppingCart"
	TopicOrder = MqttProject + "/order"
)

// 发送聊天消息
func SendChatMessage(content string, user *model.User, businessId, tableId int){
	topic := fmt.Sprintf(MqttProject + "/%v/%v" + TopicChat,businessId,tableId)
	data, _ := json.Marshal(model.NewChatMsg(user,content))
	payload := string(data)
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload,"聊天"))

}

// 发送购物车变化消息
func SendShoppingCartMessage(businessId, tableId int){
	topic := fmt.Sprintf(MqttProject + "/%v/%v" + TopicShoppingCart,businessId,tableId)
	payload := "update shopping cart"
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload,"购物车变化"))

}

// 发送订单状态变化消息
func SendOrderMessage(businessId, orderId,status int){
	topic := fmt.Sprintf(MqttProject + "/%v/%v" + TopicOrder,businessId,orderId)
	payload := fmt.Sprintf(`"status":%v`,status)
	GetMqttInstance().Publish(model.NewMqttMessage(topic,payload,"订单状态变化"))

}
