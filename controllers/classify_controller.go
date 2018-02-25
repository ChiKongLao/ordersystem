package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/api/middleware/authentication"
)

// 分类
type ClassifyController struct {
	Ctx iris.Context
	services.ClassifyService
	UserService services.UserService

}

// 获取分类
func (c *ClassifyController) GetBy(businessId int) (int,interface{}) {
	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}


	var list []model.Classify
	status, list, err = c.GetClassifyList(businessId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:list,
		constant.NameCount:len(list),
		}
}

// 获取分类详情
func (c *ClassifyController) GetByBy(businessId, classifyId int) (int,interface{}) {

	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}
	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.Classify
	status, item, err = c.GetClassify(businessId,classifyId)
	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
		constant.NameData:item,
		}
}

// 添加分类
func (c *ClassifyController) PostBy(businessId int) (int,interface{}) {
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

	name := c.Ctx.PostValue(constant.Name)
	sort, _ := c.Ctx.PostValueIntDefault(constant.NameSort,100)

	status, err = c.InsertClassify(name,businessId,sort)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}

// 修改分类
func (c *ClassifyController) PutByBy(businessId, classifyId int) (int,interface{}) {
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

	name := c.Ctx.PostValue(constant.Name)
	sort, _ := c.Ctx.PostValueInt(constant.NameSort)

	status, err = c.UpdateClassify(name,businessId, classifyId, sort)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


// 删除分类
func (c *ClassifyController) DeleteByBy(businessId, classifyId int) (int,interface{}) {
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

	status, err = c.DeleteClassify(businessId, classifyId)

	if err != nil{
		return status, model.NewErrorResponse(err)
	}

	return status,iris.Map{
			constant.NameIsOk:true,
		}
}


