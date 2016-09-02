package rpcsv

import (
	"bufio"
	"encoding/gob"
	"github.com/toukii/goutils"
	"io"
	"net"
	"net/rpc"
	"time"
)

type RPCgobClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}

func (c *RPCgobClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = TimeoutCoder(c.enc.Encode, r, "client write request"); err != nil {
		return
	}
	if err = TimeoutCoder(c.enc.Encode, body, "client write request body"); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *RPCgobClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *RPCgobClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *RPCgobClientCodec) Close() error {
	return c.rwc.Close()
}

func RPCClientWithCodec(tcp_addr string) *rpc.Client {
	conn, err := net.DialTimeout("tcp", tcp_addr, time.Second*5)
	if goutils.CheckErr(err) {
		return nil
	}
	encBuf := bufio.NewWriter(conn)
	codec := &RPCgobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
	c := rpc.NewClientWithCodec(codec)
	return c
}
