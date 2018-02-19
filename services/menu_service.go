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

type MenuService interface {
	GetDishesList(businessId int) (int, []model.Dishes, error)
	GetDishes(dishesId int) (int, *model.Dishes, error)
	InsertDishesOne(dishes *model.Dishes) (int, error)
	InsertDishes(dishes []*model.Dishes) (int, error)
	UpdateDishes(dishes *model.Dishes) (int, error)
	DeleteDishes(dishesId int) (int, error)
	GetOrderSumPrice(dishesList []model.Dishes) (float32, error)
	GetCollectList(userId, businessId int) (int,[]model.Dishes, error)
	UpdateCollectList(userId, businessId, dishesId int, isCollect bool) (int, error)
}

func NewMenuService() MenuService {
	return &menuService{}
}

type menuService struct {
}

// 获取菜单
func (s *menuService) GetDishesList(businessId int) (int, []model.Dishes, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.Dishes, 0)

	err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).Find(&list)
	if err != nil {
		logrus.Errorf("获取菜式失败: %s", err)
		return iris.StatusInternalServerError, nil, err
	}

	return iris.StatusOK, list, nil
}

// 获取单个菜式
func (s *menuService) GetDishes(dishesId int) (int, *model.Dishes, error) {
	if dishesId == 0 {
		return iris.StatusBadRequest, nil, errors.New("菜式id不能为空")
	}
	item := new(model.Dishes)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID), dishesId).Get(item)
	if err != nil {
		logrus.Errorf("获取菜式失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取菜式失败")
	}
	if res == false {
		logrus.Errorf("菜式不存在: %s", dishesId)
		return iris.StatusNotFound, nil, errors.New("菜式不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加菜式
func (s *menuService) InsertDishesOne(dishes *model.Dishes) (int, error) {

	if dishes.Name == "" || dishes.Price == 0 {
		return iris.StatusBadRequest, errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.InsertOne(dishes)
	if err != nil {
		logrus.Errorf("添加菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加菜式失败")
	}
	return iris.StatusOK, nil
}

// 添加菜式
func (s *menuService) InsertDishes(list []*model.Dishes) (int, error) {

	for i, subItem := range list  {
		if subItem.Name == "" || subItem.Price == 0 {
			return iris.StatusBadRequest, errors.New(
				fmt.Sprintf("菜式信息不能为空: %s",i))
		}
	}
	_, err := manager.DBEngine.Insert(list)
	if err != nil {
		logrus.Errorf("添加菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加菜式失败")
	}
	return iris.StatusOK, nil
}

// 修改菜式
func (s *menuService) UpdateDishes(dishes *model.Dishes) (int, error) {
	if dishes.Id == 0 || dishes.Name == "" || dishes.Price == 0 {
		return iris.StatusBadRequest, errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		dishes.BusinessId, dishes.Id).Update(dishes)
	if err != nil {
		logrus.Errorf("修改菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改菜式失败")
	}
	return iris.StatusOK, nil
}

// 删除菜式
func (s *menuService) DeleteDishes(dishesId int) (int, error) {
	if dishesId == 0 {
		return iris.StatusBadRequest, errors.New("菜式id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID),dishesId).Delete(new(model.Dishes))
	if err != nil {
		logrus.Errorf("删除菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除菜式失败")
	}
	return iris.StatusOK, nil
}

// 获取收藏菜式
func (s *menuService) GetCollectList(userId, businessId int) (int,[]model.Dishes, error){
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}

	status, item, err := getCollectDishes(userId,businessId)
	if err != nil {
		return status,nil, err
	}
	if item == nil {
		return iris.StatusOK, []model.Dishes{}, nil
	}

	status, dishesList, err := s.GetDishesList(businessId)
	if err != nil {
		return status,nil, err
	}

	contain := func(ids[]int, id int) bool{
		for _, subItem := range ids {
			if subItem == id{
				return true
			}
		}
		return false
	}
	ids := item.CollectDishesId
	for i := 0; i < len(dishesList); {
		subItem := dishesList[i]
		if contain(ids,subItem.Id) {
			i++
		}else{
			dishesList = append(dishesList[:i],dishesList[i+1:]...)
		}
	}

	return iris.StatusOK, dishesList, nil
}

// (取消)收藏菜式
func (s *menuService) UpdateCollectList(userId, businessId, dishesId int, isCollect bool) (int, error){
	if userId == 0 {
		return iris.StatusBadRequest, errors.New("用户id不能为空")
	}
	if dishesId == 0 {
		return iris.StatusBadRequest, errors.New("菜式id不能为空")
	}

	status, item, err := getCollectDishes(userId,businessId)
	if err != nil {
		return status, err
	}
	if item == nil {
		item = &model.CollectDishes{
			UserId:userId,
			BusinessId:businessId,
			CollectDishesId:make([]int,1),
		}
	}
	ids := item.CollectDishesId
	contain := func(ids[]int, id int) (bool,int){
		for i, subItem := range ids {
			if subItem == id{
				return true,i
			}
		}
		return false,-1
	}

	isExist, i := contain(ids,dishesId)
	if isCollect {
		if isExist {
			return iris.StatusConflict,errors.New("已经收藏")
		}
		item.CollectDishesId = append(ids,dishesId)
	}else{
		if !isExist {
			return iris.StatusBadRequest,errors.New("该菜式不在收藏列表")
		}
		if len(ids) == 0 {
			return iris.StatusBadRequest,errors.New("收藏列表为空")

		}
		item.CollectDishesId = append(ids[:i],ids[i+1:]...)
	}
	_, err = manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnUserId, constant.ColumnBusinessId),
		userId,businessId).Update(item)
	if err != nil {
		logrus.Errorf("修改收藏列表失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改收藏列表失败")
	}

	return iris.StatusOK,nil
}

// 查询收藏的菜式
func getCollectDishes(userId, businessId int)(int, *model.CollectDishes, error){

	item := new(model.CollectDishes)
	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnUserId,
			constant.ColumnBusinessId), userId,businessId).Get(item)
	if err != nil {
		logrus.Errorf("获取收藏菜式失败: %s", err)
		return iris.StatusInternalServerError,nil, errors.New("获取收藏菜式失败")
	}
	if !res {
		return iris.StatusOK,nil, nil
	}

	return iris.StatusOK,item,nil
}


// 计算订单总价
func (s *menuService)GetOrderSumPrice(dishesList []model.Dishes) (float32, error) {
	var sum float32
	for _, item := range dishesList {
		sum += item.Price * float32(item.Num)
	}
	return sum, nil

}



