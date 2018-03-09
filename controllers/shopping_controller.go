package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 购物车
type ShoppingController struct {
	Ctx iris.Context
	services.ShoppingService
	UserService services.UserService

}

// 获取购物车
func (c *ShoppingController) GetByTableBy(businessId,tableId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, item, err := c.GetShopping(businessId,userId,tableId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
	}
}

// 修改购物车
func (c *ShoppingController) PutBy(businessId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	foodType := c.Ctx.FormValue(constant.NameType)
	foodId, _ := c.Ctx.PostValueInt(constant.NameFoodId)
	tableId, _ := c.Ctx.PostValueInt(constant.NameTableId)
	num, _ := c.Ctx.PostValueInt(constant.NameNum)
	status,err = c.UpdateShopping(foodType,userId,businessId,foodId,num,tableId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
			}
}

