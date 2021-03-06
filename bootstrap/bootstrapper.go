package bootstrap

import (
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/middleware/logger"
	"github.com/chikong/ordersystem/constant"
	"github.com/go-xorm/xorm"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/chikong/ordersystem/manager"
	_"github.com/go-sql-driver/mysql"
	"github.com/chikong/ordersystem/network"
	"github.com/chikong/ordersystem/configs"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
	DBEngine *xorm.Engine
}

// New returns a new Bootstrapper.
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	configs.InitApplicationConfig()	// 读取配置文件
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

// SetupViews loads the templates.
func (b *Bootstrapper) SetupViews(viewsDir string) {
	//b.RegisterView(iris.HTML(viewsDir, ".html").Layout("shared/layout.html"))
	b.RegisterView(iris.HTML(viewsDir, ".html"))
}

// SetupErrorHandlers prepares the http error handlers (>=400).
func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			//"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"msg": ctx.Values().GetString("message"),
		}
		ctx.JSON(err)
	})
}

// SetupDatabase engine default use MySQL for xorm.io
func (b *Bootstrapper) SetupDatabaseEngine() {
	// 创建 ORM 引擎与数据库
	initDBEngine(b)
}
func initDBEngine(b *Bootstrapper){
	// 创建 ORM 引擎与数据库
	engine, err := xorm.NewEngine("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			configs.GetConfig().MySQL.UserName,
			configs.GetConfig().MySQL.Password,
			configs.GetConfig().MySQL.Host,
			configs.GetConfig().MySQL.Name))
	if err != nil {
		logrus.Errorf("连接数据库失败: %s",err)
		return
	}

	b.DBEngine = engine
	manager.SetDBEngine(engine)
}

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.Configure(iris.WithConfiguration(iris.YAML(constant.ROOT+"/configs/test.yml")))
	b.SetupViews("./web/views")
	b.SetupErrorHandlers()
	b.SetupDatabaseEngine()
	go network.GetMqttInstance()
	//network.SetupWebSocket(b.Application)
	go network.GetSocketInstance() // Socket

	// middleware, after static files
	b.Use(recover.New())
	b.Use(logger.New())
	logrus.SetLevel(logrus.DebugLevel)

	return b
}

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	if err := b.Run(iris.Addr(addr), cfgs...); err != nil {
		b.Logger().Warn("Shutdown with error: " + err.Error())
	}
}