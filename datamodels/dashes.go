package datamodels

const (
	NameNum   = "num"
	NamePic   = "pic"
	NamePrice = "price"
	NameType  = "type"
)

// 菜式
type Dashes struct {
	Id         int    `xorm:"not null pk autoincr unique INT"` // 菜式id
	BusinessId int    `xorm:"not null index INT(11)"`          // 商家id
	Name       string `xorm:"not null VARCHAR(255)"`
	Num        int    `xorm:"not null INT"`
	Pic        string `xorm:"VARCHAR(255)"`
	Price      string `xorm:"not null VARCHAR(20)"`
	Type       string `xorm:"VARCHAR(20)"` 					   // 种类
}
