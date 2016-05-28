package rpcsv

import (
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
)

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
