package services

import (
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/model"
	"errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/util"
	"github.com/chikong/ordersystem/network"
)

type TableService interface {
	GetTableList(businessId ,status int) (int, []model.TableInfo, error)
	GetTable(businessId, tableId int) (int, *model.TableInfo, error)
	InsertTable(table *model.TableInfo) (int, error)
	UpdateTable(table *model.TableInfo) (int, error)
	DeleteTable(businessId, tableId int) (int, error)
	ChangeTable(businessId, oldTableId, newTableId int) (int, error)

	UpdateTableStatus(businessId, userId, tableId, tableStatus, personNum int) (int, error)
}

func NewTableService(userService UserService) TableService {
	return &tableService{
		UserService:  userService,
	}
}

type tableService struct {
	UserService  UserService
}

// 获取餐桌列表
func (s *tableService) GetTableList(businessId, status int) (int, []model.TableInfo, error) {
	if businessId == 0 {
		return iris.StatusBadRequest, nil, errors.New("商家id不能为空")
	}

	list := make([]model.TableInfo, 0)
	var err error
	if status == constant.TableStatusUnknown {
		err = manager.DBEngine.Where(
			fmt.Sprintf("%s=?", constant.ColumnBusinessId),
			businessId).Find(&list)
	}else{
		err = manager.DBEngine.Where(
			fmt.Sprintf("%s=? and %s=?", constant.ColumnBusinessId,constant.NameStatus),
			businessId,status).Find(&list)
	}

	if err != nil {
		logrus.Errorf("获取餐桌失败: %s", err)
		return iris.StatusInternalServerError, list, errors.New("获取餐桌失败")
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
	dbItem.Time = table.Time
	dbItem.UserId = table.UserId
	//if table.PersonNum != 0 {
		dbItem.PersonNum = table.PersonNum
	//}
	dbItem.Status = table.Status
	if dbItem.Status == constant.TableStatusEmpty { // 闲置状态清空旧的信息
		dbItem.ClearTable()

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


// 换桌
func (s *tableService) ChangeTable(businessId, oldTableId, newTableId int) (int, error) {
	status, oldTable, err := s.GetTable(businessId,oldTableId) // 旧桌
	if err != nil {
		return status,nil
	}
	oldTable.Status = constant.TableStatusWaitClean
	status, err = s.UpdateTable(oldTable)
	if err != nil {
		return status,nil
	}

	status, newTable, err := s.GetTable(businessId,newTableId) // 新桌
	if err != nil {
		return status,nil
	}
	newTable.Status = constant.TableStatusUsing
	newTable.UserId = oldTable.UserId
	newTable.PersonNum = oldTable.PersonNum
	newTable.Time = util.GetCurrentTime()
	status, err = s.UpdateTable(newTable)
	if err != nil {
		return status,nil
	}

	return iris.StatusOK, nil
}





/////////////////////// 客户相关


// 更新餐桌状态
func (s *tableService) UpdateTableStatus(businessId, userId, tableId, tableStatus, personNum int) (int, error) {
	status, table, err := s.GetTable(businessId,tableId)
	if err != nil{
		return status, err
	}
	table.Status = tableStatus
	table.PersonNum = personNum
	isExist := false
	for _, subItem := range table.UserId {
		if subItem == userId {
			isExist = true
			break
		}
	}
	if !isExist { // 餐桌不存在该客户, 则添加
		table.UserId = append(table.UserId, userId)
	}

	status, err = s.UpdateTable(table)
	if err != nil{
		return status, err
	}
	if tableStatus == constant.TableStatusWaitClean {
		_, user, _ := s.UserService.GetUserById(userId)
		network.SendChatMessage("麻烦清理下餐桌",user, businessId, tableId)
	}

	return iris.StatusOK, nil
}





