package routes

import (
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/logger"
	"github.com/chikong/ordersystem/web/controllers"
	"github.com/chikong/ordersystem/api/middleware/authentication"
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
	//v1.Use(authentication.JWTHandler.Serve)
	//v1.Use(crs)

	auth := authentication.JWTHandler.Serve

	{
		userParty := v1.Party("/user")
		mvc.Configure(userParty, func(mvcApp *mvc.Application) {
			service := services.NewUserService()
			mvcApp.Register(service)
			mvcApp.Handle(new(controllers.UserController))

		})

		mvc.New(v1.Party("/message",auth)).Handle(new(controllers.MessageController))


	}
}