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
	// setup cross domain credentials
	//crs := cors.New(cors.Options{
	//	AllowedOrigins:		[]string{"*"},
	//	AllowCredentials:	true,
	//})

	//b.Logger().SetOutput(logger.NewLogFile())
	//logger.ConfigLogger()

	//setup version 1 routes
	v1 := b.Party("/v1")
	//v1.Use(crs)

	v1.Use(recover.New())
	//v1.Use(languages.CurrentLanguage)
	v1.Use(logger.GetRequestLogger())
	//v1.Use(authentication.JWTHandler.Serve)


	// addTestData
	//	b.SetupDatabaseEngine()
	//	addTestData()




	auth := authentication.JWTHandler.Serve
	{

		userService := services.NewUserService()

		printerService := services.NewPrinterService()

		classifyService := services.NewClassifyService(userService)
		shopService := services.NewShopService(userService)
		menuService := services.NewMenuService(userService,classifyService)
		tableService := services.NewTableService(userService)
		shoppingCartService := services.NewShoppingService(userService, menuService)
		orderService := services.NewOrderService(userService, menuService,tableService,shoppingCartService,printerService)
		chatService := services.NewChatService(userService,tableService)


		mvc.Configure(v1.Party("/user"), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService)
			mvcApp.Handle(new(controllers.UserController))
		})

		// ###########   系统相关开始
		mvc.Configure(v1.Party("/upload",auth), func(mvcApp *mvc.Application) {
			uploadService := services.NewUploadService()
			mvcApp.Register(userService,uploadService)
			mvcApp.Handle(new(controllers.UploadController))
		})


		// ###########   系统相关结束

		mvc.Configure(v1.Party("/menu",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, menuService)
			mvcApp.Handle(new(controllers.MenuController))

		})
		mvc.Configure(v1.Party("/table",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, tableService)
			mvcApp.Handle(new(controllers.TableController))

		})
		mvc.Configure(v1.Party("/order",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, orderService)
			mvcApp.Handle(new(controllers.OrderController))

		})
		mvc.Configure(v1.Party("/home",auth), func(mvcApp *mvc.Application) {
			service := services.NewHomeService(userService, menuService,tableService,orderService,shopService)
			mvcApp.Register(userService,service)
			mvcApp.Handle(new(controllers.HomeController))
		})
		mvc.Configure(v1.Party("/shopping",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, shoppingCartService)
			mvcApp.Handle(new(controllers.ShoppingController))
		})
		mvc.Configure(v1.Party("/shop",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, shopService)
			mvcApp.Handle(new(controllers.ShopController))

		})
		mvc.Configure(v1.Party("/classify",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, classifyService)
			mvcApp.Handle(new(controllers.ClassifyController))
		})
		mvc.Configure(v1.Party("/chat",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, chatService)
			mvcApp.Handle(new(controllers.ChatController))
		})
		mvc.Configure(v1.Party("/printer",auth), func(mvcApp *mvc.Application) {
			mvcApp.Register(userService, printerService)
			mvcApp.Handle(new(controllers.PrinterController))
		})

		mvc.New(v1.Party("/message",auth)).Handle(new(controllers.MessageController))

	}
}

func addTestData()  {

	userService := services.NewUserService()
	classifyService := services.NewClassifyService(userService)
	foodService := services.NewMenuService(userService,classifyService)
	tableService := services.NewTableService(userService)

	addUser(userService)
	addFood(userService,foodService)
	addTable(userService,tableService)





}

func addUser(userService services.UserService)  {
	userService.InsertUser(constant.RoleManager,"admin","admin","admin","")
	for i := 0; i < 10 ; i = i + 1 {
		userService.InsertUser(constant.RoleBusiness,
			"business"+strconv.Itoa(i),"111",
			fmt.Sprintf("商家%v",i),"")
	}
	for i := 0; i < 100 ; i = i + 1 {
		userService.InsertUser(constant.RoleCustomer,
			"customer"+strconv.Itoa(i),"111",
			fmt.Sprintf("客户%v",i),"")

	}
}

func addFood(userService services.UserService, foodService services.MenuService){
	userList, err := userService.GetUserList()
	if err != nil {
		return
	}
	for _, user := range userList {
		if user.Role != constant.RoleBusiness{
			continue
		}
		list := make([]*model.Food,100)

		for j := 0; j < len(list) ; j = j+1  {
			list[j] = &model.Food{
				BusinessId:user.Id,
				Name:fmt.Sprintf("食物%v",j),
				Num:100,
				Pic:"https://www.baidu.com/img/bd_logo1.png",
				Price:float32(j),
				Desc:fmt.Sprintf("%s的食物%v",user.NickName,j),

			}
		}
		foodService.InsertFood(list)

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