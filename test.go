package main

import (
	tp "github.com/henrylee2cn/teleport"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

func main() {
	svr := tp.NewPeer(tp.PeerConfig{
		CountTime:     true,
		ListenAddress: ":8091",
	})
	svr.SetUnknownPull(func(ctx tp.UnknownPullCtx) (interface{}, *tp.Rerror) {
		var v = struct {
			RawMessage json.RawMessage
			Bytes      []byte
		}{}
		ctx.Bind(&v)
		logrus.Infoln("args=",v)
		return "Unknown",nil
	})

	svr.Listen()
}