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
	"time"
)

type OrderService interface {
	GetOrderList(businessId int) (int, []model.Order, error)
	GetOrder(businessId, orderId int) (int, *model.Order, error)
	InsertOrder(order *model.Order) (int, int, error)
	UpdateOrder(order *model.Order) (int, error)
	DeleteOrder(businessId, orderId int) (int, error)
	GetOldCustomer(businessId int) (int, interface{}, error)
}

func NewOrderService(UserService UserService, MenuService MenuService) OrderService {
	return &orderService{
		MenuService: MenuService,
		UserService:   UserService,
	}
}

type orderService struct {
	MenuService MenuService
	UserService   UserService
}

// 获取订单列表
func (s *orderService) GetOrderList(businessId int) (int, []model.Order, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.Order, 0)

	//err := manager.DBEngine.Where(
	//	fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).Find(&list)
	sql := ""
	err := manager.DBEngine.SQL(sql).Find(&list)
	if err != nil {
		logrus.Errorf("获取订单失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取订单失败")
	}

	return iris.StatusOK, list, nil
}

// 获取单个订单
func (s *orderService) GetOrder(businessId, orderId int) (int, *model.Order, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if orderId == 0 {
		return iris.StatusBadRequest, nil, errors.New("订单id不能为空")
	}
	item := new(model.Order)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID), businessId, orderId).Get(item)
	if err != nil {
		logrus.Errorf("获取订单失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取订单失败")
	}
	if res == false {
		logrus.Errorf("订单不存在: %s", orderId)
		return iris.StatusNotFound, nil, errors.New("订单不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加订单
func (s *orderService) InsertOrder(order *model.Order) (int, int, error) {

	if order.TableName == "" || order.PersonNum == 0 {
		return iris.StatusBadRequest, 0, errors.New("订单信息不能为空")
	}
	order.Time = strconv.FormatInt(time.Now().Unix(), 10)
	order.Status = constant.OrderStatusWaitPay

	// 设置菜单信息
	foodList := order.FoodList
	for i, subItem := range foodList {
		status, dbItem, err := s.MenuService.GetFood(subItem.Id)
		if err != nil {
			return status, 0, err
		}
		subItem.Price = dbItem.Price
		subItem.Name = dbItem.Name
		subItem.Type = dbItem.Type
		foodList[i] = subItem
	}
	order.FoodList = foodList

	sumPrice, err := s.MenuService.GetOrderSumPrice(order.FoodList)
	if err != nil {
		return iris.StatusInternalServerError, 0, err
	}
	order.Price = sumPrice

	_, err = manager.DBEngine.InsertOne(order)

	if err != nil {
		logrus.Errorf("添加订单失败: %s", err)
		return iris.StatusInternalServerError, 0, errors.New("添加订单失败")
	}
	return iris.StatusOK, order.Id, nil
}

// 修改订单
func (s *orderService) UpdateOrder(order *model.Order) (int, error) {
	if order.Id == 0 || order.TableName == "" ||
		order.PersonNum == 0 || order.Status == 0 ||
		order.BusinessId == 0 {
		return iris.StatusBadRequest, errors.New("订单信息不能为空")
	}
	status, dbItem, err := s.GetOrder(order.BusinessId, order.Id)
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.TableName = order.TableName
	dbItem.Status = order.Status
	dbItem.PersonNum = order.PersonNum
	dbItem.Time = strconv.FormatInt(time.Now().Unix(), 10)
	var sumPrice float32
	sumPrice, err = s.MenuService.GetOrderSumPrice(order.FoodList)
	if err != nil {
		return iris.StatusInternalServerError, err
	}
	order.Price = sumPrice

	_, err = manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		order.BusinessId, order.Id).Update(dbItem)
	if err != nil {
		logrus.Errorf("修改订单失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改订单失败")
	}

	return iris.StatusOK, nil
}

// 删除订单
func (s *orderService) DeleteOrder(businessId, orderId int) (int, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, errors.New("商家id不能为空")
	}
	if orderId == 0 {
		return iris.StatusBadRequest, errors.New("订单id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		businessId, orderId).Delete(new(model.Order))
	if err != nil {
		logrus.Errorf("删除订单失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除订单失败")
	}
	return iris.StatusOK, nil
}

// 获取商家的老用户列表
func (s *orderService) GetOldCustomer(businessId int) (int, interface{}, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	type OldCustomer struct {
		NickName string `json:"nickName"`
		Head     string `json:"head"`
		Num      int    `json:"num"`
	}
	var userList []OldCustomer

	err := manager.DBEngine.Table("user").Select("user.nick_name, Count(user.nick_name) AS num").
		Join("INNER", "order", "order.user_id = user.id").
		GroupBy("order.user_id").
		Find(&userList)
	if err != nil {
		logrus.Errorf("获取老用户列表失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取老用户列表失败")
	}
	return iris.StatusOK, userList, nil
}

