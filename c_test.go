package rpcsv

import (
	"fmt"
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
	"testing"
)

var (
	lis net.Listener
)

func init() {
	// lis, _ = RPCServe("8800")
}
func TestC(t *testing.T) {
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8800")
	// addr, err := net.ResolveTCPAddr("tcp", "182.254.132.59:8800")
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
	c := RPCClient("127.0.0.1:8800")
	// c := RPCClient("182.254.132.59:8800")
	defer c.Close()
	in := []byte("#   [Hi](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)
	Markdown(c, &in, &out)
	fmt.Println(goutils.ToString(out))
	// defer lis.Close()
}
