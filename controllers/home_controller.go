package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/chikong/ordersystem/constant"
)

// 首页
type HomeController struct {
	Ctx iris.Context
	services.HomeService
	UserService services.UserService
}

// 获取首页,自动识别为用户身份,再获取不同数据
func (c *HomeController) GetBy(businessId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	status,userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	status, user,err := c.UserService.GetUserById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var data interface{}
	if user.IsManagerOrBusiness() {
		status, data, err = c.GetBusinessHome(businessId)
	}else{
		status, data, err = c.GetCustomerHome(businessId,userId)
	}
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:data,
	}
}


