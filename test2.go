package main

import (
	"time"

	tp "github.com/henrylee2cn/teleport"
	"log"
	"github.com/chikong/ordersystem/network/proto"
)

type Home struct {
	tp.PushCtx
}


func main() {
	// Server
	svr := tp.NewPeer(tp.PeerConfig{ListenAddress: ":8091"})
	svr.SetUnknownPush(func(ctx tp.UnknownPushCtx) *tp.Rerror {
			println(" 接收到: ",string(ctx.InputBodyBytes()))
		return nil
	})

	svr.RoutePushFunc(func(ctx tp.PushCtx, args *string) *tp.Rerror {
		tp.Printf("RoutePushFunc 接收到: %s",*args)
		return nil
	})

	go svr.Listen(proto.NewJsonProtoFunc2)
	time.Sleep(1e9)

	// Client
	cli := tp.NewPeer(tp.PeerConfig{})
	sess, err := cli.Dial(":8091", proto.NewJsonProtoFunc2)
	if err != nil {
		log.Fatalln("err=",err)
	}

	sess.Push("/",[]byte("hello"))


	time.Sleep(3e9)
	//time.Sleep(3111e9)
}
