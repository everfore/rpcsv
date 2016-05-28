package rpcsv

import (
	"github.com/shaalx/goutils"
	md "github.com/shurcooL/github_flavored_markdown"
	"net"
	"net/rpc"

	// "bytes"
	// "fmt"
	"html/template"
	"os"
)

type RPC struct {
}

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

func (r *RPC) Markdown(in, out *([]byte)) error {
	// fmt.Println(goutils.ToString(*in))
	html := md.Markdown(*in)
	goutils.ReWriteFile("tempory.tmp", nil)
	of, _ := os.OpenFile("tempory.tmp", os.O_CREATE|os.O_WRONLY, 0666)
	defer of.Close()
	data := make(map[string]interface{})
	data["MDContent"] = template.HTML(goutils.ToString(html))
	err := theme.Execute(of, data)
	if goutils.CheckErr(err) {
		return err
	}
	*out = goutils.ReadFile("tempory.tmp")
	// fmt.Println(goutils.ToString(html))
	// fmt.Println("out:", goutils.ToString(*out))
	return nil
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
