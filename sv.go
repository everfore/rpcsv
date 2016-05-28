package rpcsv

import (
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"

	// "bytes"
	// "fmt"
	"html/template"
)

func NewBufWriter() *BufWriter {
	return &BufWriter{
		buf: make([]byte, 1000, 1024),
	}
}

type BufWriter struct {
	buf []byte
}

func (w *BufWriter) Write(p []byte) (n int, err error) {
	w.buf = p
	return len(w.buf), nil
}

func (w *BufWriter) Bytes() []byte {
	return w.buf
}

func RPCServe(port string) (net.Listener, error) {
	_rpc := new(RPC)
	server := rpc.NewServer()
	server.Register(_rpc)
	lis, err := net.Listen("tcp", ":"+port)
	if goutils.CheckErr(err) {
		return nil, err
	}
	go server.Accept(lis)
	return lis, nil
}

var (
	theme    *template.Template
	theme_bs []byte
)

func init() {
	theme_bs = goutils.ReadFile("theme.thm")
	if nil == theme_bs {
		theme_bs = goutils.ToByte(theme_s)
	}
	var err error
	theme, err = template.New("theme.thm").Parse(goutils.ToString(theme_bs))
	if err != nil {
		panic("theme error")
	}
}

func UpdataTheme() bool {
	theme_bs = goutils.ReadFile("theme.thm")
	if nil == theme_bs {
		theme_bs = goutils.ToByte(theme_s)
	}
	var err error
	theme, err = template.New("theme.thm").Parse(goutils.ToString(theme_bs))
	if err != nil {
		panic("theme error")
	}
	return true
}

const (
	theme_s = `{{.MDContent}}`
)
