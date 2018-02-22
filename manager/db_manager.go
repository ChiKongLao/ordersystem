package manager

import (
	"github.com/go-xorm/xorm"
	"github.com/chikong/ordersystem/model"
	"github.com/sirupsen/logrus"
)

// 根据表生成model
// xorm reverse mysql root:@/order_system?charset=utf8 /Users/chikong/go_workspace/src/github.com/go-xorm/cmd/xorm/templates/goxorm modeltmp
// xorm reverse mysql root:@/order_system?charset=utf8 H:\GOPATH\src\github.com\go-xorm\cmd\xorm\templates\goxorm modeltmp



type DbManager interface {
	SetDBEngine(engine *xorm.Engine)
}

var DBEngine *xorm.Engine //定义引擎全局变量

// 设置数据库引擎
func SetDBEngine(engine *xorm.Engine) {
	DBEngine = engine
	sync(engine)
}

// 同步结构体到数据表，创建对应的表
func sync(engine *xorm.Engine){
	unSuccessTableName := ""
	var err error
	defer func() {
		if unSuccessTableName != "" {
			logrus.Errorf("创建%s表失败: %v\n",unSuccessTableName, err)

		}
	}()

	if err = engine.Sync2(new(model.User)); err != nil {
		unSuccessTableName = "user"
		return
	}
	if err = engine.Sync2(new(model.Food)); err != nil {
		unSuccessTableName = "food"
		return
	}
	if err = engine.Sync2(new(model.TableInfo)); err != nil {
		unSuccessTableName = "table_info"
		return
	}
	if err = engine.Sync2(new(model.Order)); err != nil {
		unSuccessTableName = "order"
		return
	}
	if err = engine.Sync2(new(model.ShoppingCart)); err != nil {
		unSuccessTableName = "shopping_cart"
		return
	}
	if err = engine.Sync2(new(model.CollectFood)); err != nil {
		unSuccessTableName = "collect_food"
		return
	}

}



