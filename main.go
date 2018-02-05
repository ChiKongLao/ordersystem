package main

import (
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/chikong/ordersystem/middleware/identity"
	"github.com/chikong/ordersystem/routes"
	"github.com/chikong/ordersystem/constant"
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
	app.Listen(":8090")

}