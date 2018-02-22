package model

// 食物
type Food struct {
	Id          int     `json:"id" xorm:"not null pk autoincr unique INT"` // 食物id
	BusinessId  int     `json:"-" xorm:"not null index INT(11)"`           // 商家id
	Name        string  `json:"name" xorm:"not null VARCHAR(255)"`
	Num         int     `json:"num" xorm:"not null INT"`
	IsRecommend bool    `json:"isRecommend" xorm:"not null default 0 TINYINT(1)"` // 是否推荐
	IsCollect   bool    `json:"isCollect"`                                        // 是否收藏
	Pic         string  `json:"pic" xorm:"VARCHAR(255)"`
	Price       float32 `json:"price" xorm:"not null FLOAT"`
	Type        string  `json:"type" xorm:"VARCHAR(20)"` // 种类
	Desc        string  `json:"desc" xorm:"VARCHAR(255)"`
	//Desc       string `json:"desc,omitempty" xorm:"VARCHAR(255)"`
}

// 是否为同一道菜. 包含同样的味道
func (item *Food) IsSameFood(targetItem *Food) bool {
	return item.Id == targetItem.Id &&
		item.Type == targetItem.Type
}
