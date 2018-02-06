package manager

import (
	"github.com/go-xorm/xorm"
	"github.com/chikong/ordersystem/datamodels"
	"github.com/sirupsen/logrus"
)

// 根据表生成model
// xorm reverse mysql root:@/order_system?charset=utf8 /Users/chikong/go_workspace/src/github.com/go-xorm/cmd/xorm/templates/goxorm



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
		if unSuccessTableName == "user" {
			logrus.Errorf("创建用户表失败: %v\n", err)
		}else{
			logrus.Info("初始化数据表成功")

		}
	}()

	if err = engine.Sync2(new(datamodels.User)); err != nil {
		unSuccessTableName = "user"
		return
	}

}



