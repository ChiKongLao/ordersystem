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

	var list []model.FoodResponse
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
	var item *model.FoodResponse
	status, item, err = c.GetFood(userId,foodId)
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
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	name := c.Ctx.FormValue(constant.Name)
	num,_ := c.Ctx.PostValueInt(constant.NameNum)
	classifyId,_ := c.Ctx.PostValueInt(constant.NameType)
	pic := c.Ctx.FormValue(constant.NamePic)
	price, _ := strconv.ParseFloat(c.Ctx.FormValue(constant.NamePrice),10)
	desc := c.Ctx.FormValue(constant.NameDesc)

	//var classify model.Classify
	//err = json.Unmarshal([]byte(c.Ctx.FormValue(constant.NameType)), &classify)

	if err != nil {
		return iris.StatusBadRequest, iris.Map{constant.NameMsg: "分类格式错误"}
	}

	status, err = c.InsertFoodOne(&model.Food{
		BusinessId:userId,
		Name:name,
		Num:num,
		Pic:pic,
		Price:float32(price),
		ClassifyId:classifyId,
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
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	name := c.Ctx.FormValue(constant.Name)
	num,_ := strconv.Atoi(c.Ctx.FormValue(constant.NameNum))
	classifyId,_ := c.Ctx.PostValueInt(constant.NameType)
	pic := c.Ctx.FormValue(constant.NamePic)
	price, _ := strconv.ParseFloat(c.Ctx.FormValue(constant.NamePrice),10)
	desc := c.Ctx.FormValue(constant.NameDesc)

	var classify model.Classify
	err = json.Unmarshal([]byte(c.Ctx.FormValue(constant.NameType)), &classify)

	if err != nil {
		return iris.StatusBadRequest, iris.Map{constant.NameMsg: "分类格式错误"}
	}


	status, err = c.UpdateFood(&model.Food{
		Id:foodId,
		BusinessId:userId,
		Name:name,
		Num:num,
		Pic:pic,
		Price:float32(price),
		ClassifyId:classifyId,
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
func (c *MenuController) DeleteByBy(businessId, foodId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	status, user, err := c.UserService.GetUserById(businessId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	status, err = c.DeleteFood(businessId,foodId)

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

