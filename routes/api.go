package routes

import (
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/api/middleware/logger"
	"github.com/chikong/ordersystem/controllers"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/chikong/ordersystem/constant"
	"strconv"
	"github.com/chikong/ordersystem/model"
	"fmt"
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

	// addTestData
	//	b.SetupDatabaseEngine()
	//	addTestData()


	userService := services.NewUserService()
	dashesService := services.NewDashesService()


	auth := authentication.JWTHandler.Serve
	{

		userParty := v1.Party("/user")
		mvc.Configure(userParty, func(mvcApp *mvc.Application) {
			mvcApp.Register(userService)
			mvcApp.Handle(new(controllers.UserController))
		})
		mvc.Configure(v1.Party("/menu",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(dashesService,userService)
			mvcApp.Handle(new(controllers.DashesController))

		})
		mvc.Configure(v1.Party("/table",auth), func(mvcApp *mvc.Application) {
			service := services.NewTableService()
			mvcApp.Register(service,userService)
			mvcApp.Handle(new(controllers.TableController))

		})
		mvc.Configure(v1.Party("/order",auth), func(mvcApp *mvc.Application) {
			service := services.NewOrderService(dashesService)
			mvcApp.Register(service,userService)
			mvcApp.Handle(new(controllers.OrderController))

		})

		mvc.New(v1.Party("/message",auth)).Handle(new(controllers.MessageController))

	}
}

func addTestData()  {

	userService := services.NewUserService()
	dashesService := services.NewDashesService()
	tableService := services.NewTableService()

	addUser(userService)
	addDashes(userService,dashesService)
	addTable(userService,tableService)





}

func addUser(userService services.UserService)  {
	userService.InsertUser(constant.RoleManager,"admin","admin","admin")
	for i := 0; i < 10 ; i = i + 1 {
		userService.InsertUser(constant.RoleBusiness,
			"business"+strconv.Itoa(i),"111",
			fmt.Sprintf("商家%v",i))
	}
	for i := 0; i < 100 ; i = i + 1 {
		userService.InsertUser(constant.RoleCustomer,
			"customer"+strconv.Itoa(i),"111",
			fmt.Sprintf("客户%v",i))

	}
}

func addDashes(userService services.UserService, dashesService services.DashesService){
	userList, err := userService.GetUserList()
	if err != nil {
		return
	}
	for _, user := range userList {
		if user.Role != constant.RoleBusiness{
			continue
		}
		list := make([]*model.Dashes,100)

		for j := 0; j < len(list) ; j = j+1  {
			list[j] = &model.Dashes{
				BusinessId:user.Id,
				Name:fmt.Sprintf("菜式%v",j),
				Num:100,
				Pic:"https://www.baidu.com/img/bd_logo1.png",
				Price:strconv.Itoa(j),
				Desc:fmt.Sprintf("%s的菜式%v",user.NickName,j),

			}
		}
		dashesService.InsertDashes(list)

	}

}

func addTable(userService services.UserService, tableService services.TableService){
	userList, err := userService.GetUserList()
	if err != nil {
		return
	}
	for _, user := range userList {
		if user.Role != constant.RoleBusiness{
			continue
		}

		for j := 0; j < 100 ; j = j+1  {
			tableService.InsertTable(&model.TableInfo{
				BusinessId:user.Id,
				Name:fmt.Sprintf("%s的餐桌%v",user.NickName,j),
				Capacity:j % 8,

			})
		}
	}

}