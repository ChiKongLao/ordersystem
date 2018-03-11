package network

import (
	"fmt"
	"github.com/chikong/ordersystem/model"
	"encoding/json"
	"github.com/chikong/ordersystem/util"
)

type MqttApi interface {

}

type mqttApi struct {
}

const(
	MqttProject  = "orderSystem"
	TopicChat = "/chat"
	TopicShoppingCart = "/shoppingCart"
	TopicOrder = "/order"
)

// 发送聊天消息
func SendChatMessage(content string, user *model.User, businessId, tableId int){
	topic := fmt.Sprintf(MqttProject + "/%v/%v" + TopicChat,businessId,tableId)
	payload, _ := json.Marshal(model.NewChatMsg(user,content))
	GetMqttInstance().Publish(model.NewMqttMessage(topic,string(payload),"聊天"))

}

// 发送购物车变化消息
func SendShoppingCartMessage(businessId, tableId int){
	topic := fmt.Sprintf(MqttProject + "/%v/%v" + TopicShoppingCart,businessId,tableId)
	payload, _ := json.Marshal(struct {
		Data string `json:"data"`
		Time int64	`json:"time"`
	}{
		Data:"update shopping cart",
		Time:util.GetCurrentTime(),
	})
	GetMqttInstance().Publish(model.NewMqttMessage(topic,string(payload),"购物车变化"))

}

// 发送订单状态变化消息
func SendOrderMessage(businessId int, order *model.Order){
	topic := fmt.Sprintf(MqttProject + "/%v/%v" + TopicOrder,businessId,order.Id)
	payload,_ := json.Marshal(struct {
		Data interface{} `json:"data"`
		}{order})
	GetMqttInstance().Publish(model.NewMqttMessage(topic,string(payload),"订单状态变化"))

}
