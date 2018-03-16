package routes

import (
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/chikong/ordersystem/services"
)

var mWebServer *websocket.Server

func LoadWebSocketRoutes(b *bootstrap.Bootstrapper) {
	setupWebSocket(b.Application)

}

func setupWebSocket(app *iris.Application) {
	mWebServer = websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	printerService := services.NewPrinterService()

	mWebServer.OnConnection(func(c websocket.Connection) {
		printerService.HandleConnection(c)

	})

	app.Get("/echo", mWebServer.Handler())


}

func GetWebServer() *websocket.Server {
	return mWebServer
}