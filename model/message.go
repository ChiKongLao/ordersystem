package model

import "strings"

type Message struct {
	// 数据内容载体
	Payload string
	// 发送的设备号
	Token string
	// 发送到平台
	Platform string
	Topic    string
	Qos      int
	Desc     string // 描述
}

func (msg *Message) GetTokens() []string {
	return strings.Split(msg.Token, ",")
}

// 创建mqtt消息
func NewMqttMessage(topic, payload,desc string) *Message {
	return &Message{
		Topic:   topic,
		Payload: payload,
		Desc:desc,
		Qos:     1,
	}

}
