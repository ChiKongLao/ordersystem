package routes

import (
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/logger"
	"github.com/chikong/ordersystem/web/controllers"
)

func LoadAPIRoutes(b *bootstrap.Bootstrapper) {
	//setup cross domain credentials
	//crs := cors.New(cors.Options{
	//	AllowedOrigins:		[]string{"*"},
	//	AllowCredentials:	true,
	//})

	//b.Logger().SetOutput(logger.NewLogFile())

	//setup version 1 routes
	v1 := b.Party("/v1")

	v1.Use(recover.New())
	//v1.Use(languages.CurrentLanguage)
	v1.Use(logger.GetRequestLogger())
	//v1.Use(basicauth.New(authentication.JWTAuth))
	//v1.Use(crs)

	{
		mvc.Configure(v1.Party("/user"), func(mvcApp *mvc.Application) {
			service := services.NewUserService()
			mvcApp.Register(service)
			mvcApp.Handle(new(controllers.UserController))

		})


	}
}