package services

import (
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/model"
	"errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/network"
)

type ShoppingService interface {
	GetShopping(businessId, userId, table int) (int, *model.ShoppingCart, error)
	UpdateShopping(foodType string, userId, businessId, foodId, num, tableId int) (int, error)
}

func NewShoppingService(userService UserService, menuService MenuService) ShoppingService {
	return &shoppingService{
		MenuService: menuService,
		UserService: userService,
	}
}

type shoppingService struct {
	MenuService MenuService
	UserService UserService
}

// 获取单个购物车
func (s *shoppingService) GetShopping(businessId, userId, tableId int) (int, *model.ShoppingCart, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}
	item := new(model.ShoppingCart)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.ColumnTableId), businessId, tableId).
		Desc(constant.NameID).Get(item)
	if err != nil {
		logrus.Errorf("获取购物车失败: %s", err)
		return iris.StatusInternalServerError,  nil, errors.New("获取购物车失败")
	}
	if res == false {
		return iris.StatusNotFound, nil, errors.New("购物车为空")
	}

	var count int
	for i, subItem := range item.FoodList {
		status,dbFood, err := s.MenuService.GetFood(businessId, userId, subItem.Id)
		if err != nil {
			return status, nil, err
		}
		dbFood.Type = subItem.Type
		dbFood.Num = subItem.Num
		count += dbFood.Num
		item.FoodList[i] = *dbFood.GetFood()
	}

	price, err := s.MenuService.GetOrderSumPrice(item.FoodList)
	if err != nil {
		return iris.StatusInternalServerError, nil, err
	}

	item.Count = count
	item.Price = price

	return iris.StatusOK, item, nil
}

// 修改购物车
func (s *shoppingService) UpdateShopping(foodType string, userId, businessId,
	foodId, num, tableId int) (int, error) {
	status, shoppingCart, err := s.GetShopping(businessId, userId,tableId)
	if status == iris.StatusInternalServerError {
		return status, err
	}

	status,dbFood, err := s.MenuService.GetFood(businessId, userId, foodId)
	if err != nil {
		return status,err
	}


	// 购物车数据库不需要保存这些信息
	dbFood.Price = 0
	dbFood.Pic = ""
	dbFood.Desc = ""
	dbFood.Name = ""
	dbFood.ClassifyId = nil
	dbFood.SaleCount = 0

	if shoppingCart != nil {
		isExist := false
		// 设置修改信息
		for i, subItem := range shoppingCart.FoodList {
			// 发生变化才更新
			if subItem.Id == foodId && foodType == foodType {
				isExist = true
				if subItem.Num != num {
					if num == 0 { // 删除食物
						shoppingCart.FoodList = append(shoppingCart.FoodList[:i],shoppingCart.FoodList[i+1:]...)
					}else {
						subItem.Num = num
						shoppingCart.FoodList[i] = subItem
					}
					_, err = manager.DBEngine.AllCols().Where(
						fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
							businessId,shoppingCart.Id).Update(shoppingCart)
					if err != nil {
						logrus.Errorf("修改购物车失败: %s", err)
						return iris.StatusInternalServerError, errors.New("修改购物车失败")
					}
					break
				}

			}
		}
		if !isExist {
			mySlice := shoppingCart.FoodList[:]
			dbFood.Num = num
			list := append(mySlice,*dbFood.GetFood())
			shoppingCart.FoodList = list
			_, err = manager.DBEngine.AllCols().Where(
				fmt.Sprintf("%s=?", constant.NameID), shoppingCart.Id).Update(shoppingCart)
			if err != nil {
				logrus.Errorf("修改购物车失败: %s", err)
				return iris.StatusInternalServerError, errors.New("修改购物车失败")
			}

		}
	} else {
		dbFood.Num = num
		list := []model.Food{
			*dbFood.GetFood(),
		}
		_, err = manager.DBEngine.Insert(&model.ShoppingCart{
			UserId:     userId,
			BusinessId: businessId,
			FoodList: list,
		})
		if err != nil {
			logrus.Errorf("添加食物到购物车失败: %s", err)
			return iris.StatusInternalServerError, errors.New("添加食物到购物车失败")
		}
	}
	network.SendShoppingCartMessage(businessId,tableId)
	return iris.StatusOK, nil
}
