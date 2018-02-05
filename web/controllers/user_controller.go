package controllers

import (
	"github.com/chikong/ordersystem/datamodels"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
)

type UserController struct {
	Ctx iris.Context
	services.UserService
}

// 注册
func (c *UserController) PostRegister() (int,interface{}) {
	userName := c.Ctx.FormValue(datamodels.UserName)
	password := c.Ctx.FormValue(datamodels.Password)

	//var status int
	//var err	error

	status,err := c.UserService.InsertUser(userName,password)

	if err != nil{
		return status,datamodels.NewErrorResponse(err)
	}

	return status,iris.Map{datamodels.KeyIsOk:true}
}


// 登录
func (c *UserController) PostLogin()(int,interface{}){
	userName := c.Ctx.FormValue(datamodels.UserName)
	password := c.Ctx.FormValue(datamodels.Password)

	status, token, err := c.UserService.Login(userName,password)
	if err != nil{
		return status,datamodels.NewErrorResponse(err)
	}

	return status,map[string]string{
		datamodels.Token:token,
	}

}

