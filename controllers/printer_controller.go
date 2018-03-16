package controllers

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
)

// 打印机
type PrinterController struct {
	Ctx iris.Context
	services.PrinterService
	UserService services.UserService

}

//// 获取打印机
//func (c *PrinterController) GetBy(businessId int) (int,interface{}) {
//	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
//	if !isOwn {
//		return iris.StatusUnauthorized, model.NewErrorResponse(err)
//	}
//	status, _, err := c.UserService.GetBusinessById(businessId)
//	if err != nil {
//		return status, model.NewErrorResponse(err)
//	}
//
//
//	var list []model.Printer
//	status, list, err = c.GetPrinterList(businessId)
//	if err != nil{
//		return status, model.NewErrorResponse(err)
//	}
//
//	return status,iris.Map{
//		constant.NameData:list,
//		constant.NameCount:len(list),
//		}
//}
//
//// 获取打印机详情
//func (c *PrinterController) GetByBy(businessId, classifyId int) (int,interface{}) {
//
//	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
//	if !isOwn {
//		return iris.StatusUnauthorized, model.NewErrorResponse(err)
//	}
//	status, _, err := c.UserService.GetBusinessById(businessId)
//	if err != nil {
//		return status, model.NewErrorResponse(err)
//	}
//	var item *model.Printer
//	status, item, err = c.GetPrinter(businessId,classifyId)
//	if err != nil{
//		return status, model.NewErrorResponse(err)
//	}
//
//	return status,iris.Map{
//		constant.NameData:item,
//		}
//}
//
//// 添加打印机
//func (c *PrinterController) PostBy(businessId int) (int,interface{}) {
//	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
//	if !isOwn {
//		return iris.StatusUnauthorized, model.NewErrorResponse(err)
//	}
//
//	status, user, err := c.UserService.GetUserById(businessId)
//
//	if err != nil {
//		return status, model.NewErrorResponse(err)
//	}
//
//	if !user.IsManager() && !user.IsBusiness(){
//		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
//	}
//
//	name := c.Ctx.PostValue(constant.Name)
//	sort, _ := c.Ctx.PostValueIntDefault(constant.NameSort,100)
//
//	status, err = c.InsertPrinter(name,businessId,sort)
//
//	if err != nil{
//		return status, model.NewErrorResponse(err)
//	}
//
//	return status,iris.Map{
//			constant.NameIsOk:true,
//		}
//}
//
//// 修改打印机
//func (c *PrinterController) PutByBy(businessId, classifyId int) (int,interface{}) {
//	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
//	if !isOwn {
//		return iris.StatusUnauthorized, model.NewErrorResponse(err)
//	}
//
//	status, user, err := c.UserService.GetUserById(businessId)
//
//	if err != nil {
//		return status, model.NewErrorResponse(err)
//	}
//
//	if !user.IsManager() && !user.IsBusiness(){
//		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
//	}
//
//	name := c.Ctx.PostValue(constant.Name)
//	sort, _ := c.Ctx.PostValueInt(constant.NameSort)
//
//	status, err = c.UpdatePrinter(name,businessId, classifyId, sort)
//
//	if err != nil{
//		return status, model.NewErrorResponse(err)
//	}
//
//	return status,iris.Map{
//			constant.NameIsOk:true,
//		}
//}
//
//
//// 删除打印机
//func (c *PrinterController) DeleteByBy(businessId, classifyId int) (int,interface{}) {
//	isOwn, err := authentication.IsOwnWithToken(c.Ctx, businessId)
//	if !isOwn {
//		return iris.StatusUnauthorized, model.NewErrorResponse(err)
//	}
//
//	status, user, err := c.UserService.GetUserById(businessId)
//
//	if err != nil {
//		return status, model.NewErrorResponse(err)
//	}
//
//	if !user.IsManager() && !user.IsBusiness(){
//		return iris.StatusUnauthorized, model.NewErrorResponseWithMsg("没有该权限")
//	}
//
//	status, err = c.DeletePrinter(businessId, classifyId)
//
//	if err != nil{
//		return status, model.NewErrorResponse(err)
//	}
//
//	return status,iris.Map{
//			constant.NameIsOk:true,
//		}
//}


