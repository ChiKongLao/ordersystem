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
	GetOrderList(businessId string) (int, []model.Order, error)
	GetOrder(businessId, orderId string) (int, *model.Order, error)
	InsertOrder(order *model.Order) (int, error)
	UpdateOrder(order *model.Order) (int, error)
	DeleteOrder(businessId, orderId string) (int, error)
}

func NewOrderService(DashesService DashesService) OrderService {
	return &orderService{
		DashesService:DashesService,
	}
}

type orderService struct {
	DashesService DashesService
}

// 获取订单列表
func (s *orderService) GetOrderList(businessId string) (int, []model.Order, error) {
	if businessId == "" {
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
func (s *orderService) GetOrder(businessId, orderId string) (int, *model.Order, error) {
	if businessId == "" {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if orderId == "" {
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
func (s *orderService) InsertOrder(order *model.Order) (int, error) {

	if order.TableName == "" || order.PersonNum == 0{
		return iris.StatusBadRequest, errors.New("订单信息不能为空")
	}
	order.Time = strconv.FormatInt(time.Now().Unix(),10)
	order.Status = constant.OrderStatusWaitPay

	// 设置菜单信息
	dashesList := order.DashesList
	for i, subItem := range dashesList {
		status,dbItem, err := s.DashesService.GetDashes(strconv.Itoa(subItem.Id))
		if err != nil{
			return status,err
		}
		subItem.Price = dbItem.Price
		subItem.Name = dbItem.Name
		subItem.Type = dbItem.Type
		dashesList[i] = subItem
	}
	order.DashesList = dashesList

	sumPrice, err := s.getOrderSumPrice(order.BusinessId,order.DashesList)
	if err != nil {
		return iris.StatusInternalServerError,err
	}
	order.Price = sumPrice

	_, err = manager.DBEngine.InsertOne(order)
	if err != nil {
		logrus.Errorf("添加订单失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加订单失败")
	}
	return iris.StatusOK, nil
}

// 修改订单
func (s *orderService) UpdateOrder(order *model.Order) (int, error) {
	if order.Id == 0 || order.TableName == "" ||
		order.PersonNum == 0  || order.Status == 0 ||
			order.BusinessId == 0{
		return iris.StatusBadRequest, errors.New("订单信息不能为空")
	}
	status, dbItem, err := s.GetOrder(strconv.Itoa(order.BusinessId),strconv.Itoa(order.Id))
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.TableName = order.TableName
	dbItem.Status = order.Status
	dbItem.PersonNum = order.PersonNum
	dbItem.Time = strconv.FormatInt(time.Now().Unix(),10)
	var sumPrice float32
	sumPrice, err = s.getOrderSumPrice(order.BusinessId,order.DashesList)
	if err != nil {
		return iris.StatusInternalServerError,err
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
func (s *orderService) DeleteOrder(businessId, orderId string) (int, error) {
	if businessId == "" {
		return iris.StatusBadRequest, errors.New("商家id不能为空")
	}
	if orderId == "" {
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

// 计算订单总价
func (s *orderService)getOrderSumPrice(businessId int, dashesList []model.Dashes) (float32, error){
	priceMap := make(map[int]float32)
	for _, subItem := range dashesList {
		value,ok := priceMap[subItem.Id]
		if ok {
			priceMap[subItem.Id] = value +1
		}else{
			priceMap[subItem.Id] = 1
		}
	}
	var sum float32
	for key, value := range priceMap {
		_, dashes, err := s.DashesService.GetDashes(strconv.Itoa(key))
		if err != nil{
			return 0,err
		}
		res, err :=  strconv.ParseFloat(dashes.Price,10)
		if err !=nil {
			return 0,errors.New("价格格式出错")
		}
		sum += float32(res) * value
	}
	return sum,nil

}
