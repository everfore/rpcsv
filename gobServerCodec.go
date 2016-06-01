package rpcsv

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/shaalx/goutils"
	"io"
	"net"
	"net/rpc"
)

type gobServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	closed bool
}

func (c *gobServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return TimeoutCoder(c.dec.Decode, r, "server read request header")
}

func (c *gobServerCodec) ReadRequestBody(body interface{}) error {
	return TimeoutCoder(c.dec.Decode, body, "server read request body")
}

func (c *gobServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = TimeoutCoder(c.enc.Encode, r, "server write response"); err != nil {
		if c.encBuf.Flush() == nil {
			fmt.Println("rpc: gob error encoding response:", err)
			c.Close()
		}
		return
	}
	if err = TimeoutCoder(c.enc.Encode, body, "server write response body"); err != nil {
		if c.encBuf.Flush() == nil {
			fmt.Println("rpc: gob error encoding body:", err)
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *gobServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

func RPCServeWithCode(port string) (net.Listener, error) {
	rpc.Register(&RPC{})
	lis, err := net.Listen("tcp", ":"+port)
	goutils.CheckErr(err)
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				fmt.Println("Error: accept rpc connection", err.Error())
				continue
			}
			go func(conn net.Conn) {
				buf := bufio.NewWriter(conn)
				srv := &gobServerCodec{
					rwc:    conn,
					dec:    gob.NewDecoder(conn),
					enc:    gob.NewEncoder(buf),
					encBuf: buf,
				}
				defer srv.Close()
				err = rpc.ServeRequest(srv)
				if err != nil {
					fmt.Println("Error: server rpc request", err.Error())
				}
			}(conn)
		}
	}()
	return lis, nil
}
