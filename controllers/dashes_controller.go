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

// 菜式
type DashesController struct {
	Ctx iris.Context
	services.MenuService
	UserService services.UserService

}

// 获取菜单
func (c *DashesController) GetBy(userId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	var list []model.Dashes
	status, list, err = c.GetDashesList(userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}

// 获取菜式详情
func (c *DashesController) GetByBy(userId, dashId int) (int,interface{}) {

	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.Dashes
	status, item, err = c.GetDashes(dashId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 添加菜式
func (c *DashesController) PostBy(userId int) (int,interface{}) {
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
	//num,_ := strconv.Atoi(c.Ctx.FormValue(constant.NameNum)) // 暂时不用数量
	pic := c.Ctx.FormValue(constant.NamePic)
	price, _ := strconv.ParseFloat(c.Ctx.FormValue(constant.NamePrice),10)
	dashesType := c.Ctx.FormValue(constant.NameType)
	desc := c.Ctx.FormValue(constant.NameDesc)


	status, err = c.InsertDashesOne(&model.Dashes{
		BusinessId:userId,
		Name:name,
		Num:100,
		Pic:pic,
		Price:float32(price),
		Type:dashesType,
		Desc:desc,
	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 修改菜式
func (c *DashesController) PutByBy(userId, dashId int) (int,interface{}) {
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
	num,_ := strconv.Atoi(c.Ctx.FormValue(constant.NameNum))
	pic := c.Ctx.FormValue(constant.NamePic)
	price, _ := strconv.ParseFloat(c.Ctx.FormValue(constant.NamePrice),10)
	dashesType := c.Ctx.FormValue(constant.NameType)
	desc := c.Ctx.FormValue(constant.NameDesc)


	status, err = c.UpdateDashes(&model.Dashes{
		Id:dashId,
		BusinessId:userId,
		Name:name,
		Num:num,
		Pic:pic,
		Price:float32(price),
		Type:dashesType,
		Desc:desc,
	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 删除菜式
func (c *DashesController) DeleteByBy(userId, dashId int) (int,interface{}) {
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

	status, err = c.DeleteDashes(dashId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


