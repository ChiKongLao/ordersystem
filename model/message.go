package model

import "strings"


type Message struct {
	// 数据内容载体
	Payload  string
	// 发送的设备号
	Token    string
	// 发送到平台
	Platform string
	Topic string
	Qos int
}

func (msg *Message) GetTokens() []string {
	return strings.Split(msg.Token,",")
}

// 创建mqtt消息
func NewMqttMessage(topic string, payload string) *Message{
	return &Message{
		Topic:topic,
		Payload:payload,
		Qos:1,
	}

}