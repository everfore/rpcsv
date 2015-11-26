package main

import (
	"github.com/shaalx/goutils"
	md "github.com/shurcooL/github_flavored_markdown"
	"net"
	"net/rpc"
)

type RPC struct {
}

func (r *RPC) Markdown(in, out *([]byte)) error {
	*out = md.Markdown(*in)
	return nil
}

func main() {
	_rpc := new(RPC)
	server := rpc.NewServer()
	server.Register(_rpc)
	lis, err := net.Listen("tcp", ":8099")
	if goutils.CheckErr(err) {
		return
	}
	server.Accept(lis)
}
