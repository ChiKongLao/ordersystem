package model

// 收藏的菜式
type CollectDishes struct {
	Id              int   `json:"id" xorm:"not null pk autoincr unique INT"`
	UserId          int   `json:"userId" xorm:"not null index INT(11)"` // 用户id
	BusinessId      int   `json:"-" xorm:"not null index INT(11)"`      // 商家id
	CollectDishesId []int `json:"dishesIds" xorm:"not null"`            // 菜式id列
}
