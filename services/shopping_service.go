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
	GetShoppingList(businessId, userId, table int) (int, []model.ShoppingCartResponse, error)
	GetShopping(businessId, userId, table int) (int, *model.ShoppingCartResponse, error)
	UpdateShopping(foodType string, userId, businessId, foodId, num, tableId int) (int, error)
	DeleteShopping(businessId, shoppingCartId int) (int, error)
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

// 获取购物车
func (s *shoppingService) GetShoppingList(businessId, userId, tableId int) (int, []model.ShoppingCartResponse, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}
	if tableId == 0 {
		return iris.StatusBadRequest, nil, errors.New("餐桌id不能为空")
	}
	list := make([]model.ShoppingCartResponse, 0)

	err := manager.DBEngine.Table("shopping_cart").Select("shopping_cart.*,`user`.nick_name").
		Join("INNER", "`user`", "shopping_cart.user_id = `user`.id").
		Where(fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.ColumnTableId), businessId, tableId).
		Desc("shopping_cart.id").
		Find(&list)

	if err != nil {
		logrus.Errorf("获取购物车失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取购物车失败")
	}

	for i, subItem := range list {
		var count int
		for i, foodItem := range subItem.FoodList {
			status, dbFood, err := s.MenuService.GetFood(businessId, userId, foodItem.Id)
			if err != nil {
				return status, nil, err
			}
			dbFood.Type = foodItem.Type
			dbFood.Num = foodItem.Num
			count += dbFood.Num
			subItem.FoodList[i] = *dbFood.GetFood()
		}

		price, err := s.MenuService.GetOrderSumPrice(subItem.FoodList)
		if err != nil {
			return iris.StatusInternalServerError, nil, err
		}
		subItem.Count = count
		subItem.Price = price
		list[i] = subItem
	}
	return iris.StatusOK, list, nil
}

// 获取单个购物车
func (s *shoppingService) GetShopping(businessId, userId, tableId int) (int, *model.ShoppingCartResponse, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}
	item := new(model.ShoppingCartResponse)

	res, err := manager.DBEngine.Table("shopping_cart").Select("shopping_cart.*,`user`.nick_name").
		Join("INNER", "`user`", "shopping_cart.user_id = `user`.id").
		Where(fmt.Sprintf("%s=? and %s=? and %s=?", constant.ColumnBusinessId, constant.ColumnUserId, constant.ColumnTableId),
		businessId, userId, tableId).
		Get(item)

	if err != nil {
		logrus.Errorf("获取购物车失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取购物车失败")
	}
	if res == false {
		return iris.StatusNotFound, nil, errors.New("购物车为空")
	}

	var count int
	for i, subItem := range item.FoodList {
		status, dbFood, err := s.MenuService.GetFood(businessId, userId, subItem.Id)
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
	status, shoppingCartResponse, err := s.GetShopping(businessId, userId, tableId)
	if status == iris.StatusInternalServerError {
		return status, err
	}

	status, dbFood, err := s.MenuService.GetFood(businessId, userId, foodId)
	if err != nil {
		return status, err
	}

	// 购物车数据库不需要保存这些信息
	dbFood.Price = 0
	dbFood.Pic = ""
	dbFood.Desc = ""
	dbFood.Name = ""
	dbFood.ClassifyId = nil
	dbFood.SaleCount = 0

	if shoppingCartResponse != nil {
		isExist := false
		// 设置修改信息
		for i, subItem := range shoppingCartResponse.FoodList {
			// 发生变化才更新
			if subItem.Id == foodId && foodType == foodType {
				isExist = true
				if subItem.Num != num {
					if num == 0 { // 删除食物
						shoppingCartResponse.FoodList = append(shoppingCartResponse.FoodList[:i], shoppingCartResponse.FoodList[i+1:]...)
					} else {
						subItem.Num = num
						shoppingCartResponse.FoodList[i] = subItem
					}
					_, err = manager.DBEngine.AllCols().Where(
						fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
						businessId, shoppingCartResponse.Id).Update(shoppingCartResponse.ShoppingCart)
					if err != nil {
						logrus.Errorf("修改购物车失败: %s", err)
						return iris.StatusInternalServerError, errors.New("修改购物车失败")
					}
					break
				}

			}
		}
		if !isExist {
			mySlice := shoppingCartResponse.FoodList[:]
			dbFood.Num = num
			list := append(mySlice, *dbFood.GetFood())
			shoppingCartResponse.FoodList = list
			_, err = manager.DBEngine.AllCols().Where(
				fmt.Sprintf("%s=?", constant.NameID), shoppingCartResponse.Id).Update(shoppingCartResponse.ShoppingCart)
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
			FoodList:   list,
			TableId:    tableId,
		})
		if err != nil {
			logrus.Errorf("添加食物到购物车失败: %s", err)
			return iris.StatusInternalServerError, errors.New("添加食物到购物车失败")
		}
	}
	network.SendShoppingCartMessage(businessId, tableId)
	return iris.StatusOK, nil
}

// 删除购物车
func (s *shoppingService) DeleteShopping(businessId, shoppingCartId int) (int, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, errors.New("商家id不能为空")
	}
	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID), shoppingCartId).Delete(new(model.ShoppingCart))
	if err != nil {
		logrus.Errorf("删除购物车失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除购物车失败")
	}
	return iris.StatusOK, nil
}
