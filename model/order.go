package model

// 订单
type Order struct {
	Id         int      `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	TableName  string   `json:"name" xorm:"not null VARCHAR(20)"`
	PersonNum  int      `json:"personNum" xorm:"not null INT(11)"`
	Price      int      `json:"price" xorm:"not null INT(11)"`
	Status     int      `json:"status" xorm:"INT(11)"`
	Time       string   `json:"time" xorm:"not null VARCHAR(25)"`
	BusinessId int      `json:"businessId" xorm:"not null index INT(11)"`
	UserId     int      `json:"userId" xorm:"not null index INT(11)"` // 下单的用户id
	DashesList []Dashes `json:"dashesList" `                          // 菜单
}
