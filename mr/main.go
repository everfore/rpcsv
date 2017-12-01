package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/everfore/exc"
	"github.com/everfore/rpcsv"
	"github.com/toukii/bytes"
	"github.com/toukii/goutils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// rpc_tcp_server = "tcphub.t0.daoapp.io:61142"
	rpc_tcp_server = "127.0.0.1:8800"
	buf            []byte
	bytesWriter    *bytes.Writer

	tplFile = kingpin.Arg("tpl", "mr file.tpl").Default("README.md").String()
	outFile = kingpin.Flag("out", "mr -o render").Short('o').String()
)

func init() {
	// rpcsv.UpdataTheme()

	buf = make([]byte, 1024)
	bytesWriter = bytes.NewWriter(buf)
}

func main() {
	kingpin.Parse()

	if *tplFile == "" {
		fmt.Println("tpl file is nil.")
		return
	}
	if *outFile == "" {
		*outFile = renderfileName(*tplFile, "Render-", ".html")
	}

	in := goutils.ReadFile(*tplFile)
	if len(in) <= 0 {
		fmt.Println("tplFile with nil content")
		return
	}
	out, err := os.OpenFile(*outFile, os.O_CREATE|os.O_WRONLY, 0644)
	if goutils.CheckErr(err) {
		return
	}
	err = RenderWriter(in, out)
	goutils.CheckErr(err)
	exc.NewCMD("open -b com.google.Chrome " + *outFile).Debug().Execute()
}

func renderfileName(filename, prefix, suffix string) string {
	names := strings.Split(filename, ".")
	return strings.Join([]string{
		prefix,
		names[0],
		suffix,
	}, "")
}

func RenderWriter(in []byte, wr io.Writer) (err error) {
	client := rpcsv.RPCClient(rpc_tcp_server)
	defer client.Close()

	out := make([]byte, 0, 100)
	err = rpcsv.Markdown(client, &in, &out)
	if goutils.CheckErr(err) {
		return
	}

	data := make(map[string]interface{})
	data["MDContent"] = template.HTML(goutils.ToString(out))
	err = rpcsv.Theme.Execute(wr, data)
	goutils.CheckErr(err)

	return
}

func RenderedBytes(in []byte) ([]byte, error) {
	client := rpcsv.RPCClient(rpc_tcp_server)
	defer client.Close()

	out := make([]byte, 0, 1024)
	err := rpcsv.Markdown(client, &in, &out)
	if goutils.CheckErr(err) {
		return nil, err
	}

	bytesWriter.Reset()

	data := make(map[string]interface{})
	data["MDContent"] = template.HTML(goutils.ToString(out))
	err = rpcsv.Theme.Execute(bytesWriter, data)
	if goutils.CheckErr(err) {
		return nil, err
	}
	return bytesWriter.Bytes(), nil
}
