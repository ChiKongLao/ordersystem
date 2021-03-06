package model

import (
	"github.com/chikong/ordersystem/util"
)

//　聊天具体
type ChatDetail struct {
	Content  string `json:"content"`
	Head     string `json:"head"`
	UserId   int    `json:"userId"`
	NickName string `json:"nickName"`
	Time     int64 `json:"time"`
	MsgId    string `json:"msgId"`
}

//　聊天消息
type ChatMsg struct {
	Data     ChatDetail `json:"data"`
	DataType string     `json:"dataType"`
	Action   string     `json:"action"`
}

// 创建mqtt消息
func NewChatMsg(user *User, content string) *ChatMsg {
	return &ChatMsg{
		Data: ChatDetail{Content: content,
			Head: user.Head,
			UserId: user.Id,
			NickName: user.NickName,
			Time: util.GetCurrentTime(),
			MsgId: util.GetUUID(),
			},

		DataType: "text",
		Action:   "say",
	}

}
