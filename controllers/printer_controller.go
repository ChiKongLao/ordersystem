package controllers

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/model"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/util"
)

// 打印机
type PrinterController struct {
	Ctx         iris.Context
	services.PrinterService
	UserService services.UserService
}

// 获取打印机
func (c *PrinterController) Get() (int, interface{}) {
	status, ok, err := c.UserService.CheckRoleIsBusinessWithToken(c.Ctx)
	if err != nil || !ok {
		return status, model.NewErrorResponse(err)
	}

	var list []model.Printer
	status, list, err = c.GetPrinterList()
	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameData:  list,
		constant.NameCount: len(list),
	}
}

// 添加打印机
func (c *PrinterController) Post() (int, interface{}) {
	status, ok, err := c.UserService.CheckRoleIsBusinessWithToken(c.Ctx)
	if err != nil || !ok {
		return status, model.NewErrorResponse(err)
	}

	name := c.Ctx.PostValue(constant.Name)
	businessId, _ := c.Ctx.PostValueInt(constant.NameBusinessID)

	status, err = c.InsertPrinter(businessId, name)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 修改打印机
func (c *PrinterController) PutBy(id int) (int, interface{}) {
	status, ok, err := c.UserService.CheckRoleIsBusinessWithToken(c.Ctx)
	if err != nil || !ok {
		return status, model.NewErrorResponse(err)
	}

	name := c.Ctx.PostValue(constant.Name)
	businessId, _ := c.Ctx.PostValueInt(constant.NameBusinessID)

	status, err = c.UpdatePrinter(id, businessId, name)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 删除打印机
func (c *PrinterController) DeleteByBy(id int) (int, interface{}) {
	status, ok, err := c.UserService.CheckRoleIsBusinessWithToken(c.Ctx)
	if err != nil || !ok {
		return status, model.NewErrorResponse(err)
	}
	status, err = c.DeletePrinter(id)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return status, iris.Map{
		constant.NameIsOk: true,
	}
}

// 测试打印机
func (c *PrinterController) PostTest() (int, interface{}) {
	//printerService.HandlePayload("A*888*1*AS01#")

	foodList := make([]model.Food, 0)

	foodList = append(foodList, model.Food{
		Name:  "玉米牛肉肠",
		Num:   2,
		Price: 10,
	})
	foodList = append(foodList, model.Food{
		Name:  "生菜鱼片粥",
		Num:   1,
		Price: 7,
	})

	order := model.OrderPrint{
		OrderResponse: model.OrderResponse{
			Order: model.Order{
				Id:         123,
				OrderNo:    "123",
				PersonNum:  3,
				Price:      43,
				CreateTime: util.GetCurrentTime(),
				BusinessId: 227,
				UserId:     2,
				FoodList:   foodList,
			},
			TableName: "桃花阁",
		},
		Customer: model.User{
			UserName: "大爷",
		},
		Business: model.User{
			Address: "广东省广州市天河区天河城88号",
			Phone:   "18888888888",
		},
	}
	c.SendOrder(order)
	return iris.StatusOK, nil
}
