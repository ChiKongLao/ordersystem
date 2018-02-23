package services

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/model"
	"errors"
	"time"
	"github.com/chikong/ordersystem/manager"
	"strings"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/chikong/ordersystem/constant"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	InsertUser(role int, userName, password, nickName, head string) (int, error)
	Login(userName, password string) (int, *model.User, error)
	GetUserByName(userName string) (int, *model.User, error)
	GetUserById(id int) (int, *model.User, error)
	GetBusinessById(id int) (int, *model.User, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GetUserList() ([]model.User, error)
	CheckRoleIsManagerWithToken(ctx iris.Context) (int, bool, error)
	CheckRoleWithToken(ctx iris.Context, role int) (int, bool, error)
	GetUserFormToken(ctx iris.Context) (int, *model.User, error)
}

func NewUserService() UserService {
	return &userService{
	}
}

type userService struct {
}

// 注册
func (s *userService) InsertUser(role int, userName, password, nickName, head string) (int, error) {
	if userName == "" || password == "" {
		return iris.StatusBadRequest, errors.New("用户名或密码不能为空")
	}
	if nickName == "" {
		return iris.StatusBadRequest, errors.New("昵称不能为空")
	}
	user := &model.User{
		UserName:    userName,
		Password:    password,
		NickName:    nickName,
		Role:        role,
		Head:        head,
		CreatedTime: time.Now().Unix(),
	}
	_, err := manager.DBEngine.InsertOne(*user)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return iris.StatusBadRequest, errors.New("用户已存在")
		}
		return iris.StatusInternalServerError, err
	}
	return iris.StatusOK, nil
}

// 登录
func (s *userService) Login(userName, password string) (int, *model.User, error) {
	if userName == "" || password == "" {
		return iris.StatusBadRequest, nil, errors.New("用户名或密码不能为空")
	}

	status, user, err := s.GetUserByName(userName)
	if err != nil {
		return status, nil, errors.New("没有找到该用户")
	}
	if user.Password != password {
		return iris.StatusBadRequest, nil, errors.New("密码不正确")
	}
	token, err := authentication.MakeToken(user)
	if err != nil {
		return iris.StatusInternalServerError, nil, err
	}
	setToken(user, token)

	return iris.StatusOK, user, nil

}

// 查询
func (s *userService) GetUserByName(userName string) (int, *model.User, error) {
	user := new(model.User)
	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.ColumnUserName), userName).Get(user)
	if err != nil {
		logrus.Errorf("查找用户失败:%s", err)
		return iris.StatusInternalServerError, nil, errors.New("查找用户失败")
	}
	if res == false {
		return iris.StatusNotFound, nil, errors.New("没有找到该用户")
	}
	return iris.StatusOK, user, nil

}

// 查询商家
func (s *userService) GetBusinessById(id int) (int, *model.User, error) {
	status, user, err := s.GetUserById(id)
	if err != nil {
		return status, nil, err
	}
	if !user.IsBusiness() {
		return iris.StatusNotFound, nil, errors.New("没有找到该用户")
	}
	return status, user, nil
}

// 查询
func (s *userService) GetUserById(id int) (int, *model.User, error) {
	user := new(model.User)
	res, err := manager.DBEngine.Where(
		fmt.Sprintf("%s=?", constant.NameID), id).Get(user)
	if err != nil {
		logrus.Errorf("查找用户失败:%s", err)
		return iris.StatusInternalServerError, nil, errors.New("查找用户失败")
	}
	if res == false {
		return iris.StatusNotFound, nil, errors.New("没有找到该用户")
	}
	return iris.StatusOK, user, nil
}

// 获取所有用户
func (s *userService) GetUserList() ([]model.User, error) {
	list := make([]model.User, 0)
	err := manager.DBEngine.Find(&list)
	if err != nil {
		logrus.Errorf("获取所有用户失败:%s", err)
		return nil, errors.New("获取所有用户失败")
	}
	return list, nil
}

// 检测token的角色是否为管理员
func (s *userService) CheckRoleIsManagerWithToken(ctx iris.Context) (int, bool, error) {
	status, res, err := s.CheckRoleWithToken(ctx, constant.RoleManager)
	if !res {
		return status, false, err
	}
	return iris.StatusOK, true, nil

}

// 检测token的角色
func (s *userService) CheckRoleWithToken(ctx iris.Context, role int) (int, bool, error) {
	status, user, err := s.GetUserFormToken(ctx)
	if err != nil {
		return status, false, err
	}

	if user.Role != role {
		return iris.StatusUnauthorized, false, errors.New("角色不匹配")
	}
	return iris.StatusOK, true, nil

}

// 从请求头获取token的用户
func (s *userService) GetUserFormToken(ctx iris.Context) (int, *model.User, error) {
	status, id, err := authentication.GetUserIDFormHeaderToken(ctx)
	if err != nil {
		return status, nil, err
	}
	var user *model.User
	status, user, err = s.GetUserById(id)
	if err != nil {
		return status, nil, err
	}
	return status, user, nil

}

// 设置token
func setToken(user *model.User, token string) {
	user.Token = token
	_, err := manager.DBEngine.Id(user.Id).Update(user)
	if err != nil {
		logrus.Errorf("更新用户token失败: %s", err)
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
