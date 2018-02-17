package model

import "github.com/chikong/ordersystem/constant"

type User struct {
	Id          int    `json:"id" xorm:"not null pk autoincr unique INT(11)"`          // id
	UserName    string `json:"userName" xorm:"VARCHAR(25) notnull unique 'user_name'"` // 名字
	Password    string `json:"password" xorm:"not null VARCHAR(20) "`                  // 密码
	NickName    string `json:"nickName" xorm:"VARCHAR(30)"`                            // 昵称
	CreatedTime int64  `json:"createdTime" xorm:"INT(11) "`                            // 创建时间
	Token       string `json:"token" orm:"VARCHAR(255)"`                               // token
	Head        string `json:"head" orm:"VARCHAR(255)"`                                // 头像
	Role        int    `xorm:"not null INT(11)"`                                       //角色
}

// 是否为管理员
func (user *User) IsManager() bool {
	return user.Role == constant.RoleManager
}

// 是否为客户
func (user *User) IsCustomer() bool {
	return user.Role == constant.RoleCustomer
}

// 是否为商家
func (user *User) IsBusiness() bool {
	return user.Role == constant.RoleBusiness
}

// 是否为商家或者管理员
func (user *User) IsManagerOrBusiness() bool {
	return user.IsManager() || user.IsBusiness()
}
