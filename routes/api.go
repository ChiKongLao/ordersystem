package routes

import (
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/logger"
	"github.com/chikong/ordersystem/controllers"
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
		userService := services.NewUserService()

		userParty := v1.Party("/user")
		mvc.Configure(userParty, func(mvcApp *mvc.Application) {
			mvcApp.Register(userService)
			mvcApp.Handle(new(controllers.UserController))

		})
		mvc.Configure(v1.Party("/menu",auth), func(mvcApp *mvc.Application) {
			service := services.NewDashesService()
			mvcApp.Register(service,userService)
			mvcApp.Handle(new(controllers.DashesController))

		})
		mvc.Configure(v1.Party("/table",auth), func(mvcApp *mvc.Application) {
			service := services.NewTableService()
			mvcApp.Register(service,userService)
			mvcApp.Handle(new(controllers.TableController))

		})

		mvc.New(v1.Party("/message",auth)).Handle(new(controllers.MessageController))



	}
}