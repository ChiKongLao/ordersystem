package services

import (
	"github.com/chanxuehong/rand"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chikong/ordersystem/constant"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
	"github.com/kataras/iris/core/errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/chikong/ordersystem/model"
	"github.com/kataras/iris"
)

type WechatService interface {
	GetAuthUrl() string
	GetUserInfo(code, state string) (interface{}, error)
}

func NewWechatService(userService UserService) WechatService {
	oauth2Client := oauth2.Client{
		Endpoint: mpoauth2.NewEndpoint(constant.WechatAppId, constant.WechatAppSecret),
	}
	service := &wechatService{
		UserService:userService,
		oauth2Client: oauth2Client,
	}
	return service
}

type wechatService struct {
	UserService  UserService
	oauth2Client oauth2.Client
}

// 获取个人信息url
func (s *wechatService) GetAuthUrl() string {
	state := string(rand.NewHex())
	return mpoauth2.AuthCodeURL(constant.WechatAppId, constant.WechatOauth2RedirectURI, constant.WechatOauth2Scope, state)
}

// 获取换取token的url
func (s *wechatService) GetUserInfo(code, state string) (interface{}, error) {
	if code == "" {
		return nil, errors.New("用户禁止授权")
	}
	if state == "" {
		return nil, errors.New("state 参数为空")
	}
	token, err := s.oauth2Client.ExchangeToken(code)
	//logrus.Debugf("解析微信用户授权 code=%s, state=%s, token=%v",code,state,token)
	if err != nil {
		str := fmt.Sprintf("获取用户token失败:%s",err.Error())
		logrus.Errorln(str)
		return nil, errors.New(str)
	}

	userInfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		str := fmt.Sprintf("获取用户信息失败:%s",err.Error())
		logrus.Errorln(str)
		return nil, errors.New(str)
	}
	user, err := s.exchangeUser(*userInfo)
	if err != nil {
		return nil, err

	}
	return iris.Map{
		constant.NameID:user.Id,
		constant.NameNickName:user.NickName,
		constant.NameHead:user.Head,
		constant.NameAuthorization:user.Token,
	},nil

}

// 从微信用户信息转为系统用户
func (s *wechatService)exchangeUser(wxUser mpoauth2.UserInfo) (*model.User, error) {
	password := "8888888888"
	s.UserService.InsertUser(constant.RoleCustomer,wxUser.OpenId,password,wxUser.Nickname,wxUser.HeadImageURL)

	_, user, err := s.UserService.Login(wxUser.OpenId,password)
	if err != nil {
		return nil, err
	}
	return user, nil

}