package main

import (
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/chikong/ordersystem/middleware/identity"
	"github.com/chikong/ordersystem/routes"
	"github.com/chikong/ordersystem/constant"
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/network"
)

var app = bootstrap.New(
	constant.SystemName,
	constant.SystemOwner,
	identity.Configure,
	routes.Configure,
)

func init() {
	app.Bootstrap()
}

func main() {
	go func() {
		socketApp := iris.New()
		network.SetupWebSocket(socketApp)
		socketApp.Run(iris.Addr(":8091"))
	}()

	app.Listen(":8090", iris.WithPostMaxMemory(32<<20))


}

