package main

import (
	"fmt"
	"net"
	"net/rpc"
)

func main() {
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8099")
	conn, _ := net.DialTCP("tcp", nil, addr)
	defer conn.Close()

	rc := rpc.NewClient(conn)
	err := rc.Call("RPC.Markdown", &in, &out)
	fmt.Println(err)
	fmt.Println(string(out))
}
