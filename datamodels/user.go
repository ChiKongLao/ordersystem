package datamodels

import "time"

const (
	UserName = "userName"
	Password = "password"
	Token = "token"
)

type User struct {
	Id        string    `json:"id"`                        // id
	UserName  string    `json:"userName" xorm:"varchar(25) notnull unique 'user_name'"` // 名字
	Password  string    `json:"password"`                  // 密码
	NickName  string    `json:"nickName"`                  // 昵称
	CreatedAt time.Time `json:"createdTime"`               // 创建时间
	UpdateAt  time.Time `json:"updatedTime"`               // 更新时间
	Token     string    `json:"token"`                     // token
}

func NewLoginUser(userName, password string) *User {
	user := &User{}
	user.UserName = userName
	user.Password = password
	return user

}
