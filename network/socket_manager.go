package network

import (
	"github.com/sirupsen/logrus"
	"github.com/henrylee2cn/teleport/socket"
	"net"
	"sync"
	"time"
)


type SocketManager interface {
	RegisterSocketCallback(func(s socket.Socket, payload string))
}

type socketManager struct {
	// 回调列
	socketCallbackList []func(s socket.Socket, payload string)
}

var mSocketInstance *socketManager
var mSocketOnce sync.Once


func GetSocketInstance() SocketManager {
	mSocketOnce.Do(func() {
		mSocketInstance = &socketManager{}
		go mSocketInstance.initSocket()
	})
	return mSocketInstance

}

func(m *socketManager) initSocket() {
	socket.SetNoDelay(false)
	socket.SetPacketSizeLimit(512)
	lis, err := net.Listen("tcp", "0.0.0.0:8091")
	if err != nil {
		logrus.Errorln("启动Socket失败: %v", err)
	}

	logrus.Infoln("启动Socket成功, listen tcp 0.0.0.0:8091")
	for {
		conn, err := lis.Accept()
		if err != nil {
			logrus.Errorln("接受socket失败, err: ", err)

		}
		go func(s socket.Socket) {
			logrus.Infoln("打印机连接服务器: ", s.Id())
			defer s.Close()
			for {
				data := make([]byte, 2048)
				dataLen, err := s.Read(data)

				if err != nil {
					if err.Error() == "EOF" {
						logrus.Infof("打印机断开连接", s.Id())
						break
					}
					logrus.Errorln("解析内容出错, err: ", err)
					continue
				}

				if len(data) != 0 {
					data = data[:dataLen]
					for _, value := range mSocketInstance.socketCallbackList {
						value(s,string(data))
					}
				}
				time.Sleep(10*time.Millisecond)
			}
		}(socket.GetSocket(conn))
	}

}

func (m *socketManager) RegisterSocketCallback(fun func(s socket.Socket, payload string)) {
	m.socketCallbackList = append(m.socketCallbackList, fun)
}
