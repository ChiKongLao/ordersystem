package main

import (
	"github.com/henrylee2cn/teleport"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/chikong/ordersystem/network/proto"
	"github.com/henrylee2cn/teleport/socket"
	"net"
	"log"
	"time"
)

func main() {
	go read2()
	time.Sleep(1e9)
	client()
	time.Sleep(3e9)
}

func read() {
	svr := tp.NewPeer(tp.PeerConfig{
		CountTime:     true,
		ListenAddress: ":8091",
	})

	//svr.SetUnknownPush(func(ctx tp.UnknownPushCtx, args string) (string, *tp.Rerror) {
	//	data := string(ctx.InputBodyBytes())
	//	logrus.Infoln("push=", data)
	//	return "hello",nil
	//})
	//svr.RoutePushFunc(func(data string) {
	//	logrus.Infoln("RoutePushFunc=", data)
	//
	//})
	svr.RoutePushFunc(func(ctx tp.PushCtx, args *string) (*tp.Rerror) {
		//data := string(ctx.InputBodyBytes())
		logrus.Infoln("push=", args)
		return nil
	})

	//svr.RoutePullFunc(func(data string) {
	//
	//	logrus.Infoln("RoutePullFunc=", data)
	//
	//})

	svr.Listen(proto.NewStringProtoFunc)


}

func read2()  {
	socket.SetNoDelay(false)
	socket.SetPacketSizeLimit(512)
	lis, err := net.Listen("tcp", "0.0.0.0:8091")
	if err != nil {
		log.Fatalf("[SVR] listen err: %v", err)
	}
	log.Printf("listen tcp 0.0.0.0:8091")
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("[SVR] accept err: %v", err)
		}
		go func(s socket.Socket) {
			log.Printf("accept %s", s.Id())
			defer s.Close()
			for {
				var data []byte
				_, err = s.Read(data)

				if err != nil {
					log.Fatalf("err: %s",err)
					continue
				}

				if len(data) != 0 {
					println(string(data))
				}
			}
		}(socket.GetSocket(conn,proto.NewJsonProtoFunc2))
	}
}

func push() {

	svr := tp.NewPeer(tp.PeerConfig{
		CountTime:     true,
		ListenAddress: ":8091",
	})

	svr.Listen()
	svr.SetUnknownPush(func(ctx tp.UnknownPushCtx) *tp.Rerror {
		data := string(ctx.InputBodyBytes())
		logrus.Infoln("SetUnknownPush=", data)
		return nil
	})
}

func pull() {

	svr := tp.NewPeer(tp.PeerConfig{
		CountTime:     true,
		ListenAddress: ":8091",
	})

	svr.Listen()
	svr.SetUnknownPull(func(ctx tp.UnknownPullCtx) (interface{}, *tp.Rerror) {
		var v = struct {
			RawMessage json.RawMessage
			Bytes      []byte
		}{}
		ctx.Bind(&v)
		logrus.Infoln("args=", v)
		return "Unknown", nil
	})
}

func client()  {
	conn, err := net.Dial("tcp", "0.0.0.0:8091")
	if err != nil {
		log.Fatalf("[CLI] dial err: %v", err)
	}
	s := socket.GetSocket(conn)
	defer s.Close()
	s.Write([]byte("hello"))
}