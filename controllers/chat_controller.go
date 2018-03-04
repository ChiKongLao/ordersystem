package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
)

// 聊天室
type ChatController struct {
	Ctx iris.Context
	services.ChatService
	UserService services.UserService

}

// 获取聊天记录
func (c *ChatController) GetByTableBy(businessId,tableId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var list []model.ChatMsg
	status, list, err =	c.GetChatLog(businessId,tableId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}



