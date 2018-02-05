package services

import (
	"github.com/chikong/ordersystem/datamodels"
)


type MessageService interface {
	PostBy(msg datamodels.Message) (bool ,error)
}


func NewMessageService() MessageService {
	return &messageService{}
}

type messageService struct {
}

/** 接收申请发送的消息 */
func (s *messageService) PostBy(msg datamodels.Message) (bool,error) {
	return true ,nil

}
