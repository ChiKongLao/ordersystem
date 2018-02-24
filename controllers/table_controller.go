package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 餐桌
type TableController struct {
	Ctx iris.Context
	services.TableService
	UserService services.UserService

}

// 获取餐桌
func (c *TableController) GetBy(userId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	var list []model.TableInfo
	status, list, err = c.GetTableList(userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}

// 获取餐桌详情
func (c *TableController) GetByBy(userId, tableId int) (int,interface{}) {

	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.TableInfo
	status, item, err = c.GetTable(userId,tableId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 添加餐桌
func (c *TableController) PostBy(userId int) (int,interface{}) {
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

	name := c.Ctx.FormValue(constant.Name)
	capacity, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameCapacity))

	status, err = c.InsertTable(&model.TableInfo{
		BusinessId:userId,
		Name:name,
		Capacity:capacity,
		Status:constant.TableStatusEmpty,

	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 使用餐桌
func (c *TableController) PostByJoinBy(businessId,tableId int) (int,interface{}) {
	status, err := c.TableService.JoinTable(businessId,tableId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 修改餐桌
func (c *TableController) PutByBy(userId, tableId int) (int,interface{}) {
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

	name := c.Ctx.FormValue(constant.Name)
	tableStatus, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameStatus))
	capacity, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameCapacity))

	status, err = c.UpdateTable(&model.TableInfo{
		Id:tableId,
		BusinessId:userId,
		Name:name,
		Status:tableStatus,
		Capacity:capacity,


	} )


	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 删除餐桌
func (c *TableController) DeleteByBy(userId, tableId int) (int,interface{}) {
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

	status, err = c.DeleteTable(userId,tableId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


