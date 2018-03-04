package manager

import (
	"github.com/garyburd/redigo/redis"
	"github.com/chikong/ordersystem/constant"
	"github.com/sirupsen/logrus"
)

const (
	RedisGet    = "GET"
	RedisSet    = "SET"
	RedisExists = "EXISTS"
	RedisDel    = "DEL"
	RedisRPush  = "RPUSH"  // 列表,表尾
	RedisLPush  = "LPUSH"  // 列表,表头
	RedisLRange = "LRANGE" // 列表
	RedisLLen   = "LLEN"   // 列表

	RedisEx     = "EX"     // key-value同时设置过期
	RedisExpire = "EXPIRE" // 单独设置过期

	RedisSetNX = "SETNX" // SET if Not eXists

)

type RedisManager struct {
}

var mConn redis.Conn

func InitRedis() {
	conn, err := redis.Dial("tcp", constant.RedisHost)
	if err != nil {
		logrus.Error("连接redis失败", err)
		return
	}
	conn.Do("SELECT",2) // 选择数据库
	mConn = conn

}

func GetRedisConn() redis.Conn {
	if mConn == nil {
		InitRedis()
	}
	return mConn
}
