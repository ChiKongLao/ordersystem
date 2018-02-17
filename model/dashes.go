package model

// 菜式
type Dashes struct {
	Id          int    `json:"id" xorm:"not null pk autoincr unique INT"` // 菜式id
	BusinessId  int    `json:"-" xorm:"not null index INT(11)"`           // 商家id
	Name        string `json:"name" xorm:"not null VARCHAR(255)"`
	Num         int    `json:"num" xorm:"not null INT"`
	IsRecommend bool   `json:"isRecommend" xorm:"not null default 0 TINYINT(1)"` // 是否推荐
	Pic         string `json:"pic" xorm:"VARCHAR(255)"`
	Price       string `json:"price" xorm:"not null VARCHAR(20)"`
	Type        string `json:"type" xorm:"VARCHAR(20)"` // 种类
	Desc        string `json:"desc" xorm:"VARCHAR(255)"`
	//Desc       string `json:"desc,omitempty" xorm:"VARCHAR(255)"`
}
