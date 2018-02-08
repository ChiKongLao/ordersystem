package model

// 餐桌
type TableInfo struct {
	Id         int    `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	Name       string `json:"name" xorm:"not null VARCHAR(255)"`
	Price      string `json:"price" xorm:"TINYTEXT"`                      // 消费总额
	BusinessId int    `json:"businessId" xorm:"not null index INT(11)"`   // 商家id
	Time       string `json:"time" xorm:"VARCHAR(15)"`                    // 就餐时间
	Status     int    `json:"status" xorm:"not null default 0 INT(11)"`   // 餐桌状态
	PersonNum  int    `json:"personNum" xorm:"default 0 INT(11)"`         // 就餐人数
	Capacity   int    `json:"capacity" xorm:"not null default 2 INT(11)"` // 容纳人数
	OrderList  []Dashes `json:"orderList" `                                       // 订单列
}
