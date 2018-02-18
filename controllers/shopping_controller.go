package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 购物车
type ShoppingController struct {
	Ctx iris.Context
	services.ShoppingService
	UserService services.UserService

}

// 获取购物车
func (c *ShoppingController) GetBy(businessId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, item, err := c.GetShopping(businessId,userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,item
}

// 修改购物车,只能商家操作
func (c *ShoppingController) PutBy(businessId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	dashesType := c.Ctx.FormValue(constant.NameType)
	dashesId, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameDashesId))
	num, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameNum))

	status,err = c.UpdateShopping(userId,businessId, dashesId,num,dashesType)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
			}
}

//
//// 删除购物车
//func (c *ShoppingController) DeleteByBy(userId, orderId int) (int,interface{}) {
//	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
//	if !isOwn {
//		return iris.StatusUnauthorized, model.NewErrorResponse(err)
//	}
//
//	status, user, err := c.UserService.GetUserById(userId)
//
//	if err != nil {
//		return status, model.NewErrorResponse(err)
//	}
//
//	if !user.IsManager() && !user.IsBusiness(){
//		return iris.StatusUnauthorized,errors.New("没有该权限")
//	}
//
//	status, err = c.DeleteShopping(userId,orderId)
//
//	if err != nil{
//		return status, model.NewErrorResponse(err)
//	}
//
//	return status,iris.Map{
//			constant.NameIsOk:true,
//		}
//}


