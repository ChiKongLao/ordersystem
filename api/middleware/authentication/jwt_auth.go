package authentication

import (
	"github.com/kataras/iris/middleware/basicauth"
	"time"
)

var JWTAuth = basicauth.Config{
	Users:   map[string]string{"userName": "password"},
	Expires: time.Duration(24) * time.Hour * 7,
}

