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

type MenuService interface {
	GetDashesList(businessId int) (int, []model.Dashes, error)
	GetDashes(dashId int) (int, *model.Dashes, error)
	InsertDashesOne(dashes *model.Dashes) (int, error)
	InsertDashes(dashes []*model.Dashes) (int, error)
	UpdateDashes(dashes *model.Dashes) (int, error)
	DeleteDashes(dashId int) (int, error)
	GetOrderSumPrice(dashesList []model.Dashes) (float32, error)
}

func NewMenuService() MenuService {
	return &menuService{}
}

type menuService struct {
}

// 获取菜单
func (s *menuService) GetDashesList(businessId int) (int, []model.Dashes, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.Dashes, 0)

	err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).Find(&list)
	if err != nil {
		logrus.Errorf("获取菜式失败: %s", err)
		return iris.StatusInternalServerError, nil, err
	}

	return iris.StatusOK, list, nil
}

// 获取单个菜式
func (s *menuService) GetDashes(dashId int) (int, *model.Dashes, error) {
	if dashId == 0 {
		return iris.StatusBadRequest, nil, errors.New("菜式id不能为空")
	}
	item := new(model.Dashes)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID), dashId).Get(item)
	if err != nil {
		logrus.Errorf("获取菜式失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取菜式失败")
	}
	if res == false {
		logrus.Errorf("菜式不存在: %s", dashId)
		return iris.StatusNotFound, nil, errors.New("菜式不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加菜式
func (s *menuService) InsertDashesOne(dashes *model.Dashes) (int, error) {

	if dashes.Name == "" || dashes.Price == 0 {
		return iris.StatusBadRequest, errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.InsertOne(dashes)
	if err != nil {
		logrus.Errorf("添加菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加菜式失败")
	}
	return iris.StatusOK, nil
}

// 添加菜式
func (s *menuService) InsertDashes(list []*model.Dashes) (int, error) {

	for i, subItem := range list  {
		if subItem.Name == "" || subItem.Price == 0 {
			return iris.StatusBadRequest, errors.New(
				fmt.Sprintf("菜式信息不能为空: %s",i))
		}
	}
	_, err := manager.DBEngine.Insert(list)
	if err != nil {
		logrus.Errorf("添加菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加菜式失败")
	}
	return iris.StatusOK, nil
}

// 修改菜式
func (s *menuService) UpdateDashes(dashes *model.Dashes) (int, error) {
	if dashes.Id == 0 || dashes.Name == "" || dashes.Price == 0 {
		return iris.StatusBadRequest, errors.New("菜式信息不能为空")
	}

	_, err := manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId, constant.NameID),
		dashes.BusinessId, dashes.Id).Update(dashes)
	if err != nil {
		logrus.Errorf("修改菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改菜式失败")
	}
	return iris.StatusOK, nil
}

// 删除菜式
func (s *menuService) DeleteDashes(dashId int) (int, error) {
	if dashId == 0 {
		return iris.StatusBadRequest, errors.New("菜式id不能为空")
	}

	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID),dashId).Delete(new(model.Dashes))
	if err != nil {
		logrus.Errorf("删除菜式失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除菜式失败")
	}
	return iris.StatusOK, nil
}


// 计算订单总价
func (s *menuService)GetOrderSumPrice(dashesList []model.Dashes) (float32, error) {
	var sum float32
	for _, item := range dashesList {
		sum += item.Price * float32(item.Num)
	}
	return sum, nil

}



