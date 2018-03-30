package network

import (
	"github.com/sirupsen/logrus"
	"github.com/henrylee2cn/teleport/socket"
	"net"
	"sync"
	"time"
	"strings"
)


type SocketManager interface {
	OnConnect(func(s socket.Socket))
	OnMessage(func(s socket.Socket, payload string))
}

type socketManager struct {
	// 回调列
	socketCallbackList []func(s socket.Socket, payload string)
	connectCallbackList []func(s socket.Socket)
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
		soc := socket.GetSocket(conn)
		for _, value := range mSocketInstance.connectCallbackList {
			value(soc)
		}
		go func(s socket.Socket) {
			logrus.Infoln("打印机连接服务器: ", s.Id())
			defer s.Close()
			for {
				data := make([]byte, 2048)
				dataLen, err := s.Read(data)

				if err != nil {
					errString := err.Error()
					if strings.Contains(errString,"EOF") {
						logrus.Infof("打印机断开连接", s.Id())
						break
					}
					if strings.Contains(errString,"closed by the remote host") {
						logrus.Infof("远程主机强制关闭现有的连接", s.Id())
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
		}(soc)
	}

}

func (m *socketManager) OnMessage(fun func(s socket.Socket, payload string)) {
	m.socketCallbackList = append(m.socketCallbackList, fun)
}

func (m *socketManager) OnConnect(fun func(s socket.Socket)) {
	m.connectCallbackList = append(m.connectCallbackList, fun)
}
