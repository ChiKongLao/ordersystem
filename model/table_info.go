package model

import "github.com/chikong/ordersystem/constant"

// 餐桌
type TableInfo struct {
	Id         int    `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	Name       string `json:"name" xorm:"not null VARCHAR(255)"`
	Price      string `json:"price" xorm:"TINYTEXT"`                      // 消费总额
	BusinessId int    `json:"-" xorm:"not null index INT(11)"`            // 商家id
	Time       int64  `json:"time" xorm:"not null BIGINT(20)"`                                       // 就餐时间
	Status     int    `json:"status" xorm:"not null default 0 INT(11)"`   // 餐桌状态
	PersonNum  int    `json:"personNum" xorm:"default 0 INT(11)"`         // 就餐人数
	Capacity   int    `json:"capacity" xorm:"not null default 2 INT(11)"` // 容纳人数
	UserId     string `json:"userId" xorm:"VARCHAR(255)"`                 // 用户id列
	//Desc       string  `json:"desc" xorm:"VARCHAR(255)"`                   // 描述
	OrderList []Order `json:"orderList" ` // 订单列
}

// 重置餐桌信息
func (item *TableInfo) ClearTable() {
	item.PersonNum = 0
	item.Price = ""
	item.UserId = ""
	item.Time = 0
	item.OrderList = nil
}

// 获取餐桌状态描述
func (item *TableInfo) GetStatusString() string {
	switch item.Status {
	case constant.TableStatusEmpty:
		return "闲置"
	case constant.TableStatusUsing:
		return "正在使用"
	case constant.TableStatusWaitClean:
		return "待清理"
	case constant.TableStatusCleaning:
		return "清理中"
	case constant.TableStatusOrdering:
		return "点餐中"
	default:
		return "未知"
	}
}
