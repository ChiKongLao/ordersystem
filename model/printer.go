package model

// 打印机
type Printer struct {
	Id          int     `json:"id" xorm:"not null pk autoincr unique INT"`
	Name        string  `json:"name" xorm:"not null VARCHAR(255)"`
	//Id          int     `json:"id"`
	//Name        string  `json:"name"`
	BusinessId  int     `json:"-" xorm:"not null index INT(11)"`

}