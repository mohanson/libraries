package main

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/mohanson/libraries/go/doa"
)

type Math struct{}

func (m *Math) Add(arg []int, ret *int) error {
	for _, e := range arg {
		*ret += e
	}
	return nil
}

func mainServer() {
	doa.Nil(rpc.Register(&Math{}))
	rpc.HandleHTTP()
	l := doa.Try(net.Listen("tcp", "127.0.0.1:8080"))
	go http.Serve(l, nil)
}

func mainClient() {
	client := doa.Try(rpc.DialHTTP("tcp", "127.0.0.1:8080"))
	defer client.Close()
	ret := 0
	doa.Nil(client.Call("Math.Add", []int{1, 2, 3, 4}, &ret))
	doa.Doa(ret == 10)
}

func main() {
	mainServer()
	mainClient()
}
