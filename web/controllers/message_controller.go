package controllers

import (
	"github.com/chikong/ordersystem/datamodels"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
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
		token := c.Ctx.Values().Get(authentication.JWTHandler.Config.ContextKey).(*jwt.Token)
		logrus.Infof("token = %v",token.Raw)
		return status,iris.Map{datamodels.KeyIsOk:true}
	} else{
		return status,&datamodels.Response{err.Error()}
	}
	return iris.StatusOK,""
}

