package services

import (
	"github.com/chanxuehong/rand"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chikong/ordersystem/constant"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/kataras/iris/core/errors"
	"fmt"
)

type WechatService interface {
	GetAuthUrl() string
	GetUserInfo(code, state string) (string, error)

}

func NewWechatService() WechatService {
	service := &wechatService{}
	service.oauth2Endpoint = mpoauth2.NewEndpoint(constant.WechatAppId,constant.WechatAppSecret)
	return service
}

type wechatService struct {
	UserService  UserService
	TableService TableService
	oauth2Endpoint oauth2.Endpoint
}

// 获取个人信息url
func (s *wechatService)GetAuthUrl() string {
	state := string(rand.NewHex())
	return mpoauth2.AuthCodeURL(constant.WechatAppId, constant.WechatOauth2RedirectURI, constant.WechatOauth2Scope, state)
}

// 获取个人信息
func (s *wechatService)GetUserInfo(code, state string) (string, error) {
	if code == "" {
		return "",errors.New("用户禁止授权")
	}
	if state == "" {
		return "",errors.New("state 参数为空")
	}

	oauth2Client := oauth2.Client{
		Endpoint: s.oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	logrus.Debugf("解析微信用户授权 code=%s, state=%s, token=%v",code,state,token)
	if err != nil {
		str := fmt.Sprintf("获取用户token失败:%s",err.Error())
		logrus.Errorln(str)
		return "", errors.New(str)
	}

	userInfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		str := fmt.Sprintf("获取用户信息失败:%s",err.Error())
		logrus.Errorln(str)
		return "", errors.New(str)
	}

	logrus.Debugln("userinfo: %+v\r\n", userInfo)
	data, _ := json.Marshal(userInfo)

	return string(data),nil
}

