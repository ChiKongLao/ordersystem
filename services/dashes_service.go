package services

import (
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/datamodels"
	"errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/chikong/ordersystem/constant"
)

type DashesService interface {
	GetDashesList(businessId string) (int, []datamodels.Dashes, error)
	GetDashes(businessId, dashId string) (int, *datamodels.Dashes, error)
	InsertDashes(user *datamodels.User,dashes *datamodels.Dashes) (int, error)
	UpdateDashes(user *datamodels.User,dashes *datamodels.Dashes) (int, error)
	DeleteDashes(user *datamodels.User, dashId int) (int, error)
}

func NewDashesService() DashesService {
	return &dashesService{}
}

type dashesService struct {
}

// 获取菜单
func (s *dashesService) GetDashesList(businessId string) (int, []datamodels.Dashes, error) {
	if businessId == ""{
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]datamodels.Dashes,0)

	err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?",constant.ColumnBusinessId),businessId).Find(&list)
	if err != nil{
		logrus.Errorf("获取菜式失败: %s",err)
		return iris.StatusInternalServerError, nil, err
	}

	return iris.StatusOK, list, nil
}

// 获取单个菜式
func (s *dashesService) GetDashes(businessId, dashId string) (int, *datamodels.Dashes, error) {
	if businessId == ""{
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}
	if dashId == ""{
		return iris.StatusBadRequest, nil, errors.New("菜式id不能为空")
	}
	item := new(datamodels.Dashes)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?",constant.ColumnBusinessId,constant.NameID),businessId,dashId).Get(item)
	if err != nil{
		logrus.Errorf("获取菜式失败: %s",err)
		return iris.StatusInternalServerError, nil, errors.New("获取菜式失败")
	}
	if res == false{
		logrus.Errorf("菜式不存在: %s",dashId)
		return iris.StatusNotFound, nil, errors.New("菜式不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加菜式
func (s *dashesService) InsertDashes(user *datamodels.User,dashes *datamodels.Dashes) (int, error) {

	if dashes.Name == "" || dashes.Price == ""{
		return iris.StatusBadRequest,errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.InsertOne(dashes)
	if err != nil{
		logrus.Errorf("添加菜式失败: %s",err)
		return iris.StatusInternalServerError,err
	}
	return iris.StatusOK, nil
}

// 修改菜式
func (s *dashesService) UpdateDashes(user *datamodels.User,dashes *datamodels.Dashes) (int, error) {
	if dashes.Id == 0 || dashes.Name == "" || dashes.Price == ""{
		return iris.StatusBadRequest,errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?",constant.ColumnBusinessId,constant.NameID),
		user.Id,dashes.Id).Update(dashes)
	if err != nil{
		logrus.Errorf("修改菜式失败: %s",err)
		return iris.StatusInternalServerError,err
	}
	return iris.StatusOK, nil
}

// 删除菜式
func (s *dashesService) DeleteDashes(user *datamodels.User, dashId int) (int, error) {
	if dashId == 0 {
		return iris.StatusBadRequest,errors.New("菜式id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?",constant.ColumnBusinessId,constant.NameID),
			user.Id,dashId).Delete(new (datamodels.Dashes))
	if err != nil{
		logrus.Errorf("删除菜式失败: %s",err)
		return iris.StatusInternalServerError,err
	}
	return iris.StatusOK, nil
}

