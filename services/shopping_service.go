package services

import (
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/model"
	"errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/chikong/ordersystem/constant"
)

type ShoppingService interface {
	GetShopping(businessId, userId int) (int, *model.ShoppingCart, error)
	UpdateShopping(userId int, businessId int, foodId int, num int, foodType string) (int, error)
	//DeleteShopping(id int) (int, error)
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
func (s *shoppingService) GetShopping(businessId, userId int) (int, *model.ShoppingCart, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}
	item := new(model.ShoppingCart)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.ColumnUserId), businessId, userId).
		Desc(constant.NameID).Get(item)
	if err != nil {
		logrus.Errorf("获取购物车失败: %s", err)
		return iris.StatusInternalServerError,  nil, errors.New("获取购物车失败")
	}
	if res == false {
		return iris.StatusNotFound, nil, errors.New("购物车为空")
	}

	price, err := s.MenuService.GetOrderSumPrice(item.FoodList)
	if err != nil {
		return iris.StatusInternalServerError, nil, err
	}

	var count int
	for _, subItem := range item.FoodList {
		count += subItem.Num
	}

	item.Count = count
	item.Price = price

	return iris.StatusOK, item, nil
}

// 修改购物车
func (s *shoppingService) UpdateShopping(userId int, businessId int,
	foodId int, num int, foodType string) (int, error) {
	status, shoppingCart, err := s.GetShopping(businessId, userId)
	if status == iris.StatusInternalServerError {
		return status, err
	}

	status,food, err := s.MenuService.GetFood(foodId)
	if err != nil {
		return status,err
	}

	if shoppingCart != nil {
		isExist := false
		// 设置修改信息
		for i, subItem := range shoppingCart.FoodList {
			// 发生变化才更新
			if subItem.Id == foodId && subItem.Type == subItem.Type {
				isExist = true
				if subItem.Num != num {
					if num == 0 { // 删除食物
						shoppingCart.FoodList = append(shoppingCart.FoodList[:i],shoppingCart.FoodList[i+1:]...)
					}else {
						subItem.Num = num
						shoppingCart.FoodList[i] = subItem
					}
					_, err = manager.DBEngine.AllCols().Where(
						fmt.Sprintf("%s=?", constant.NameID), shoppingCart.Id).Update(shoppingCart)
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
			food.Num = num
			list := append(mySlice,*food)
			shoppingCart.FoodList = list
			_, err = manager.DBEngine.AllCols().Where(
				fmt.Sprintf("%s=?", constant.NameID), shoppingCart.Id).Update(shoppingCart)
			if err != nil {
				logrus.Errorf("修改购物车失败: %s", err)
				return iris.StatusInternalServerError, errors.New("修改购物车失败")
			}

		}
	} else {
		food.Num = num
		list := []model.Food{
			*food,
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

	return iris.StatusOK, nil
}

//// 删除购物车
//func (s *shoppingService) DeleteShopping(businessId, userId, foodId int) (int, error) {
//	if businessId == 0 {
//		return iris.StatusBadRequest, errors.New("商家id不能为空")
//	}
//	if userId == 0 {
//		return iris.StatusBadRequest, errors.New("用户id不能为空")
//	}
//	_, err := manager.DBEngine.Where(
//		fmt.Sprintf("%s=?", constant.NameID), id).Delete(new(model.ShoppingCart))
//	if err != nil {
//		logrus.Errorf("删除购物车失败: %s", err)
//		return iris.StatusInternalServerError, errors.New("删除购物车失败")
//	}
//	return iris.StatusOK, nil
//}
