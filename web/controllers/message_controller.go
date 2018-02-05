package controllers

import (
	"github.com/chikong/ordersystem/datamodels"
	"github.com/kataras/iris"
)

type MessageController struct {
}

// 接收需要推送的信息
func (c *MessageController) Post() (int,interface{}) {
	//platform := c.Ctx.FormValue(datamodels.NamePlatForm)
	//msg := datamodels.Message{
	//	Payload: c.Ctx.FormValue(datamodels.NamePayload),
	//	Token: c.Ctx.FormValue(datamodels.NameToken),
	//	Platform: platform}

	var status int
	var err	error

	if  status == iris.StatusOK{
		return status,iris.Map{datamodels.KeyIsOk:true}
	} else{
		return status,&datamodels.Response{err.Error()}
	}
	return iris.StatusOK,""
}

