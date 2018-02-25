package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"encoding/json"
)

// 订单
type OrderController struct {
	Ctx         iris.Context
	services.OrderService
	UserService services.UserService
}

// 获取订单,
func (c *OrderController) GetBy(userId int) (int, interface{}) {

	status, user, err := c.UserService.GetUserFormToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	if user.Role != constant.RoleManager && user.Role != constant.RoleBusiness {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	status, _, err = c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	var item *model.OrderListResponse
	status, item, err = c.GetOrderList(userId,0,user.Role)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameData:        item,
	}

}

// 获取订单
func (c *OrderController) GetByTableBy(userId, tableId int) (int, interface{}) {
	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	var item *model.OrderListResponse
	status, item, err = c.GetOrderList(userId,tableId,constant.RoleCustomer)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	return status, iris.Map{
		constant.NameData:        item,
		constant.NameTotalCount:  len(item.List),
	}

}

// 获取订单详情
func (c *OrderController) GetByBy(userId, orderId int) (int, interface{}) {

	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.OrderResponse
	status, item, err = c.GetOrder(orderId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameData: item,
	}
}

// 获取商家的老用户列表
func (c *OrderController) GetCustomerBy(businessId int) (int, interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}
	status, user, err := c.UserService.GetUserById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	if !user.IsManagerOrBusiness() {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}
	status, userList, err := c.OrderService.GetOldCustomer(businessId)
	return status, iris.Map{
		constant.NameData: userList,
	}
}

// 添加订单, 商家不能自己下单,只能客户下
func (c *OrderController) PostBy(businessId int) (int, interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	tableId, _ := c.Ctx.PostValueInt(constant.NameTableId)
	personNum, _ := c.Ctx.PostValueInt(constant.NamePersonNum)
	var list = new([]model.Food)
	err = json.Unmarshal([]byte(c.Ctx.FormValue(constant.NameFood)), &list)

	if err != nil {
		return iris.StatusBadRequest, iris.Map{constant.NameMsg: "菜单格式错误"}
	}

	var orderId int

	status, orderId, err = c.InsertOrder(&model.Order{
		BusinessId: businessId,
		UserId:     userId,
		TableId:    tableId,
		PersonNum:  personNum,
		FoodList:   *list,
	})

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
		constant.NameID:   orderId,
	}
}

// 修改订单,只能商家操作
func (c *OrderController) PutByBy(userId, orderId int) (int, interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManagerOrBusiness() {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	tableId, _ := c.Ctx.PostValueInt(constant.NameTableId)
	orderStatus, _ := c.Ctx.PostValueInt(constant.NameStatus)
	personNum, _ := c.Ctx.PostValueInt(constant.NamePersonNum)

	status, err = c.UpdateOrder(&model.Order{
		Id:         orderId,
		BusinessId: userId,
		TableId:    tableId,
		Status:     orderStatus,
		PersonNum:  personNum,
	})

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 修改订单,只能商家操作
func (c *OrderController) PutByStatusBy(userId, orderId int) (int, interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManagerOrBusiness() {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	orderStatus, _ := c.Ctx.PostValueInt(constant.NameStatus)

	status, err = c.UpdateOrderStatus(orderId,orderStatus)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 删除订单
func (c *OrderController) DeleteByBy(userId, orderId int) (int, interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness() {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	status, err = c.DeleteOrder(userId, orderId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}
