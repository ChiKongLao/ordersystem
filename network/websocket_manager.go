package network

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/sirupsen/logrus"
)

var mWebServer *websocket.Server

func SetupWebSocket(app *iris.Application) {
	mWebServer = websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	mWebServer.OnConnection(func(c websocket.Connection) {
		logrus.Debugf("打印机连接服务器: %s",c.ID())
		c.OnDisconnect(func() {
			logrus.Debugf("打印机断开连接: %s",c.ID())
		})
		println("origin= ",c.Context().GetHeader("Origin"))
		logrus.Infoln(c)


	})

	app.Get("/", mWebServer.Handler())


}

func GetWebServer() *websocket.Server {
	return mWebServer
}