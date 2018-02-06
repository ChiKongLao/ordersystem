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
	"github.com/chikong/ordersystem/datamodels"
	"strings"
)

const (
	SecretKey = "Hello World!!!"
)

//var JWTAuth = basicauth.Config{
//	Users:   map[string]string{"userName": "password"},
//	Realm:   basicauth.DefaultBasicAuthRealm,
//	Expires: time.Duration(24) * time.Hour * 7,
//}

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
func MakeToken(userName, password string) (string,error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24) * 7).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userName"] = userName
	claims["password"] = password
	token.Claims = claims

	signedString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logrus.Errorf("生成token失败: %s",err)
		return "",errors.New("生成token失败")
	}

	return "Bearer "+signedString,nil
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
	tokenString := strings.Replace(ctx.GetHeader(datamodels.NameAuthorization),"Bearer ","",1)
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
	token, err := GetTokenFormHeader(ctx)
	if err != nil {
		logrus.Errorf("从请求头获取token信息失败: %s",err)
		return "",errors.New("解析token信息失败")
	}
	claim := token.Claims.(jwt.MapClaims)
	return claim[datamodels.NameUserName].(string),nil
}


