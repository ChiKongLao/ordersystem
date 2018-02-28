package services

import (
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/manager"
	"github.com/sirupsen/logrus"
	"errors"
	"fmt"
)

type HomeService interface {
	GetBusinessHome(userId int) (int, interface{}, error)
	GetCustomerHome(businessId, userId int) (int, interface{}, error)
}

func NewHomeService(userService UserService, foodService MenuService,
	tableService TableService, orderService OrderService,
	shopService ShopService) HomeService {
	return &homeService{
		UserService:  userService,
		MenuService:  foodService,
		TableService: tableService,
		OrderService: orderService,
		ShopService:  shopService,
	}
}

type homeService struct {
	MenuService  MenuService
	UserService  UserService
	TableService TableService
	OrderService OrderService
	ShopService  ShopService
}

// 获取商家端首页
func (s *homeService) GetBusinessHome(userId int) (int, interface{}, error) {
	status, tableList, err := s.TableService.GetTableList(userId, constant.TableStatusUnknown)
	if err != nil {
		return status, nil, err
	}

	var eatingNum, eatingPerson, emptyTable, saleOutNum int // 就餐桌数, 就餐人数,空桌数,售罄食物数量
	for _, subItem := range tableList {
		switch subItem.Status {
		case constant.TableStatusEmpty:
			emptyTable ++
		case constant.TableStatusUsing:
			eatingNum ++
			eatingPerson += subItem.PersonNum
		}
	}

	_, err = manager.DBEngine.Table("food").
		Select("Count(num) AS saleOutNum").
		Where(fmt.Sprintf("%s=0", constant.ColumnNum)).
		Get(&saleOutNum)

	if err != nil {
		logrus.Errorf("获取售罄食物数量失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取首页数据失败")
	}

	type Home struct {
		EatingNum    int               `json:"eatingNum"`
		EatingPerson int               `json:"eatingPerson"`
		EmptyTable   int               `json:"emptyTable"`
		SaleOutNum   int               `json:"saleOutNum"`
		Data         []model.TableInfo `json:"data"`
	}

	return iris.StatusOK, &Home{
		Data:         tableList,
		EatingNum:    eatingNum,
		EatingPerson: eatingPerson,
		EmptyTable:   emptyTable,
		SaleOutNum:   saleOutNum,
	}, nil
}

// 获取用户端首页
func (s *homeService) GetCustomerHome(businessId, userId int) (int, interface{}, error) {
	status, foodMap, err := s.MenuService.GetFoodList(businessId, userId)
	if err != nil {
		return status, nil, err
	}
	status, user, err := s.UserService.GetUserById(businessId)
	if err != nil {
		return status, nil, err
	}
	status, shop, err := s.ShopService.GetShop(businessId)
	if err != nil {
		return status, nil, err
	}
	type Home struct {
		Name string                          `json:"name"`
		Desc string                          `json:"desc"`
		Pic  string                          `json:"pic"`
		Food []model.FoodResponse `json:"food"`
	}
	foodList := make([]model.FoodResponse,0)

	for _, value := range foodMap {
		foodList = append(foodList,value...)
	}
	return iris.StatusOK, &Home{
		Food: foodList,
		Name: user.NickName,
		Pic:  shop.Pic,
	}, nil
}
