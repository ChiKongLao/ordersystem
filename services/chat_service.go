package services

import (
	"github.com/kataras/iris"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/chikong/ordersystem/network"
	"github.com/chikong/ordersystem/manager"
	"strings"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/kataras/iris/core/errors"
	"github.com/chikong/ordersystem/constant"
)

type ChatService interface {
	GetChatLog(businessId, tableId int) (int, []string, error)
}

func NewChatService(userService UserService) ChatService {
	handleChatMessage()
	return &chatService{
		UserService: userService,
	}
}

type chatService struct {
	UserService UserService
}

// 处理聊天消息
func handleChatMessage() {
	network.GetMqttInstance().RegisterCallback(func(msg mqtt.Message) {
		if !strings.Contains(msg.Topic(), network.TopicChat) {
			return
		}
		key := msg.Topic()
		if ok, _ := redis.Bool(manager.GetRedisConn().Do(manager.RedisRPush, key, string(msg.Payload()))); ok {
			manager.GetRedisConn().Do(manager.RedisExpire, key, constant.TimeCacheChatLog)

		}

	})

}

// 获取聊天记录
func (s *chatService) GetChatLog(businessId, tableId int) (int, []string, error) {
	key := fmt.Sprintf(network.MqttProject+"/%v/%v"+network.TopicChat, businessId, tableId)
	len, _ := redis.Int(manager.GetRedisConn().Do(manager.RedisLLen, key))
	list, err := redis.Strings(manager.GetRedisConn().Do(manager.RedisLRange, key, 0, len))
	if err != nil {
		logrus.Warnf("获取redis聊天记录失败. %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取redis聊天记录失败")
	}
	if list == nil {
		list = make([]string, 0)
	}
	return iris.StatusOK, list, nil
}
