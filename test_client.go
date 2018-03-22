package main

import (
	"github.com/henrylee2cn/teleport"
	"net"
	"github.com/henrylee2cn/teleport/socket"
	"log"
)

func main() {

	test()
}

func test()  {
	conn, err := net.Dial("tcp", "127.0.0.1:8091")
	if err != nil {
		log.Fatalf("[CLI] dial err: %v", err)
	}
	s := socket.GetSocket(conn)
	defer s.Close()
	s.Write([]byte("hello"))

}

func test2()  {
	tp.SetLoggerLevel("ERROR")

	cli := tp.NewPeer(tp.PeerConfig{})
	defer cli.Close()

	sess, err := cli.Dial(":8091")
	if err != nil {
		tp.Fatalf("%v", err)
	}

	rerr := sess.Push("/",
		"hello",
	)

	if rerr != nil {
		tp.Fatalf("err= %v", rerr)
	}

}
