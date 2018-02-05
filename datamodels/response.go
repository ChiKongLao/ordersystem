package datamodels

// 请求回调载体
type Response struct {
	//IsOk bool	`json:"isOk"`
	Msg  string	`json:"msg"`
}

const (
	KeyIsOk = "isOk"
)

//func NewSuccess(msg string) *Response{
//	return &Response{true,msg}
//}
//func NewSuccess(msg string) *Response{
//	return &Response{true,msg}
//}
//func (*Response) New(isOk bool,msg string) *Response{
//	return &Response{isOk,msg}
//}
func NewSuccess(msg string) *Response{
	return &Response{msg}
}

// 请求失败的回调
func NewErrorResponse(err error) *Response{
	return &Response{err.Error()}

}


