package services

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/chikong/ordersystem/manager"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/datamodels"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
)

type UserService interface {
	InsertUser(userName string, password string) (int, error)
	Login(userName, password string) (int, string, error)
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
func (s *userService) InsertUser(userName, password string) (int, error) {
	if userName == "" || password == ""{
		return iris.StatusBadRequest,errors.New("用户名或密码不能为空")
	}
	user := datamodels.NewLoginUser(userName,password)
	_, err := manager.DBEngine.InsertOne(user)
	if err != nil{
		if strings.Contains(err.Error(),"Duplicate entry") {
			return iris.StatusBadRequest,errors.New("用户已存在")
		}
		return iris.StatusInternalServerError,err
	}

	return iris.StatusOK, nil
}

// 查询
func getUserByName(userName string) (*datamodels.User, error) {
	user := new(datamodels.User)
	res, err := manager.DBEngine.Where("user_name=?",userName).Get(user)
	if err != nil{
		logrus.Errorf("查找用户失败:%s",err)
		return nil,err
	}
	if res == false{
		return nil,errors.New("没有找到该用户")
	}
	return user,nil

}
// 登录
func (s *userService) Login(userName,password string) (int,string, error) {
	user,err := getUserByName(userName)
	if err != nil{
		return iris.StatusInternalServerError,"",errors.New("没有找到该用户")
	}
	if user.Password != password {
		return iris.StatusBadRequest,"",errors.New("密码不正确")
	}
	return iris.StatusOK,user.Token,nil

}



func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *userService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
