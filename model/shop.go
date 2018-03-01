package model

// 店铺
type Shop struct {
	Id         int    `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	BusinessId int    `json:"business_id" xorm:"not null unique INT(11)"`
	Name       string `json:"name" xorm:"not null VARCHAR(25)"`
	Desc       string `json:"desc" xorm:"not null VARCHAR(255)"`
	Pic        []string `json:"pic" xorm:"TEXT"`
}
