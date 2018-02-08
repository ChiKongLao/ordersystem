package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"errors"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 餐桌
type TableController struct {
	Ctx iris.Context
	services.TableService
	UserService services.UserService

}

// 获取餐桌
func (c *TableController) GetBy(userId string) (int,interface{}) {
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
		constant.NameSize:len(list),
		}
}

// 获取餐桌详情
func (c *TableController) GetByBy(userId, tableId string) (int,interface{}) {

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
func (c *TableController) PostBy(userId string) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized,errors.New("没有该权限")
	}

	name := c.Ctx.FormValue(constant.Name)
	capacity, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameCapacity))

	userIdInt,_ := strconv.Atoi(userId)

	status, err = c.InsertTable(&model.TableInfo{
		BusinessId:userIdInt,
		Name:name,
		Capacity:capacity,

	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 修改餐桌
func (c *TableController) PutByBy(userId, tableId string) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized,errors.New("没有该权限")
	}

	name := c.Ctx.FormValue(constant.Name)
	tableStatus, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameStatus))
	capacity, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameCapacity))

	userIdInt,_ := strconv.Atoi(userId)
	tableIdInt,_ := strconv.Atoi(tableId)

	status, err = c.UpdateTable(&model.TableInfo{
		Id:tableIdInt,
		BusinessId:userIdInt,
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
func (c *TableController) DeleteByBy(userId, tableId string) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized,errors.New("没有该权限")
	}

	tableIdInt,_ := strconv.Atoi(tableId)

	status, err = c.DeleteTable(user.Id,tableIdInt)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


