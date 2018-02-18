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
	UpdateShopping(userId int, businessId int, dashesId int, num int, dashesType string) (int, error)
	//DeleteShopping(id int) (int, error)
}

func NewShoppingService(UserService UserService, MenuService MenuService) ShoppingService {
	return &shoppingService{
		MenuService: MenuService,
		UserService: UserService,
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

	price, err := s.MenuService.GetOrderSumPrice(item.DashesList)
	if err != nil {
		return iris.StatusInternalServerError, nil, err
	}

	var count int
	for _, subItem := range item.DashesList {
		count += subItem.Num
	}

	item.Count = count
	item.Price = price

	return iris.StatusOK, item, nil
}

// 修改购物车
func (s *shoppingService) UpdateShopping(userId int, businessId int,
	dashesId int, num int, dashesType string) (int, error) {
	status, shoppingCart, err := s.GetShopping(businessId, userId)
	if status == iris.StatusInternalServerError {
		return status, err
	}

	status,dashes, err := s.MenuService.GetDashes(dashesId)
	if err != nil {
		return status,err
	}

	if shoppingCart != nil {
		isExist := false
		// 设置修改信息
		for i, subItem := range shoppingCart.DashesList {
			// 发生变化才更新
			if subItem.Id == dashesId && subItem.Type == subItem.Type {
				isExist = true
				if subItem.Num != num {
					if num == 0 { // 删除菜式
						shoppingCart.DashesList = append(shoppingCart.DashesList[:i],shoppingCart.DashesList[i+1:]...)
					}else {
						subItem.Num = num
						shoppingCart.DashesList[i] = subItem
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
			mySlice := shoppingCart.DashesList[:]
			dashes.Num = num
			list := append(mySlice,*dashes)
			shoppingCart.DashesList = list
			_, err = manager.DBEngine.AllCols().Where(
				fmt.Sprintf("%s=?", constant.NameID), shoppingCart.Id).Update(shoppingCart)
			if err != nil {
				logrus.Errorf("修改购物车失败: %s", err)
				return iris.StatusInternalServerError, errors.New("修改购物车失败")
			}

		}
	} else {
		dashes.Num = num
		list := []model.Dashes{
			*dashes,
		}
		_, err = manager.DBEngine.Insert(&model.ShoppingCart{
			UserId:     userId,
			BusinessId: businessId,
			DashesList: list,
		})
		if err != nil {
			logrus.Errorf("添加菜式到购物车失败: %s", err)
			return iris.StatusInternalServerError, errors.New("添加菜式到购物车失败")
		}
	}

	return iris.StatusOK, nil
}

//// 删除购物车
//func (s *shoppingService) DeleteShopping(businessId, userId, dashesId int) (int, error) {
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
