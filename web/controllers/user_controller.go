package controllers

import (
	"github.com/chikong/ordersystem/datamodels"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	services.UserService
}

// 注册
func (c *UserController) PostRegister() (int,interface{}) {
	userName := c.Ctx.FormValue(datamodels.NameUserName)
	password := c.Ctx.FormValue(datamodels.NamePassword)
	nickName := c.Ctx.FormValue(datamodels.NameNickName)
	role,_ := strconv.Atoi(c.Ctx.FormValue(datamodels.NameRole))

	status,err := c.UserService.InsertUser(role,userName,password,nickName)

	if err != nil{
		return status,datamodels.NewErrorResponse(err)
	}

	return status,iris.Map{datamodels.KeyIsOk:true}
}


// 登录
func (c *UserController) PostLogin()(int,interface{}){
	userName := c.Ctx.FormValue(datamodels.NameUserName)
	password := c.Ctx.FormValue(datamodels.NamePassword)

	status, token, err := c.UserService.Login(c.Ctx,userName,password)
	if err != nil{
		return status,datamodels.NewErrorResponse(err)
	}

	return status,map[string]string{
		datamodels.NameAuthorization:token,
	}

}

// 获取用户信息
func (c *UserController) Get()(int,interface{}){
	userName,err := authentication.GetUserNameFormHeaderToken(c.Ctx)
	if err != nil{
		return iris.StatusInternalServerError,datamodels.NewErrorResponse(err)
	}
	user, err := c.UserService.GetUserByName(userName)
	if err != nil{
		return iris.StatusInternalServerError,datamodels.NewErrorResponse(err)
	}

	return iris.StatusOK,iris.Map{
		datamodels.NameID:user.Id,
		datamodels.NameNickName:user.NickName,
	}

}

