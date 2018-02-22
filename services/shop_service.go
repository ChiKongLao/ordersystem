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

type ShopService interface {
	GetShopList() (int, []model.Shop, error)
	GetShop(businessId int) (int, *model.Shop, error)
	InsertShop(businessId int, name, desc, pic string) (int, error)
	UpdateShop(businessId int, name, desc, pic string) (int, error)
	DeleteShop(businessId int) (int, error)
}

func NewShopService(userService UserService) ShopService {
	return &shopService{
		UserService: userService,
	}
}

type shopService struct {
	UserService UserService
}

// 获取店铺列表
func (s *shopService) GetShopList() (int, []model.Shop, error) {

	list := make([]model.Shop, 0)

	err := manager.DBEngine.Find(&list)
	if err != nil {
		logrus.Errorf("获取店铺失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取店铺失败")
	}

	return iris.StatusOK, list, nil
}

// 获取单个店铺
func (s *shopService) GetShop(businessId int) (int, *model.Shop, error) {

	status,_,err := s.UserService.GetUserById(businessId)
	if err != nil {
		logrus.Errorf("获取店铺详情失败: %s", err)
		return status, nil, err
	}

	item := new(model.Shop)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).Get(item)
	if err != nil {
		logrus.Errorf("获取店铺详情失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取店铺详情失败")
	}
	if res == false {
		logrus.Errorf("店铺不存在: %s", businessId)
		return iris.StatusNotFound, nil, errors.New("店铺不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加店铺
func (s *shopService) InsertShop(businessId int, name, desc, pic string) (int, error) {
	if name == "" {
		return iris.StatusBadRequest, errors.New("店铺名不能为空")
	}

	status, user, err := s.UserService.GetUserById(businessId)
	if err != nil {
		logrus.Errorf("获取商家信息失败: %s", err)
		return status, err
	}
	if !user.IsBusiness() {
		return iris.StatusBadRequest,errors.New("该用户非商家")
	}

	_, err = manager.DBEngine.InsertOne(&model.Shop{
		BusinessId: businessId,
		Name:       name,
		Desc:       desc,
		Pic:        pic,
	})
	if err != nil {
		logrus.Errorf("添加店铺失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加店铺失败")
	}
	return iris.StatusOK, nil
}

// 修改店铺
func (s *shopService) UpdateShop(businessId int, name, desc, pic string) (int, error) {
	if name == "" {
		return iris.StatusBadRequest, errors.New("店铺名不能为空")
	}
	status, user, err := s.UserService.GetUserById(businessId)
	if err != nil {
		logrus.Errorf("获取商家信息失败: %s", err)
		return status, err
	}
	if !user.IsBusiness() {
		return iris.StatusBadRequest,errors.New("该用户非商家")
	}

	status, dbItem, err := s.GetShop(businessId)
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.Name = name
	dbItem.Desc = desc
	dbItem.Pic = pic

	_, err = manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId),
		businessId).Update(dbItem)
	if err != nil {
		logrus.Errorf("修改店铺失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改店铺失败")
	}

	return iris.StatusOK, nil
}

// 删除店铺
func (s *shopService) DeleteShop(businessId int) (int, error) {
	status, user, err := s.UserService.GetUserById(businessId)
	if err != nil {
		logrus.Errorf("获取商家信息失败: %s", err)
		return status, err
	}
	if !user.IsBusiness() {
		return iris.StatusBadRequest,errors.New("该用户非商家")
	}

	_, err = manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId),
		businessId).Delete(new(model.Shop))
	if err != nil {
		logrus.Errorf("删除店铺失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除店铺失败")
	}
	return iris.StatusOK, nil
}
