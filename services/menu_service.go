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
	GetFoodList(businessId, userId int) (int, []model.Food, error)
	GetFood(foodId int) (int, *model.Food, error)
	InsertFoodOne(food *model.Food) (int, error)
	InsertFood(food []*model.Food) (int, error)
	UpdateFood(food *model.Food) (int, error)
	DeleteFood(foodId int) (int, error)
	GetOrderSumPrice(foodList []model.Food) (float32, error)
	GetCollectList(userId, businessId int) (int,[]model.Food, error)
	UpdateCollectList(userId, businessId, foodId int, isCollect bool) (int, error)
}

func NewMenuService(userService UserService) MenuService {
	return &menuService{
		UserService: userService,
	}
}

type menuService struct {
	UserService UserService
}

// 获取菜单
func (s *menuService) GetFoodList(businessId, userId int) (int, []model.Food, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.Food, 0)

	err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).Find(&list)
	if err != nil {
		logrus.Errorf("获取食物失败: %s", err)
		return iris.StatusInternalServerError, nil, err
	}
	var user *model.User
	if _, user, err = s.UserService.GetUserById(userId); user.IsCustomer() { // 客户才查看收藏
		_, collectItem, _ := getCollectFood(userId, businessId) // 获取收藏的食物
		if collectItem != nil {
			for i, subItem := range list {
				for _, collectId := range collectItem.CollectFoodId {
					if subItem.Id == collectId {
						subItem.IsCollect = true
						list[i] = subItem
						break
					}
				}
			}
		}
	}
	return iris.StatusOK, list, nil
}

// 获取单个食物
func (s *menuService) GetFood(foodId int) (int, *model.Food, error) {
	if foodId == 0 {
		return iris.StatusBadRequest, nil, errors.New("食物id不能为空")
	}
	item := new(model.Food)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID), foodId).Get(item)
	if err != nil {
		logrus.Errorf("获取食物失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取食物失败")
	}
	if res == false {
		logrus.Errorf("食物不存在: %s", foodId)
		return iris.StatusNotFound, nil, errors.New("食物不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加食物
func (s *menuService) InsertFoodOne(food *model.Food) (int, error) {

	if food.Name == "" || food.Price == 0 {
		return iris.StatusBadRequest, errors.New("食物信息不能为空")
	}

	_, err := manager.DBEngine.InsertOne(food)
	if err != nil {
		logrus.Errorf("添加食物失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加食物失败")
	}
	return iris.StatusOK, nil
}

// 添加食物
func (s *menuService) InsertFood(list []*model.Food) (int, error) {

	for i, subItem := range list  {
		if subItem.Name == "" || subItem.Price == 0 {
			return iris.StatusBadRequest, errors.New(
				fmt.Sprintf("食物信息不能为空: %s",i))
		}
	}
	_, err := manager.DBEngine.Insert(list)
	if err != nil {
		logrus.Errorf("添加食物失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加食物失败")
	}
	return iris.StatusOK, nil
}

// 修改食物
func (s *menuService) UpdateFood(food *model.Food) (int, error) {
	if food.Id == 0 || food.Name == "" || food.Price == 0 {
		return iris.StatusBadRequest, errors.New("食物信息不能为空")
	}

	_, err := manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		food.BusinessId, food.Id).Update(food)
	if err != nil {
		logrus.Errorf("修改食物失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改食物失败")
	}
	return iris.StatusOK, nil
}

// 删除食物
func (s *menuService) DeleteFood(foodId int) (int, error) {
	if foodId == 0 {
		return iris.StatusBadRequest, errors.New("食物id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID),foodId).Delete(new(model.Food))
	if err != nil {
		logrus.Errorf("删除食物失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除食物失败")
	}
	return iris.StatusOK, nil
}

// 获取收藏食物
func (s *menuService) GetCollectList(userId, businessId int) (int,[]model.Food, error){
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}

	status, item, err := getCollectFood(userId,businessId)
	if err != nil {
		return status,nil, err
	}
	if item == nil {
		return iris.StatusOK, []model.Food{}, nil
	}

	status, foodList, err := s.GetFoodList(businessId,userId)
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
	ids := item.CollectFoodId
	for i := 0; i < len(foodList); {
		subItem := foodList[i]
		if contain(ids,subItem.Id) {
			i++
		}else{
			foodList = append(foodList[:i],foodList[i+1:]...)
		}
	}

	return iris.StatusOK, foodList, nil
}

// (取消)收藏食物
func (s *menuService) UpdateCollectList(userId, businessId, foodId int, isCollect bool) (int, error){
	if userId == 0 {
		return iris.StatusBadRequest, errors.New("用户id不能为空")
	}
	if foodId == 0 {
		return iris.StatusBadRequest, errors.New("食物id不能为空")
	}

	status, item, err := getCollectFood(userId,businessId)
	if err != nil {
		return status, err
	}
	if item == nil {
		item = &model.CollectFood{
			UserId:userId,
			BusinessId:businessId,
			CollectFoodId:make([]int,1),
		}
	}
	ids := item.CollectFoodId
	contain := func(ids[]int, id int) (bool,int){
		for i, subItem := range ids {
			if subItem == id{
				return true,i
			}
		}
		return false,-1
	}

	isExist, i := contain(ids,foodId)
	if isCollect {
		if isExist {
			return iris.StatusConflict,errors.New("已经收藏")
		}
		item.CollectFoodId = append(ids,foodId)
	}else{
		if !isExist {
			return iris.StatusBadRequest,errors.New("该食物不在收藏列表")
		}
		if len(ids) == 0 {
			return iris.StatusBadRequest,errors.New("收藏列表为空")

		}
		item.CollectFoodId = append(ids[:i],ids[i+1:]...)
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

// 查询收藏的食物
func getCollectFood(userId, businessId int)(int, *model.CollectFood, error){

	item := new(model.CollectFood)
	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnUserId,
			constant.ColumnBusinessId), userId,businessId).Get(item)
	if err != nil {
		logrus.Errorf("获取收藏食物失败: %s", err)
		return iris.StatusInternalServerError,nil, errors.New("获取收藏食物失败")
	}
	if !res {
		return iris.StatusOK,nil, nil
	}

	return iris.StatusOK,item,nil
}


// 计算订单总价
func (s *menuService)GetOrderSumPrice(foodList []model.Food) (float32, error) {
	var sum float32
	for _, item := range foodList {
		sum += item.Price * float32(item.Num)
	}
	return sum, nil

}



