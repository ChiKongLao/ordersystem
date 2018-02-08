package services

import (
	"github.com/chikong/ordersystem/model"
)


type MessageService interface {
	PostBy(msg model.Message) (bool ,error)
}


func NewMessageService() MessageService {
	return &messageService{}
}

type messageService struct {
}

/** 接收申请发送的消息 */
func (s *messageService) PostBy(msg model.Message) (bool,error) {
	return true ,nil

}
