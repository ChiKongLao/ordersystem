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

type ClassifyService interface {
	GetClassifyList(businessId int) (int, []model.Classify, error)
	GetClassify(businessId, classifyId int) (int, *model.Classify, error)
	InsertClassify(name string, businessId, sort int) (int, error)
	UpdateClassify(name string, businessId, classifyId, sort int) (int, error)
	DeleteClassify(businessId, classifyId int) (int, error)
}

func NewClassifyService(userService UserService) ClassifyService {
	return &classifyService{
		UserService: userService,
	}
}

type classifyService struct {
	UserService UserService
}

// 获取分类列表
func (s *classifyService) GetClassifyList(businessId int) (int, []model.Classify, error) {

	list := make([]model.Classify, 0)

	err := manager.DBEngine.
		Where(fmt.Sprintf("%s=?", constant.ColumnBusinessId), businessId).
		OrderBy(constant.ColumnSort).
		Find(&list)
	if err != nil {
		logrus.Errorf("获取分类失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取分类失败")
	}

	return iris.StatusOK, list, nil
}

// 获取单个分类
func (s *classifyService) GetClassify(businessId, classifyId int) (int, *model.Classify, error) {

	item := new(model.Classify)

	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?",
			constant.ColumnBusinessId, constant.NameID),businessId, classifyId).Get(item)
	if err != nil {
		logrus.Errorf("获取分类详情失败: %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取分类详情失败")
	}
	if res == false {
		logrus.Errorf("分类不存在: %s", classifyId)
		return iris.StatusNotFound, nil, errors.New("分类不存在")
	}

	return iris.StatusOK, item, nil
}

// 添加分类
func (s *classifyService) InsertClassify(name string, businessId, sort int) (int, error) {
	if name == "" {
		return iris.StatusBadRequest, errors.New("分类名不能为空")
	}

	_, err := manager.DBEngine.InsertOne(&model.Classify{
		BusinessId: businessId,
		Name:       name,
		Sort:       sort,
	})
	if err != nil {
		logrus.Errorf("添加分类失败: %s", err)
		return iris.StatusInternalServerError, errors.New("添加分类失败")
	}
	return iris.StatusOK, nil
}

// 修改分类
func (s *classifyService) UpdateClassify(name string, businessId, classifyId, sort int) (int, error) {
	if name == "" {
		return iris.StatusBadRequest, errors.New("分类名不能为空")
	}
	status, dbItem, err := s.GetClassify(businessId,classifyId)
	if err != nil {
		return status, err
	}
	// 设置修改信息
	dbItem.Name = name
	dbItem.Sort = sort

	_, err = manager.DBEngine.AllCols().Where(
		fmt.Sprintf("%s=? and %s=?",
			constant.ColumnBusinessId, constant.NameID),businessId, classifyId).Update(dbItem)
	if err != nil {
		logrus.Errorf("修改分类失败: %s", err)
		return iris.StatusInternalServerError, errors.New("修改分类失败")
	}

	return iris.StatusOK, nil
}

// 删除分类
func (s *classifyService) DeleteClassify(businessId, classifyId int) (int, error) {
	_, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=? and %s=?",
			constant.ColumnBusinessId, constant.NameID),businessId, classifyId).Delete(new(model.Classify))
	if err != nil {
		logrus.Errorf("删除分类失败: %s", err)
		return iris.StatusInternalServerError, errors.New("删除分类失败")
	}
	return iris.StatusOK, nil
}
