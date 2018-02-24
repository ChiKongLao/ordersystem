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
	GetOrderList(businessId, tableId, role int) (int, *model.OrderListResponse, error)
	GetOrder(orderId int) (int, *model.OrderResponse, error)
	InsertOrder(order *model.Order) (int, int, error)
	UpdateOrder(order *model.Order) (int, error)
	UpdateOrderStatus(orderId, status int) (int, error)
	DeleteOrder(businessId, orderId int) (int, error)
	GetOldCustomer(businessId int) (int, interface{}, error)
}

func NewOrderService(userService UserService, menuService MenuService, tableService TableService) OrderService {
	return &orderService{
		MenuService:  menuService,
		UserService:  userService,
		TableService: tableService,
	}
}

type orderService struct {
	MenuService  MenuService
	UserService  UserService
	TableService TableService
}

// 获取订单列表
func (s *orderService) GetOrderList(businessId, tableId, role int) (int, *model.OrderListResponse, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.OrderResponse, 0)

	session := manager.DBEngine.Table("`order`").Select("`order`.*,table_info.name AS table_name").
		Join("INNER", "table_info", "`order`.table_id=table_info.id")

	if role == constant.RoleCustomer {
		session = session.Where(fmt.Sprintf("%s=?", constant.ColumnTableId), tableId)
	}

	err := session.Find(&list)

	if err != nil {
		logrus.Errorf("获取订单失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取订单失败")
	}
	return iris.StatusOK, model.ConvertOrderResponseData(list), nil
}

// 获取单个订单
func (s *orderService) GetOrder(orderId int) (int, *model.OrderResponse, error) {
	if orderId == 0 {
		return iris.StatusBadRequest, nil, errors.New("订单id不能为空")
	}
	item := new(model.OrderResponse)

	res, err := manager.DBEngine.Table("`order`").Select("`order`.*,table_info.name AS table_name").
		Join("INNER", "table_info", "`order`.table_id = table_info.id").
		GroupBy("`order`.user_id").
		Where("`order`.id=?", orderId).
		Get(item)
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

	if order.TableId == 0 || order.PersonNum == 0 {
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
		if subItem.Num > dbItem.Num {
			return iris.StatusBadRequest,0,errors.New("数量不足")
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
	if order.Id == 0 || order.TableId == 0 ||
		order.PersonNum == 0 || order.Status == 0 ||
		order.BusinessId == 0 {
		return iris.StatusBadRequest, errors.New("订单信息不能为空")
	}
	status, dbItem, err := s.GetOrder(order.Id)
	if err != nil {
		return status, err
	}
	if dbItem.Status == constant.OrderStatusFinish {
		return iris.StatusBadRequest, errors.New("已完成的订单不能修改")
	}
	if order.Status == constant.OrderStatusPaid { // 订单已付款,减少库存
		foodList := order.FoodList
		for _, subItem := range foodList {
			status, err = s.MenuService.ReduceFoodNum(subItem.Id,subItem.Num)
			if err != nil {
				return status, err
			}
		}
	}


	// 设置修改信息
	dbItem.TableId = order.TableId
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
		order.BusinessId, order.Id).Update(model.ConvertOrderResponseToOrder(*dbItem))
	if err != nil {
		logrus.Errorf("修改订单失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改订单失败")
	}

	return iris.StatusOK, nil
}

// 修改订单状态
func (s *orderService) UpdateOrderStatus(orderId, orderStatus int) (int, error) {
	status, dbItem, err := s.GetOrder(orderId)
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.Status = orderStatus
	status, err = s.UpdateOrder(model.ConvertOrderResponseToOrder(*dbItem))
	if err != nil {
		return status, err
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
