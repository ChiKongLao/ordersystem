package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
)

// 订单
type WeChatNotifyController struct {
	Ctx         iris.Context
	services.OrderService
	UserService services.UserService
}

// 获取订单,
func (c *WeChatNotifyController) GetBy(orderId int) (int, interface{}) {

	status, err := c.UpdateOrderStatus(orderId, constant.OrderStatusPaid)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	return status, iris.Map{
		constant.NameIsOk: iris.StatusOK,
	}

}
