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
	"github.com/chikong/ordersystem/model"
	"encoding/json"
)

type ChatService interface {
	GetChatLog(businessId, tableId int) (int, []model.ChatMsg, error)
}

func NewChatService(userService UserService, tableService TableService) ChatService {
	handleChatMessage()
	return &chatService{
		UserService:  userService,
		TableService: tableService,
	}
}

type chatService struct {
	UserService  UserService
	TableService TableService
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
func (s *chatService) GetChatLog(businessId, tableId int) (int, []model.ChatMsg, error) {
	key := fmt.Sprintf(network.MqttProject+"/%v/%v"+network.TopicChat, businessId, tableId)
	count, _ := redis.Int(manager.GetRedisConn().Do(manager.RedisLLen, key))
	stringList, err := redis.Strings(manager.GetRedisConn().Do(manager.RedisLRange, key, 0, count))
	if err != nil {
		logrus.Warnf("获取redis聊天记录失败. %s", err)
		return iris.StatusInternalServerError, nil, errors.New("获取聊天记录失败")
	}
	if stringList == nil {
		stringList = make([]string, 0)
	}
	list := make([]model.ChatMsg,0)
	for _, subItem := range stringList {
		var item model.ChatMsg
		err = json.Unmarshal([]byte(subItem), &item)
		if err != nil {
			logrus.Warnf("解析聊天记录失败. %s", err)
			return iris.StatusInternalServerError, list, errors.New("解析聊天记录失败")
		}
		list = append(list, item)
	}

	return iris.StatusOK, list, nil
}
