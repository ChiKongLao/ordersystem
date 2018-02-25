package model

// 食物分类
type Classify struct {
	Id         int    `json:"id" xorm:"not null pk autoincr unique INT"`
	BusinessId int    `json:"-" xorm:"not null index INT(11)"`
	Name       string `json:"name" xorm:"not null unique VARCHAR(255)"`
	Sort       int    `json:"sort" xorm:"not null default 100 INT(11)"`
}
