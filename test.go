package main

import (
	"github.com/henrylee2cn/teleport"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"io"
	"github.com/henrylee2cn/teleport/socket"
)

func main() {
	svr := tp.NewPeer(tp.PeerConfig{
		CountTime:     true,
		ListenAddress: ":8091",
	})

	svr.SetUnknownPush(func(ctx tp.UnknownPushCtx) *tp.Rerror {
		data := string(ctx.InputBodyBytes())
		logrus.Infoln("push=",data)
		return nil
	})

	svr.Listen(func(writer io.ReadWriter) socket.Proto {

	})

	svr.Listen()
}

func push(svr tp.Peer)  {
	svr.SetUnknownPush(func(ctx tp.UnknownPushCtx) *tp.Rerror {
		data := string(ctx.InputBodyBytes())
		logrus.Infoln("push=",data)
		return nil
	})
}

func pull(svr tp.Peer)  {
	svr.SetUnknownPull(func(ctx tp.UnknownPullCtx) (interface{}, *tp.Rerror) {
		var v = struct {
			RawMessage json.RawMessage
			Bytes      []byte
		}{}
		ctx.Bind(&v)
		logrus.Infoln("args=",v)
		return "Unknown",nil
	})
}