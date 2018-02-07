package services

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/datamodels"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"time"
	"fmt"
	"github.com/chikong/ordersystem/constant"
)

type UserService interface {
	InsertUser(int, string, string, string) (int, error)
	Login(context iris.Context,userName, password string) (int, string, error)
	GetUserByName(userName string) (int, *datamodels.User, error)
	GetUserById(string) (int, *datamodels.User, error)
	GetBusinessById(string) (int, *datamodels.User, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

func NewUserService() UserService {
	return &userService{
	}
}

type userService struct {
}

// 注册
func (s *userService) InsertUser(role int, userName, password, nickName string) (int, error) {
	if userName == "" || password == ""{
		return iris.StatusBadRequest,errors.New("用户名或密码不能为空")
	}
	if nickName == ""{
		return iris.StatusBadRequest,errors.New("昵称不能为空")
	}
	 user := &datamodels.User{
		UserName:    userName,
		Password:    password,
		NickName:    nickName,
		Role:        role,
		CreatedTime: time.Now().Unix(),
	}
	_, err := manager.DBEngine.InsertOne(*user)
	if err != nil{
		if strings.Contains(err.Error(),"Duplicate entry") {
			return iris.StatusBadRequest,errors.New("用户已存在")
		}
		return iris.StatusInternalServerError,err
	}
	return iris.StatusOK, nil
}

// 登录
func (s *userService) Login(ctx iris.Context, userName, password string) (int, string, error) {
	if userName == "" || password == ""{
		return iris.StatusBadRequest,"",errors.New("用户名或密码不能为空")
	}

	status, user, err := s.GetUserByName(userName)
	if err != nil{
		return status,"",errors.New("没有找到该用户")
	}
	if user.Password != password {
		return iris.StatusBadRequest,"",errors.New("密码不正确")
	}
	token, err := authentication.MakeToken(user)
	if err != nil {
		return iris.StatusInternalServerError,"",err
	}
	setToken(user,token)

	return iris.StatusOK,token,nil

}



// 查询
func (s *userService)GetUserByName(userName string) (int, *datamodels.User, error) {
	user := new(datamodels.User)
	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?",constant.ColumnUserName),userName).Get(user)
	if err != nil{
		logrus.Errorf("查找用户失败:%s",err)
		return iris.StatusInternalServerError, nil,errors.New("查找用户失败")
	}
	if res == false{
		return iris.StatusNotFound,nil,errors.New("没有找到该用户")
	}
	return iris.StatusOK,user,nil

}

// 查询商家
func (s *userService)GetBusinessById(id string) (int, *datamodels.User, error) {
	status, user, err := s.GetUserById(id)
	if err != nil {
		return status,nil,err
	}
	if !user.IsBusiness() {
		return iris.StatusNotFound,nil,errors.New("没有找到该用户")
	}
	return status,user,nil
}

// 查询
func (s *userService)GetUserById(id string) (int, *datamodels.User, error) {
	user := new(datamodels.User)
	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?",constant.NameID),id).Get(user)
	if err != nil{
		logrus.Errorf("查找用户失败:%s",err)
		return iris.StatusInternalServerError, nil,errors.New("查找用户失败")
	}
	if res == false{
		return iris.StatusNotFound,nil,errors.New("没有找到该用户")
	}
	return iris.StatusOK,user,nil
}

// 设置token
func setToken(user *datamodels.User, token string){
	user.Token = token
	_, err := manager.DBEngine.Id(user.Id).Update(user)
	if err != nil {
		logrus.Errorf("更新用户token失败: %s",err)
		return
	}

}



func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *userService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
