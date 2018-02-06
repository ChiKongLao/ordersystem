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
)

type UserService interface {
	InsertUser(userName, password, nickName string) (int, error)
	Login(context iris.Context,userName, password string) (int, string, error)
	GetUserByName(userName string) (*datamodels.User, error)
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
func (s *userService) InsertUser(userName, password, nickName string) (int, error) {
	if userName == "" || password == ""{
		return iris.StatusBadRequest,errors.New("用户名或密码不能为空")
	}
	if nickName == ""{
		return iris.StatusBadRequest,errors.New("昵称不能为空")
	}
	user := datamodels.NewLoginUser(userName,password,nickName)
	_, err := manager.DBEngine.InsertOne(user)
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

	user,err := s.GetUserByName(userName)
	if err != nil{
		return iris.StatusInternalServerError,"",errors.New("没有找到该用户")
	}
	if user.Password != password {
		return iris.StatusBadRequest,"",errors.New("密码不正确")
	}
	token, err := authentication.MakeToken(userName,password)
	if err != nil {
		return iris.StatusInternalServerError,"",err
	}
	setToken(user,token)

	return iris.StatusOK,token,nil

}



// 查询
func (s *userService)GetUserByName(userName string) (*datamodels.User, error) {
	user := new(datamodels.User)
	res, err := manager.DBEngine.Where("user_name=?",userName).Get(user)
	if err != nil{
		logrus.Errorf("查找用户失败:%s",err)
		return nil,errors.New("查找用户失败")
	}
	if res == false{
		return nil,errors.New("没有找到该用户")
	}
	return user,nil

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
