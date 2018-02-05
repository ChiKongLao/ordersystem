package authentication

import (
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
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
		//s = fmt.Sprintf("{\"msg\":\"认证失败: %s\"",s)
		res := map[string]string{
			"msg" : s,
		}

		jwtmiddleware.OnError(ctx,s)
	},

})
