package rpcsv

import (
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
)

// func main() {
// 	in := []byte("#   [Hello](http://mdblog.daoapp.io/)")
// 	out := make([]byte, 10)

// 	// addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:80")
// 	addr, err := net.ResolveTCPAddr("tcp", "rpcsvr.daoapp.io:80")
// 	if goutils.CheckErr(err) {
// 		return
// 	}
// 	conn, err := net.DialTCP("tcp", nil, addr)
// 	defer conn.Close()
// 	if goutils.CheckErr(err) {
// 		return
// 	}

// 	rc := rpc.NewClient(conn)
// 	err = rc.Call("RPC.Markdown", &in, &out)
// 	if goutils.CheckErr(err) {
// 		return
// 	}
// 	fmt.Println(string(out))
// }

func RPCClient(tcp_addr string) *rpc.Client {
	addr, err := net.ResolveTCPAddr("tcp", tcp_addr)
	if goutils.CheckErr(err) {
		return nil
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if goutils.CheckErr(err) {
		return nil
	}

	rc := rpc.NewClient(conn)
	return rc
}

func Markdown(rc *rpc.Client, in, out *([]byte)) error {
	err := rc.Call("RPC.Markdown", in, out)
	if goutils.CheckErr(err) {
		return err
	}
	return nil
}
