package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"encoding/json"
)

// 订单
type OrderController struct {
	Ctx iris.Context
	services.OrderService
	UserService services.UserService

}

// 获取订单
func (c *OrderController) GetBy(userId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	var list []model.Order
	status, list, err = c.GetOrderList(userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	var personCount int
	var priceCount float32
	for _, subItem := range list{
		priceCount += subItem.Price
		personCount += subItem.PersonNum
	}

	return status,iris.Map{
		constant.NameTotalPerson:personCount,
		constant.NameTotalPrice:priceCount,
		constant.NameData:list,
		constant.NameTotalCount:len(list),
		}
}

// 获取订单详情
func (c *OrderController) GetByBy(userId, orderId int) (int,interface{}) {

	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.Order
	status, item, err = c.GetOrder(userId,orderId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 获取商家的老用户列表
func (c *OrderController) GetCustomerBy(businessId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx,businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}
	status, user, err := c.UserService.GetUserById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	if !user.IsManagerOrBusiness(){
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}
	status, userList, err := c.OrderService.GetOldCustomer(businessId)
	return status,iris.Map{
		constant.NameData:userList,
		}
}

// 添加订单, 商家不能自己下单,只能客户下
func (c *OrderController) PostBy(businessId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	tableName := c.Ctx.FormValue(constant.NameTableName)
	personNum, _ := strconv.Atoi(c.Ctx.FormValue(constant.NamePersonNum))
	var list = new([]model.Food)
	err = json.Unmarshal([]byte(c.Ctx.FormValue(constant.NameFood)),&list)

	if err != nil {
		return iris.StatusBadRequest,iris.Map{constant.NameMsg:"菜单格式错误"}
	}

	var orderId int

	status,orderId,err = c.InsertOrder(&model.Order{
		BusinessId:businessId,
		UserId:userId,
		TableName:tableName,
		PersonNum:personNum,
		FoodList:*list,
	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
			constant.NameID: orderId,
		}
}

// 修改订单,只能商家操作
func (c *OrderController) PutByBy(userId, orderId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManagerOrBusiness(){
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	tableName := c.Ctx.FormValue(constant.NameTableName)
	orderStatus, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameStatus))
	personNum, _ := strconv.Atoi(c.Ctx.FormValue(constant.NamePersonNum))

	status, err = c.UpdateOrder(&model.Order{
		Id:orderId,
		BusinessId:userId,
		TableName:tableName,
		Status:orderStatus,
		PersonNum:personNum,
	} )


	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 删除订单
func (c *OrderController) DeleteByBy(userId, orderId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	status, err = c.DeleteOrder(userId,orderId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

