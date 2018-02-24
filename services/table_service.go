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

type TableService interface {
	GetTableList(businessId int) (int, []model.TableInfo, error)
	GetTable(businessId, tableId int) (int, *model.TableInfo, error)
	InsertTable(table *model.TableInfo) (int, error)
	UpdateTable(table *model.TableInfo) (int, error)
	DeleteTable(businessId, tableId int) (int, error)
	JoinTable(businessId, tableId int) (int, error)
}

func NewTableService() TableService {
	return &tableService{}
}

type tableService struct {
}

// 获取餐桌列表
func (s *tableService) GetTableList(businessId int) (int, []model.TableInfo, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.TableInfo, 0)

	err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).Find(&list)
	if err != nil {
		logrus.Errorf("获取餐桌失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取餐桌失败")
	}

	return iris.StatusOK, list, nil
}

// 获取单个餐桌
func (s *tableService) GetTable(businessId, tableId int) (int, *model.TableInfo, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if tableId == 0 {
		return iris.StatusBadRequest, nil, errors.New("餐桌id不能为空")
	}
	item := new(model.TableInfo)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID), businessId, tableId).Get(item)
	if err != nil {
		logrus.Errorf("获取餐桌失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取餐桌失败")
	}
	if res == false {
		logrus.Errorf("餐桌不存在: %s", tableId)
		return iris.StatusNotFound, nil, errors.New("餐桌不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加餐桌
func (s *tableService) InsertTable(table *model.TableInfo) (int, error) {

	if table.Name == ""{
		return iris.StatusBadRequest, errors.New("餐桌信息不能为空")
	}
	if table.Capacity <= 0 {
		return iris.StatusBadRequest, errors.New("餐桌容纳人数错误")
	}

	_, err := manager.DBEngine.InsertOne(table)
	if err != nil {
		logrus.Errorf("添加餐桌失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加餐桌失败")
	}
	return iris.StatusOK, nil
}

// 修改餐桌
func (s *tableService) UpdateTable(table *model.TableInfo) (int, error) {
	if table.Id == 0 || table.Name == ""{
		return iris.StatusBadRequest, errors.New("餐桌信息不能为空")
	}
	if table.Capacity <= 0 {
		return iris.StatusBadRequest, errors.New("餐桌可容纳人数错误")
	}
	if table.Capacity < table.PersonNum {
		return iris.StatusBadRequest, errors.New("餐桌人数大于可容纳人数")
	}
	status, dbItem, err := s.GetTable(table.BusinessId,table.Id)
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.Name = table.Name
	dbItem.Capacity = table.Capacity
	if table.PersonNum != 0 {
		dbItem.PersonNum = table.PersonNum
	}
	dbItem.Status = table.Status
	if dbItem.Status == constant.TableStatusEmpty { // 闲置状态清空人数
		dbItem.PersonNum = 0
	}

	_, err = manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		table.BusinessId, table.Id).Update(dbItem)
	if err != nil {
		logrus.Errorf("修改餐桌失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改餐桌失败")
	}

	return iris.StatusOK, nil
}

// 删除餐桌
func (s *tableService) DeleteTable(businessId, tableId int) (int, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, errors.New("商家id不能为空")
	}
	if tableId == 0 {
		return iris.StatusBadRequest, errors.New("餐桌id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		businessId, tableId).Delete(new(model.TableInfo))
	if err != nil {
		logrus.Errorf("删除餐桌失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除餐桌失败")
	}
	return iris.StatusOK, nil
}



/////////////////////// 客户相关

// 使用餐桌
func (s *tableService) JoinTable(businessId, tableId int) (int, error) {
	status, table, err := s.GetTable(businessId,tableId)
	if err != nil{
		return status, err
	}
	if table.Status == constant.TableStatusEmpty{
		table.Status = constant.TableStatusUsing
	}
	status, err = s.UpdateTable(table)
	if err != nil{
		return status, err
	}

	return iris.StatusOK, nil
}





