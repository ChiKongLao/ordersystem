package main

import (
	"github.com/chikong/ordersystem/services"
	"github.com/sirupsen/logrus"
	"github.com/chikong/ordersystem/model"
	"github.com/chikong/ordersystem/util"
)

func maintest() {
	logrus.SetLevel(logrus.DebugLevel)
	printerService := services.NewPrinterService()
	printerService.HandlePayload("A*888*1*AS01#")



	foodList := make([]model.Food,0)

	foodList = append(foodList,model.Food{
		Name:"玉米牛肉肠",
		Num:2,
		Price:10,
	})
	foodList = append(foodList,model.Food{
		Name:"生菜鱼片粥",
		Num:1,
		Price:7,
	})
	foodList = append(foodList,model.Food{
		Name:"蜜汁叉烧饭",
		Num:1,
		Price:16,
	})

	order := model.OrderPrint{
		OrderResponse:model.OrderResponse{
			Order:model.Order{
				Id:123,
				OrderNo:"123",
				PersonNum:3,
				Price:43,
				CreateTime:util.GetCurrentTime(),
				BusinessId:227,
				UserId: 2,
				FoodList:foodList,
				},
			TableName:"桃花阁",
		},
		Customer:model.User{
			UserName:"大爷",
		},
		Business:model.User{
			Address:"广东省广州市天河区天河城88号",
			Phone:"18888888888",

		},
	}
	printerService.SendOrder(order)

	printerService.HandlePayload("A*888*2345*AS04#")
	printerService.HandlePayload("A*888*2345*AS05#")
	printerService.HandlePayload("A*888*2345*AS06#")
	printerService.HandlePayload("A*888*2345*AS07#")


	printerService.HandlePayload("A*888*v1.2_jsffda*AS36#")
	printerService.HandlePayload("A*888*45648475845646647*AS33#")
	printerService.HandlePayload("A*888*1,2*AS35#")
	printerService.HandlePayload("A*888*888,192.168.1.100,8090*AS34#")
	printerService.HandlePayload("A*888*31*AS37#")
	printerService.HandlePayload("A*888*AS32#")
	printerService.HandlePayload("A*888*AS48#")
	printerService.HandlePayload("A*888*1*AS47#")


}
