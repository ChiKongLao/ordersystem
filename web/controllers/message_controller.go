package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/constant"
)

type MessageController struct {
	Ctx iris.Context
}

// 接收需要推送的信息
func (c *MessageController) Post() (int,interface{}) {
	//platform := c.Ctx.FormValue(datamodels.NamePlatForm)
	//msg := datamodels.Message{
	//	Payload: c.Ctx.FormValue(datamodels.NamePayload),
	//	Token: c.Ctx.FormValue(datamodels.NameToken),
	//	Platform: platform}

	var status = iris.StatusOK
	var err	error

	if  status == iris.StatusOK{
		return status,iris.Map{constant.NameIsOk:true}
	} else{
		return status, model.NewErrorResponse(err)
	}
	return iris.StatusOK,""
}

