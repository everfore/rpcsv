package rpcsv

import (
	"encoding/json"
	"fmt"
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
	"testing"
)

var (
	lis            net.Listener
	rpc_tcp_server = "tcphub.t0.daoapp.io:61142"
)

func init() {
	// lis, _ = RPCServe("8800")
	// lis, _ = RPCServeWithCode("8800")
}

func TestC(t *testing.T) {
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)

	addr, err := net.ResolveTCPAddr("tcp", rpc_tcp_server)
	// addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8800")
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
	c := RPCClientWithCodec(rpc_tcp_server)
	// c := RPCClientWithCodec("127.0.0.1:8800")
	defer c.Close()
	in := []byte("#   [Hi](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)
	Markdown(c, &in, &out)
	fmt.Println(goutils.ToString(out))
	// defer lis.Close()
}

func TestJob(t *testing.T) {
	c := RPCClientWithCodec(rpc_tcp_server)
	// c := RPCClientWithCodec("127.0.0.1:8800")
	defer c.Close()
	in, _ := json.Marshal(Job{Name: "baidu", Target: "https://baidu.com"})
	out := make([]byte, 10)
	err := c.Call("RPC.Job", &in, &out)
	goutils.CheckErr(err)
	fmt.Println(goutils.ToString(in))
	fmt.Println(goutils.ToString(out))
}

func TestWall(t *testing.T) {
	return
	c := RPCClientWithCodec(rpc_tcp_server)
	// c := RPCClientWithCodec("127.0.0.1:8800")
	defer c.Close()
	out := make([]byte, 10)
	in := make([]byte, 1)
	err := c.Call("RPC.Wall", &in, &out)
	goutils.CheckErr(err)
	fmt.Println("out:", goutils.ToString(out))
}
