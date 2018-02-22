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

// 食物
type MenuController struct {
	Ctx iris.Context
	services.MenuService
	UserService services.UserService

}

// 获取菜单
func (c *MenuController) GetBy(businessId int) (int,interface{}) {
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	status, userId, err := authentication.GetUserIDFormHeaderToken(c.Ctx)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	var list []model.Food
	status, list, err = c.GetFoodList(businessId,userId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}

// 获取食物详情
func (c *MenuController) GetByBy(userId, foodId int) (int,interface{}) {

	status, _, err := c.UserService.GetBusinessById(userId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.Food
	status, item, err = c.GetFood(foodId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 添加食物
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
	foodType := c.Ctx.FormValue(constant.NameType)
	desc := c.Ctx.FormValue(constant.NameDesc)


	status, err = c.InsertFoodOne(&model.Food{
		BusinessId:userId,
		Name:name,
		Num:100,
		Pic:pic,
		Price:float32(price),
		Type:foodType,
		Desc:desc,
	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 修改食物
func (c *MenuController) PutByBy(userId, foodId int) (int,interface{}) {
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
	foodType := c.Ctx.FormValue(constant.NameType)
	desc := c.Ctx.FormValue(constant.NameDesc)


	status, err = c.UpdateFood(&model.Food{
		Id:foodId,
		BusinessId:userId,
		Name:name,
		Num:num,
		Pic:pic,
		Price:float32(price),
		Type:foodType,
		Desc:desc,
	} )

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 删除食物
func (c *MenuController) DeleteByBy(userId, foodId int) (int,interface{}) {
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

	status, err = c.DeleteFood(foodId)

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
	foodId, _ := strconv.Atoi(c.Ctx.FormValue(constant.NameFoodId))
	isCollect,_ := c.Ctx.PostValueBool(constant.NameIsCollect)
	status, err = c.UpdateCollectList(userId,businessId,foodId,isCollect)
	if err != nil{
		return status,model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameIsOk:true,
	}

}

