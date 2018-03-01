package services

import (
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/model"
	"errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"strings"
)

type MenuService interface {
	GetFoodList(businessId, userId int) (int, map[string][]model.FoodResponse, error)
	GetFood(businessId, userId, foodId int) (int, *model.FoodResponse, error)
	InsertFoodOne(food *model.Food) (int, error)
	InsertFood(food []*model.Food) (int, error)
	UpdateFood(food *model.Food) (int, error)
	DeleteFood(businessId, foodId int) (int, error)
	SellFood(businessId, userId, foodId, num int) (int, error)
	GetOrderSumPrice(foodList []model.Food) (float32, error)
	GetCollectList(userId, businessId int) (int,[]model.FoodResponse, error)
	UpdateCollectList(userId, businessId, foodId int, isCollect bool) (int, error)
}

func NewMenuService(userService UserService, classifyService ClassifyService) MenuService {
	return &menuService{
		UserService: userService,
		ClassifyService:classifyService,
	}
}

type menuService struct {
	UserService UserService
	ClassifyService ClassifyService
}

// 获取菜单
func (s *menuService) GetFoodList(businessId, userId int) (int, map[string][]model.FoodResponse, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	var list []model.Food
	err := manager.DBEngine.Where(fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).
		Find(&list)
	if err != nil {
		logrus.Errorf("获取食物失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取食物失败")
	}
	responseList := model.ConvertFoodResponseList(list)
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

	status, foodMap, err := s.classifyFood(businessId,responseList)
	if err != nil{
		return status,  nil, err
	}

	return status, foodMap, nil
}

// 获取单个食物
func (s *menuService) GetFood(businessId, userId, foodId int) (int, *model.FoodResponse, error) {
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
	itemResponse := &model.FoodResponse{
		Food:*item,
	}

	shoppingCart := new(model.ShoppingCart)
	_, err = manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.ColumnUserId), businessId, userId).
		Desc(constant.NameID).Get(shoppingCart)
	if err != nil{
		logrus.Errorf("获取购物车失败: %s", err)
		return iris.StatusInternalServerError,  nil, errors.New("获取购物车失败")
	}

	// 设置已选择的数量
	shoppingCartFoodList := shoppingCart.FoodList
	for _, cartItem := range shoppingCartFoodList {
		if cartItem.Id == itemResponse.Id {
			itemResponse.SelectedCount = cartItem.Num
		}
	}

	status, _, err := s.ClassifyService.GetClassify(businessId,1)
	if err != nil {
		return status,nil,err
	}
	return iris.StatusOK,itemResponse , nil
}

// 添加食物
func (s *menuService) InsertFoodOne(food *model.Food) (int, error) {

	if food.Name == ""{
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
		if subItem.Name == ""{
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
	if food.Id == 0 || food.Name == ""{
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

// 卖出食物
func (s *menuService) SellFood(businessId, userId, foodId, num int) (int, error) {
	status, item, err := s.GetFood(businessId, userId, foodId)
	if err != nil {
		return status,err
	}
	item.Num -= num
	item.SaleCount += num
	status, err = s.UpdateFood(item.GetFood())
	if err != nil {
		return status,err
	}
	return status, nil
}

// 删除食物
func (s *menuService) DeleteFood(businessId, foodId int) (int, error) {
	if foodId == 0 {
		return iris.StatusBadRequest, errors.New("食物id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		businessId,foodId).Delete(new(model.Food))
	if err != nil {
		logrus.Errorf("删除食物失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除食物失败")
	}
	return iris.StatusOK, nil
}

// 获取收藏食物
func (s *menuService) GetCollectList(userId, businessId int) (int,[]model.FoodResponse, error){
	if userId == 0 {
		return iris.StatusBadRequest, nil, errors.New("用户id不能为空")
	}

	status, item, err := getCollectFood(userId,businessId)
	if err != nil {
		return status,nil, err
	}
	if item == nil {
		return iris.StatusOK, []model.FoodResponse{}, nil
	}

	status, foodMap, err := s.GetFoodList(businessId,userId)
	if err != nil {
		return status,nil, err
	}
	foodList := make([]model.FoodResponse,0)

	for _, value := range foodMap {
		foodList = append(foodList,value...)
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

// 分类食物
func (s *menuService)classifyFood(businessId int, list []model.FoodResponse)(int, map[string][]model.FoodResponse, error){


	foodMap := make(map[string][]model.FoodResponse)

	// 设置分类
	for _, subItem := range list {
		ids := strings.Split(subItem.ClassifyId,",")
		for _, classifyId := range ids {
			var key string
			if classifyId == ""{
				key = "未分类"
			}else {
				classifyIdInt,_ := strconv.ParseInt(classifyId,10,16)
				status, classify, err := s.ClassifyService.GetClassify(businessId, int(classifyIdInt))
				if err != nil {
					return status, nil, err
				}
				key = classify.Name
			}
			foodMap[key] = append(foodMap[key],subItem)
		}
	}
	return iris.StatusOK,foodMap,nil
}



