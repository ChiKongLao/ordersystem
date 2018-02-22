package model

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/constant"
)

// 请求回调载体
type Response struct {
	//IsOk bool	`json:"isOk"`
	Msg  string	`json:"msg"`
	//Data interface{} `json:"data"`
}


//func NewSuccess(msg string) *Response{
//	return &Response{true,msg}
//}
//func NewSuccess(msg string) *Response{
//	return &Response{true,msg}
//}
//func (*Response) New(isOk bool,msg string) *Response{
//	return &Response{isOk,msg}
//}


// 请求失败的回调
func NewErrorResponse(err error) iris.Map{
	return NewErrorResponseWithMsg(err.Error())
}
// 请求失败的回调
func NewErrorResponseWithMsg(msg string) iris.Map{
	return iris.Map{constant.NameMsg:msg}

}


