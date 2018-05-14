package services

import (
	"github.com/chikong/ordersystem/payment"
	"fmt"
)

type PayService interface {
	GetPayClient(notifyUrl string) payment.Pay
}

func NewPayService() PayService {
	service := &payService{}
	return service

}

type payService struct {
}

func (s *payService) init() {

}

func (s *payService)GetPayClient(notifyUrl string) payment.Pay  {
	wxconfig := map[string]interface{}{
		payment.KEY_APP_ID:       payment.WechatAppId,
		payment.KEY_MERCHANT_ID:  payment.WechatMCHId,
		payment.KEY_APP_KEY:      payment.WechatPayKey,
		payment.KEY_PayNotifyUrl: notifyUrl,
	}

	payClient, err := payment.NewPayment().Init(payment.PAYMENTTYPE_WECHATPAY, wxconfig)
	if err != nil {
		fmt.Println("初始化微信支付失败", err)
		payClient = nil
	}
	return payClient
}
