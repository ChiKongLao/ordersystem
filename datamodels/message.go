package datamodels

import "strings"


type Message struct {
	// 数据内容载体
	Payload  string
	// 发送的设备号
	Token    string
	// 发送到平台
	Platform string
}

func (msg *Message) GetTokens() []string {
	return strings.Split(msg.Token,",")
}