package clt

import (
	"fmt"
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
	"testing"
)

func TestC(t *testing.T) {
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:88")
	// addr, err := net.ResolveTCPAddr("tcp", "rpcsvr.daoapp.io:88")
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

func TestRPC(t *testing.T) {
	c := RPCClient("127.0.0.1:88")
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)
	Markdown(c, &in, &out)
	fmt.Println(goutils.ToString(out))
}
