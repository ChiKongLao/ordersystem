package controllers

import (
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
)

// 店铺
type ShopController struct {
	Ctx         iris.Context
	services.ShopService
	UserService services.UserService
}

// 获取店铺
func (c *ShopController) Get() (int, interface{}) {
	status, res, err := c.UserService.CheckRoleIsManagerWithToken(c.Ctx)
	if !res {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	var list []model.Shop
	status, list, err = c.GetShopList()
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameData:  list,
		constant.NameCount: len(list),
	}
}

// 获取店铺详情
func (c *ShopController) GetBy(businessId int) (int, interface{}) {

	status, _, err := c.UserService.GetBusinessById(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}
	var item *model.Shop
	status, item, err = c.GetShop(businessId)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameData: item,
	}
}

// 添加店铺
func (c *ShopController) Post() (int, interface{}) {

	status, res, err := c.UserService.CheckRoleIsManagerWithToken(c.Ctx)
	if !res {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	businessId,_ := c.Ctx.PostValueInt(constant.NameBusinessID)
	name := c.Ctx.PostValue(constant.Name)
	desc := c.Ctx.PostValue(constant.NameDesc)
	pic := c.Ctx.PostValue(constant.NamePic)

	status, err = c.InsertShop(businessId,name,desc,pic)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 修改店铺
func (c *ShopController) Put() (int, interface{}) {
	status, res, err := c.UserService.CheckRoleIsManagerWithToken(c.Ctx)
	if !res {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}

	businessId,_ := c.Ctx.PostValueInt(constant.NameBusinessID)
	name := c.Ctx.PostValue(constant.Name)
	desc := c.Ctx.PostValue(constant.NameDesc)
	pic := c.Ctx.PostValue(constant.NamePic)

	status, err = c.UpdateShop(businessId,name,desc,pic)
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 删除店铺
func (c *ShopController) Delete() (int, interface{}) {
	status, res, err := c.UserService.CheckRoleIsManagerWithToken(c.Ctx)
	if !res {
		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
	}
	businessId,_ := c.Ctx.PostValueInt(constant.NameBusinessID)
	status, err = c.DeleteShop(businessId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}
