package authentication

import (
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/kataras/iris"
	"time"
	"errors"
	"strings"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/model"
	"fmt"
	"strconv"
)

const (
	SecretKey = "Hello World!!!"
)


var JWTHandler = jwtmiddleware.New(jwtmiddleware.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	},
	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(ctx context.Context, s string) {
		res := map[string]string{
			"msg" : "认证失败",
		}
		logrus.Errorf("jwt认证失败: %s",s)

		data,_ := json.Marshal(res)
		jwtmiddleware.OnError(ctx,string(data))
	},

})

// 生成token
func MakeToken(user *model.User) (string,error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24) * 7).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userName"] = user.UserName
	claims["password"] = user.Password
	claims["id"] = user.Id
	token.Claims = claims

	signedString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logrus.Errorf("生成token失败: %s",err)
		return "",errors.New("生成token失败")
	}

	return fmt.Sprintf("Bearer %s",signedString),nil
}

// 从Context获取token
func GetToken(ctx iris.Context) (token *jwt.Token){
	token = ctx.Values().Get(JWTHandler.Config.ContextKey).(*jwt.Token)
	return token
}

// 从Context获取token string
func GetTokenString(ctx iris.Context) string{
	return GetToken(ctx).Raw
}

// 从请求头获取token
func GetTokenFormHeader(ctx iris.Context) (*jwt.Token,error){
	tokenString := strings.Replace(ctx.GetHeader(constant.NameAuthorization),"Bearer ","",1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		logrus.Errorf("从请求头获取token失败: %s",err)
		return nil,errors.New("解析token失败")
	}
	return token,nil
}

// 从请求头获取token的用户名
func GetUserNameFormHeaderToken(ctx iris.Context) (string,error){
	claim, err := getClaims(ctx)
	if err != nil {
		return "",err
	}

	return claim[constant.NameUserName].(string),nil
}

// 从请求头获取token对应的用户ID
func GetUserIDFormHeaderToken(ctx iris.Context) (string,error){
	claim, err := getClaims(ctx)
	if err != nil {
		return "",err
	}
	res := strconv.Itoa(int(claim[constant.NameID].(float64)))
	return res,nil
}

// 从请求头获取token对应的信息
func getClaims(ctx iris.Context) (jwt.MapClaims,error){
	token, err := GetTokenFormHeader(ctx)
	if err != nil {
		logrus.Errorf("从请求头获取token信息失败: %s",err)
		return nil,errors.New("解析token信息失败")
	}

	return token.Claims.(jwt.MapClaims),nil
}

// 从userId和token的id对比是不为自己
func IsOwnWithToken(ctx iris.Context, userId string) (bool, error){
	id, err := GetUserIDFormHeaderToken(ctx)
	if err != nil {
		return false, err
	}

	if userId != id {
		return false, errors.New("帐号不匹配")
	}
	return userId == id,nil

}

