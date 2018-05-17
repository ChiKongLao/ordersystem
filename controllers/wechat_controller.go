package controllers

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/sirupsen/logrus"
	"github.com/chikong/ordersystem/model"
)

// 订单
type WeChatController struct {
	Ctx iris.Context
	services.WechatService
}

// 获取用户的个人信息
func (c *WeChatController) GetAuth() (int, interface{}) {
	return iris.StatusOK, iris.Map{
		constant.NameUrl: c.GetAuthUrl(),
	}

}

// 获取用户的个人信息
func (c *WeChatController) GetAuthResponse() (int, interface{}) {

	data, err := c.GetUserInfo(c.Ctx.URLParam(constant.NameCode), c.Ctx.URLParam(constant.NameState))
	if err != nil {
		return iris.StatusInternalServerError, model.NewErrorResponse(err)
	}
	return iris.StatusOK, iris.Map{
		constant.NameData: data,
	}

}

// 获取微信支付回调
func (c *WeChatController) GetBy(orderId int) (int, interface{}) {
	//logrus.Infoln("接收到微信支付回调")
	//status, err := c.UpdateOrderStatus(orderId, constant.OrderStatusPaid)
	//if err != nil {
	//	return status, model.NewErrorResponse(err)
	//}
	//return status, iris.Map{
	//	constant.NameIsOk: iris.StatusOK,
	//}

	return iris.StatusOK, nil

}

// 微信登录回调,
func (c *WeChatController) GetLogin() (int, interface{}) {
	logrus.Debugln("接收到微信登录回调: ", )
	return iris.StatusOK, c.Ctx.URLParam("echostr")

}
