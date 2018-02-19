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
type MenuController struct {
	Ctx iris.Context
	services.MenuService
	UserService services.UserService

}

// 获取菜单
func (c *MenuController) GetBy(userId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	var list []model.Dishes
	status, list, err = c.GetDishesList(userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}

// 获取菜式详情
func (c *MenuController) GetByBy(userId, dishId int) (int,interface{}) {

	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.Dishes
	status, item, err = c.GetDishes(dishId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 添加菜式
func (c *MenuController) PostBy(userId int) (int,interface{}) {
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
	dishesType := c.Ctx.FormValue(constant.NameType)
	desc := c.Ctx.FormValue(constant.NameDesc)


	status, err = c.InsertDishesOne(&model.Dishes{
		BusinessId:userId,
		Name:name,
		Num:100,
		Pic:pic,
		Price:float32(price),
		Type:dishesType,
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
func (c *MenuController) PutByBy(userId, dishId int) (int,interface{}) {
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
	dishesType := c.Ctx.FormValue(constant.NameType)
	desc := c.Ctx.FormValue(constant.NameDesc)


	status, err = c.UpdateDishes(&model.Dishes{
		Id:dishId,
		BusinessId:userId,
		Name:name,
		Num:num,
		Pic:pic,
		Price:float32(price),
		Type:dishesType,
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
func (c *MenuController) DeleteByBy(userId, dishId int) (int,interface{}) {
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

	status, err = c.DeleteDishes(dishId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 获取用户收藏的菜单
func (c *MenuController) GetCollectBy(businessId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, item, err := c.GetCollectList(userId,businessId)

	return status,iris.Map{
		constant.NameData:item,
		constant.NameCount:len(item),
	}

}

// 修改用户收藏的菜单
func (c *MenuController) PutCollectBy(businessId int) (int,interface{}) {
	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	dishesId, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameDishesId))
	//isCollect := c.Ctx.FormValue(constant.NameIsCollect)
	isCollect,_ := c.Ctx.PostValueBool(constant.NameIsCollect)
	status, err = c.UpdateCollectList(userId,businessId,dishesId,isCollect)
	if err != nil{
		return status,model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameIsOk:true,
	}

}

