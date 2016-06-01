package rpcsv

import (
	// "encoding/json"
	"fmt"
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
	"testing"
)

var (
	lis            net.Listener
	rpc_tcp_server = "tcphub.t0.daoapp.io:61142"
	// rpc_tcp_server = "127.0.0.1:8800"
)

func init() {
	// lis, _ = RPCServe("8800")
	// lis, _ = RPCServeWithCode("8800")
}

func TestC(t *testing.T) {
	return
	t.Parallel()
	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)

	addr, err := net.ResolveTCPAddr("tcp", rpc_tcp_server)
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
	return
	t.Parallel()
	c := RPCClientWithCodec(rpc_tcp_server)
	defer c.Close()
	in := []byte("#   [Hi](http://mdblog.daoapp.io/)")
	out := make([]byte, 10)
	Markdown(c, &in, &out)
	fmt.Println(goutils.ToString(out))
	// defer lis.Close()
}

func TestJob(t *testing.T) {
	t.Parallel()
	// return
	c := RPCClientWithCodec(rpc_tcp_server)
	defer c.Close()
	// in := Job{Name: "google", Target: "https://www.google.com/search?q=docker&oq=docker&aqs=chrome..69i57j69i60l4.1517j0j4&sourceid=chrome&ie=UTF-8"}
	// in := Job{Name: "docker", Target: "https://www.google.com/aclk?sa=l&ai=Cls_gs3pMV8qDGcuq9gWdmIDoBprG9PcJgp_wutwC7MbsFAgAEAFgibPGhPQToAG8ppTsA8gBAaoEJk_QWc_UQjNHKe39e-t7guvDvFwTnmO55c8m1AmJHWa40wgOFtPpgAes2esTkAcBqAemvhvYBwE&sig=AOD64_1CJPrrkYMy7tE2O8DDMb1KMt8GuQ&clui=0&ved=0ahUKEwiagI7aq4LNAhXhnaYKHX5UDeQQ0QwIEg&adurl=https://circleci.com/integrations/docker/"}
	in := Job{Name: "mdblog", Target: "http://mdblog.daoapp.io"}
	out := make([]byte, 10)
	err := c.Call("RPC.Job", &in, &out)
	goutils.CheckErr(err)
	fmt.Println(in)
	fmt.Println(goutils.ToString(out))
}

func TestWall(t *testing.T) {
	return
	t.Parallel()
	c := RPCClientWithCodec(rpc_tcp_server)
	defer c.Close()
	out := Job{}
	in := make([]byte, 1)
	err := c.Call("RPC.Wall", &in, &out)
	goutils.CheckErr(err)
	fmt.Println("out:", out)
}

func TestWallBack(t *testing.T) {
	// return
	t.Parallel()
	c := RPCClientWithCodec(rpc_tcp_server)
	defer c.Close()
	job := Job{Name: "mdblog", Result: goutils.ToByte("mdblog-result")}
	out := make([]byte, 10)
	// in := make([]byte, 1)
	err := c.Call("RPC.WallBack", &job, &out)
	goutils.CheckErr(err)
	fmt.Println("out:", goutils.ToString(out))
}
