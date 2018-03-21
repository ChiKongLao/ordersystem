package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 订单
type OrderController struct {
	Ctx         iris.Context
	services.OrderService
	UserService services.UserService
}

// 获取订单,
func (c *OrderController) GetBy(businessId int) (int, interface{}) {

	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}
	status, user, err := c.UserService.GetUserFormToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, _, err = c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	orderStatus, _ := c.Ctx.URLParamIntDefault(constant.NameStatus, constant.OrderStatusAll)

	var item *model.OrderListResponse
	status, item, err = c.GetOrderList(businessId, 0, user.Role, orderStatus)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameData: item,
	}

}

// 获取订单
func (c *OrderController) GetByTableBy(businessId, tableId int) (int, interface{}) {
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	orderStatus, _ := c.Ctx.URLParamInt(constant.NameStatus)
	status, user, err := c.UserService.GetUserFormToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	var item *model.OrderListResponse
	status, item, err = c.GetOrderList(businessId, tableId, user.Role, orderStatus)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if orderStatus == constant.OrderStatusWaitPay { // 转换未支付的格式, 以后改成为已支付一样
		if len(item.List) != 0 {
			tmpOrder := item.List[0]
			var count int
			for _, foodItem := range tmpOrder.FoodList {
				count += foodItem.Num
			}
			tmpOrder.FoodCount = count
			return status, iris.Map{
				constant.NameData: tmpOrder,
			}
		}

	}

	return status, iris.Map{
		constant.NameData: item,
	}

}

// 获取订单详情
func (c *OrderController) GetByBy(businessId, orderId int) (int, interface{}) {

	status, _, err := c.UserService.GetBusinessById(businessId)
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
	shoppingCartId, _ := c.Ctx.PostValueInt(constant.NameShoppingCartId)
	personNum, _ := c.Ctx.PostValueInt(constant.NamePersonNum)
	var orderId int

	status, orderId, err = c.InsertOrder(&model.Order{
		BusinessId: businessId,
		UserId:     userId,
		TableId:    tableId,
		PersonNum:  personNum,
	}, shoppingCartId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
		constant.NameID:   orderId,
	}
}

// 修改订单,只能商家操作
func (c *OrderController) PutByBy(businessId, orderId int) (int, interface{}) {
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

	tableId, _ := c.Ctx.PostValueInt(constant.NameTableId)
	orderStatus, _ := c.Ctx.PostValueInt(constant.NameStatus)
	personNum, _ := c.Ctx.PostValueInt(constant.NamePersonNum)

	status, err = c.UpdateOrder(&model.Order{
		Id:         orderId,
		BusinessId: businessId,
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
func (c *OrderController) PutByStatusBy(businessId, orderId int) (int, interface{}) {
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

	orderStatus, _ := c.Ctx.PostValueInt(constant.NameStatus)

	status, err = c.UpdateOrderStatus(orderId, orderStatus)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 删除订单
func (c *OrderController) DeleteByBy(businessId, orderId int) (int, interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(businessId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness() {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	status, err = c.DeleteOrder(businessId, orderId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}
