package model

// 收藏的食物
type CollectFood struct {
	Id              int   `json:"id" xorm:"not null pk autoincr unique INT"`
	UserId          int   `json:"userId" xorm:"not null index INT(11)"` // 用户id
	BusinessId      int   `json:"-" xorm:"not null index INT(11)"`      // 商家id
	CollectFoodId []int `json:"foodIds" xorm:"not null"`            // 食物id列
}
