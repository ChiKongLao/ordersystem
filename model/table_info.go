package model

// 餐桌
type TableInfo struct {
	Id         int    `xorm:"not null pk autoincr unique INT(11)"`
	Name       string `xorm:"VARCHAR(255)"`
	Price      string `xorm:"TINYTEXT"`
	BusinessId int    `xorm:"index INT(11)"`
	Time       string `xorm:"VARCHAR(15)"`
	Status     int    `xorm:"not null default 0 INT(11)"`
	Person     int    `xorm:"default 0 INT(11)"`
}
