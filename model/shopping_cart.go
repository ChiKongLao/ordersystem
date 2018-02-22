package model

// 购物车
type ShoppingCart struct {
	Id         int     `json:"id" xorm:"not null pk autoincr unique INT(11)"`
	UserId     int     `json:"userId" xorm:"not null INT(11)"`
	BusinessId int     `json:"businessId" xorm:"not null INT(11)"`
	FoodList   []Food  `json:"list" xorm:"not null TEXT"`
	Price      float32 `json:"price"`
	Count      int     `json:"count"`
}
