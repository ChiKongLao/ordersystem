package main

import (
	"github.com/henrylee2cn/teleport"
)

func main() {
	tp.SetLoggerLevel("ERROR")

	cli := tp.NewPeer(tp.PeerConfig{})
	defer cli.Close()

	sess, err := cli.Dial(":8091")
	if err != nil {
		tp.Fatalf("%v", err)
	}

	var reply string
	rerr := sess.Pull("/",
		"hello",
		&reply,
	).Rerror()

	if rerr != nil {
		tp.Fatalf("err= %v", rerr)
	}
	tp.Printf("reply2: %s", reply)
}
