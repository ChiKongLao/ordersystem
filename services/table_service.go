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
	GetTable(businessId, foodId int) (int, *model.TableInfo, error)
	InsertTable(food *model.TableInfo) (int, error)
	UpdateTable(food *model.TableInfo) (int, error)
	DeleteTable(businessId, foodId int) (int, error)
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
func (s *tableService) GetTable(businessId, foodId int) (int, *model.TableInfo, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if foodId == 0 {
		return iris.StatusBadRequest, nil, errors.New("餐桌id不能为空")
	}
	item := new(model.TableInfo)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID), businessId, foodId).Get(item)
	if err != nil {
		logrus.Errorf("获取餐桌失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取餐桌失败")
	}
	if res == false {
		logrus.Errorf("餐桌不存在: %s", foodId)
		return iris.StatusNotFound, nil, errors.New("餐桌不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加餐桌
func (s *tableService) InsertTable(food *model.TableInfo) (int, error) {

	if food.Name == ""{
		return iris.StatusBadRequest, errors.New("餐桌信息不能为空")
	}
	if food.Capacity <= 0 {
		return iris.StatusBadRequest, errors.New("餐桌容纳人数错误")
	}

	_, err := manager.DBEngine.InsertOne(food)
	if err != nil {
		logrus.Errorf("添加餐桌失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加餐桌失败")
	}
	return iris.StatusOK, nil
}

// 修改餐桌
func (s *tableService) UpdateTable(food *model.TableInfo) (int, error) {
	if food.Id == 0 || food.Name == ""{
		return iris.StatusBadRequest, errors.New("餐桌信息不能为空")
	}
	if food.Capacity <= 0 {
		return iris.StatusBadRequest, errors.New("餐桌可容纳人数错误")
	}
	if food.Capacity < food.PersonNum {
		return iris.StatusBadRequest, errors.New("餐桌人数大于可容纳人数")
	}
	status, dbItem, err := s.GetTable(food.BusinessId,food.Id)
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.Name = food.Name
	dbItem.Capacity = food.Capacity
	if food.Status != 0 {
		dbItem.Status = food.Status
	}
	if food.PersonNum != 0 {
		dbItem.PersonNum = food.PersonNum
	}

	_, err = manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		food.BusinessId, food.Id).Update(dbItem)
	if err != nil {
		logrus.Errorf("修改餐桌失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改餐桌失败")
	}

	return iris.StatusOK, nil
}

// 删除餐桌
func (s *tableService) DeleteTable(businessId, foodId int) (int, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, errors.New("商家id不能为空")
	}
	if foodId == 0 {
		return iris.StatusBadRequest, errors.New("餐桌id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		businessId, foodId).Delete(new(model.TableInfo))
	if err != nil {
		logrus.Errorf("删除餐桌失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除餐桌失败")
	}
	return iris.StatusOK, nil
}
