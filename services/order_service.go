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
	"github.com/chikong/ordersystem/util"
	"strconv"
)

type OrderService interface {
	GetOrderList(businessId, tableId, role, status int) (int, *model.OrderListResponse, error)
	GetOrder(orderId int) (int, *model.OrderResponse, error)
	InsertOrder(order *model.Order, shoppingCartId int) (int, int, error)
	UpdateOrder(order *model.Order) (int, error)
	UpdateOrderStatus(orderId, status int) (int, error)
	DeleteOrder(businessId, orderId int) (int, error)
	GetOldCustomer(businessId int) (int, interface{}, error)
}

func NewOrderService(userService UserService, menuService MenuService,
	tableService TableService, shoppingService ShoppingService, printerService PrinterService) OrderService {

	return &orderService{
		MenuService:     menuService,
		UserService:     userService,
		TableService:    tableService,
		ShoppingService: shoppingService,
		PrinterService:  printerService,
	}
}

type orderService struct {
	MenuService     MenuService
	UserService     UserService
	TableService    TableService
	ShoppingService ShoppingService
	PrinterService  PrinterService
}

// 获取订单列表
func (s *orderService) GetOrderList(businessId, tableId, role, status int) (int, *model.OrderListResponse, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.OrderResponse, 0)

	session := manager.DBEngine.Table("`order`").Select("`order`.*,table_info.name AS table_name").
		Join("INNER", "table_info", "`order`.table_id=table_info.id")

	if role == constant.RoleCustomer {
		symbol := "="
		if status == constant.OrderStatusPaid || status == constant.OrderStatusAll {
			symbol = ">="
			if status == constant.OrderStatusAll {
				status = 0
			}
		}
		session = session.Where(fmt.Sprintf("%s=? and `order`.%s=? and `order`.status%s?",
			constant.ColumnTableId, constant.ColumnBusinessId, symbol), tableId, businessId, status)
	} else if role == constant.RoleBusiness {
		session = session.Where(fmt.Sprintf("`order`.%s=?", constant.ColumnBusinessId), businessId)
	}

	err := session.Desc(constant.ColumnCreateTime).Find(&list)

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
		Where("`order`.id=?", orderId).
		Get(item)
	if err != nil {
		logrus.Errorf("获取订单失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取订单失败")
	}
	if res == false {
		logrus.Errorf("订单不存在: %v", orderId)
		return iris.StatusNotFound, nil, errors.New("订单不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加订单
func (s *orderService) InsertOrder(order *model.Order, shoppingCartId int) (int, int, error) {

	if order.TableId == 0 || order.PersonNum == 0 {
		return iris.StatusBadRequest, 0, errors.New("订单信息不能为空")
	}
	order.OrderNo, _ = s.makeOrderNo(order)
	time := util.GetCurrentTime()
	order.CreateTime = time
	order.UpdateTime = time
	order.Status = constant.OrderStatusWaitPay

	status, shopCarItem, err := s.ShoppingService.GetShopping(order.BusinessId, order.UserId, order.TableId)
	if err != nil {
		return status, 0, err
	}
	if len(shopCarItem.FoodList) == 0 {
		return iris.StatusBadRequest, 0, errors.New("购物车为空")
	}
	order.FoodList = shopCarItem.FoodList

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

	//s.ShoppingService.DeleteShopping(order.BusinessId, shopCarItem.Id) // 删除购物车

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
	order.UpdateTime = util.GetCurrentTime()

	// 设置修改信息
	dbItem.TableId = order.TableId
	dbItem.Status = order.Status
	dbItem.PersonNum = order.PersonNum
	dbItem.UpdateTime = order.UpdateTime
	var sumPrice float32
	sumPrice, err = s.MenuService.GetOrderSumPrice(order.FoodList)
	if err != nil {
		return iris.StatusInternalServerError, err
	}
	order.Price = sumPrice
	dbItem.Price = sumPrice

	if order.Status == constant.OrderStatusPaid { // 订单已付款,减少库存
		status, err = s.handlePaidOrder(dbItem)
		if err != nil {
			return status, err
		}
	}
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
	userList := make([]OldCustomer, 0)
	err := manager.DBEngine.Table("user").Select("user.nick_name, Count(user.nick_name) AS num").
		Join("INNER", "`order`", "`order`.user_id = user.id").
		Where("`order`.business_id=?", businessId).
		GroupBy("`order`.user_id").
		Find(&userList)
	if err != nil {
		logrus.Errorf("获取老用户列表失败: %s", err)
		return iris.StatusInternalServerError, userList, errors.New("获取老用户列表失败")
	}
	return iris.StatusOK, userList, nil
}

// 处理已支付的订单
func (s *orderService) handlePaidOrder(order *model.OrderResponse) (int, error) {
	foodList := order.FoodList
	for _, subItem := range foodList {
		status, err := s.MenuService.SellFood(order.BusinessId, order.UserId, subItem.Id, subItem.Num)
		if err != nil {
			return status, err
		}
	}
	go func() {
		_, orderUser, _ := s.UserService.GetUserById(order.UserId)
		_, businessUser, _ := s.UserService.GetUserById(order.BusinessId)
		network.SendChatMessage("我已经下单啦", orderUser, order.BusinessId, order.TableId)
		network.SendOrderMessage(order.BusinessId, order)

		s.PrinterService.SendOrder(model.OrderPrint{
			OrderResponse: *order,
			Customer:      *orderUser,
			Business:      *businessUser,
		})

	}()

	return iris.StatusOK, nil
}

// 生成订单号
func (s *orderService) makeOrderNo(order *model.Order) (string, error) {
	type DBItem struct {
		OrderNo    string
		CreateTime int64
	}
	item := new(DBItem)
	if res, err := manager.DBEngine.Table("`order`").
		Select("`order`.order_no,'order'.create_time").
		Where(fmt.Sprintf("%s=?", constant.ColumnBusinessId), order.BusinessId).
		Desc(constant.NameID).
		Limit(1).
		Get(&item); res == false || err != nil{
			return "1",nil
	}else{
		no, err := strconv.Atoi(item.OrderNo)
		if err != nil {
			logrus.Errorf("生成订单号失败:%s",err)
			return "1", nil
		}

		return strconv.Itoa(no+1),nil

	}

	return "1", nil
}
