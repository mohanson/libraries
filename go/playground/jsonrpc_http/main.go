package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/mohanson/libraries/go/doa"
)

type Math struct{}

func (m *Math) Add(arg []int, ret *int) error {
	for _, e := range arg {
		*ret += e
	}
	return nil
}

type Conn struct {
	io.Reader
	io.Writer
	io.Closer
}

func mainServer() {
	doa.Nil(rpc.Register(&Math{}))
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		rwc := Conn{
			Reader: r.Body,
			Writer: w,
			Closer: r.Body,
		}
		codec := jsonrpc.NewServerCodec(rwc)
		rpc.ServeRequest(codec)
	})
	l := doa.Try(net.Listen("tcp", "127.0.0.1:8080"))
	go http.Serve(l, nil)
}

func mainClient() {
	bodyJson := map[string]any{
		"jsonrpc": "2.0",
		"method":  "Math.Add",
		"params":  []any{[]int{1, 2, 3, 4}},
		"id":      1,
	}
	bodyData := doa.Try(json.Marshal(bodyJson))
	resp := doa.Try(http.Post("http://127.0.0.1:8080/rpc", "application/json", bytes.NewBuffer(bodyData)))
	body := map[string]any{}
	doa.Nil(json.NewDecoder(resp.Body).Decode(&body))
	doa.Doa(body["result"].(float64) == 10)
}

func main() {
	mainServer()
	mainClient()
}
