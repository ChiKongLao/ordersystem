package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 餐桌
type TableController struct {
	Ctx iris.Context
	services.TableService
	UserService services.UserService

}

// 获取餐桌
func (c *TableController) GetBy(businessId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, _, err = c.UserService.CheckRoleIsManagerOrBusinessWithToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	tableStatus,_ := c.Ctx.PostValueInt(constant.NameStatus)

	var list []model.TableInfo
	status, list, err = c.GetTableList(businessId,tableStatus)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}

// 获取餐桌
func (c *TableController) GetByStatusBy(businessId,tableStatus int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var list []model.TableInfo
	status, list, err = c.GetTableList(businessId,tableStatus)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
	}
}

// 获取餐桌详情
func (c *TableController) GetByBy(businessId, tableId int) (int,interface{}) {

	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.TableInfo
	status, item, err = c.GetTable(businessId,tableId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 添加餐桌
func (c *TableController) PostBy(businessId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, _, err := c.UserService.CheckRoleIsManagerOrBusinessWithToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	name := c.Ctx.FormValue(constant.Name)
	capacity, _ := c.Ctx.PostValueInt(constant.NameCapacity)

	status, err = c.InsertTable(&model.TableInfo{
		BusinessId:businessId,
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
func (c *TableController) PutByStatusBy(businessId,tableId int) (int,interface{}) {

	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	tableStatus, _ := c.Ctx.PostValueInt(constant.NameStatus)

	status, err = c.TableService.UpdateTableStatus(businessId,userId,tableId,tableStatus)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 修改餐桌
func (c *TableController) PutByBy(businessId, tableId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}


	status, _, err := c.UserService.CheckRoleIsManagerOrBusinessWithToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	name := c.Ctx.FormValue(constant.Name)
	tableStatus, _ := c.Ctx.PostValueInt(constant.NameStatus)
	capacity, _ := c.Ctx.PostValueInt(constant.NameCapacity)

	status, err = c.UpdateTable(&model.TableInfo{
		Id:tableId,
		BusinessId:businessId,
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


// 换桌
func (c *TableController) PutByChange(businessId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}
	
	status, _, err := c.UserService.CheckRoleIsManagerOrBusinessWithToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	oldTableId, _ := c.Ctx.PostValueInt(constant.NameOldTableId)
	newTableId, _ := c.Ctx.PostValueInt(constant.NameNewTableId)

	status, err = c.ChangeTable(businessId, oldTableId, newTableId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 删除餐桌
func (c *TableController) DeleteByBy(businessId, tableId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}


	status, _, err := c.UserService.CheckRoleIsManagerOrBusinessWithToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, err = c.DeleteTable(businessId,tableId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


