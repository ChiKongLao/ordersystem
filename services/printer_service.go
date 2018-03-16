package services

import (
	"github.com/kataras/iris/websocket"
	"github.com/sirupsen/logrus"
	"regexp"
	"fmt"
	"github.com/chikong/ordersystem/network"
	"strings"
)

type PrinterService interface {
	HandleConnection(c websocket.Connection)
}

func NewPrinterService() PrinterService {

	return &printerService{

	}
}

type printerService struct {

	}



func (* printerService) HandleConnection(c websocket.Connection) {

	c.OnMessage(func(bytes []byte) {
		payload := string(bytes)
		result := handlePayload(payload)
		if result == "" {
			logrus.Warnln("未找到对应的事件:",payload)
			return
		}
		c.EmitMessage([]byte(result))
	})

}

func handlePayload(payload string) string {
	ok := false
	if ok, _ = regexp.MatchString(fmt.Sprintf(network.SocketKeyPing,"\\d+","\\d"),payload); ok {
		return handlePing(payload)
	}

	return ""
}

// 处理心跳
func handlePing(payload string) string {
	size := len(payload)
	data := []rune(payload)
	deviceId := data[2:strings.Index(payload,"*AS")-2]
	status := data[size-7:size-6]
	fmt.Println(string(deviceId),string(status))

	return network.SocketKeyPingReply
}