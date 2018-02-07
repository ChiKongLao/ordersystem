package services

import (
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/datamodels"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/chikong/ordersystem/model"
)

type DashesService interface {
	InsertDashes(user *datamodels.User,dashes *model.Dashes) (int, error)
}

func NewDashesService() DashesService {
	return &dashesService{}
}

type dashesService struct {
}

// 获取菜单
func (s *dashesService) GetDashesList(user *datamodels.User,dashes *model.Dashes) (int,[]model.Dashes, error) {
	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized, nil, errors.New("没有该权限")
	}

	if dashes.Name == "" || dashes.Price == ""{
		return iris.StatusBadRequest, nil, errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.InsertOne(dashes)
	if err != nil{
		logrus.Errorf("添加菜式失败: %s",err)
		return iris.StatusInternalServerError, nil, err
	}
	return iris.StatusOK, nil, nil
}

// 添加菜
func (s *dashesService) InsertDashes(user *datamodels.User,dashes *model.Dashes) (int, error) {
	if !user.IsManager() && !user.IsBusiness(){
		return iris.StatusUnauthorized,errors.New("没有该权限")
	}

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

