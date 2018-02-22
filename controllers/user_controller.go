package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/chikong/ordersystem/constant"
)

type UserController struct {
	Ctx iris.Context
	services.UserService
}

// 注册
func (c *UserController) PostRegister() (int,interface{}) {
	userName := c.Ctx.FormValue(constant.NameUserName)
	password := c.Ctx.FormValue(constant.NamePassword)
	nickName := c.Ctx.FormValue(constant.NameNickName)
	head := c.Ctx.FormValue(constant.NameHead)
	role,_ := c.Ctx.PostValueIntDefault(constant.NameRole,3)

	status,err := c.UserService.InsertUser(role,userName,password,nickName,head)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{constant.NameIsOk:true}
}


// 登录
func (c *UserController) PostLogin()(int,interface{}){
	userName := c.Ctx.FormValue(constant.NameUserName)
	password := c.Ctx.FormValue(constant.NamePassword)

	status, token, err := c.UserService.Login(userName,password)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,map[string]string{
		constant.NameAuthorization:token,
	}

}

// 获取用户信息
func (c *UserController) Get()(int,interface{}){
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}
	status, user, err := c.UserService.GetUserById(userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameID:user.Id,
		constant.NameNickName:user.NickName,
		constant.NameHead:user.Head,

	}

}

