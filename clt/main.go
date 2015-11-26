package main

import (
	"fmt"
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
)

func main() {
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:80")
	// addr, err := net.ResolveTCPAddr("tcp", "rpcsvr.daoapp.io:80")
	if goutils.CheckErr(err) {
		return
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	defer conn.Close()
	if goutils.CheckErr(err) {
		return
	}

	rc := rpc.NewClient(conn)
	err = rc.Call("RPC.Markdown", &in, &out)
	if goutils.CheckErr(err) {
		return
	}
	fmt.Println(string(out))
}
