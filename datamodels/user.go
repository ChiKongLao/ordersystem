package datamodels

import "time"

const (
	NameID            = "id"
	NameUserName      = "userName"
	NamePassword      = "password"
	NameNickName      = "nickName"
	NameAuthorization = "authorization"
)

type User struct {
	Id          string `json:"id" xorm:"not null pk autoincr unique INT(11)"`          // id
	UserName    string `json:"userName" xorm:"varchar(25) notnull unique 'user_name'"` // 名字
	Password    string `json:"password" xorm:"not null VARCHAR(30) "`                  // 密码
	NickName    string `json:"nickName" xorm:"VARCHAR(30)"`                            // 昵称
	CreatedTime int64  `json:"createdTime" xorm:"INT(11) "`                            // 创建时间
	Token       string `json:"token" orm:"VARCHAR(255)"`                               // token
}

func NewLoginUser(userName, password, nickName string) *User {
	user := &User{}
	user.UserName = userName
	user.Password = password
	user.NickName = nickName
	user.CreatedTime = time.Now().Unix()
	return user

}
