package datamodels

const (
	Name = "name"
)
type Product struct {
	Id       string `json:"id"`			// id
	Name     string `json:"name"`		// 名字
	Desc string `json:"password"`	// 密码
	NickName string `json:"nickName"`   // 昵称
}
