package main

import (
	"fmt"
	"github.com/everfore/rpcsv"
	"github.com/shaalx/goutils"
	"html/template"
	"net/http"
	"net/rpc"
)

var (
	RPC_Client *rpc.Client
)

func init() {
	connect()
}

func connect() {
	RPC_Client = rpcsv.RPCClient("182.254.132.59:8800")
}

func main() {
	defer RPC_Client.Close()
	http.HandleFunc("/", index)
	http.HandleFunc("/markdown", markdown)
	http.ListenAndServe(":80", nil)
}

func index(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.New("index.html").ParseFiles("index.html")
	if goutils.CheckErr(err) {
		rw.Write(goutils.ToByte(err.Error()))
		return
	}
	tpl.Execute(rw, nil)
}

func markdown(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	rawContent := req.Form.Get("rawContent")
	fmt.Println(rawContent)
	out := make([]byte, 0, 100)
	in := goutils.ToByte(rawContent)
	connect()
	defer RPC_Client.Close()
	err := rpcsv.Markdown(RPC_Client, &in, &out)
	if goutils.CheckErr(err) {
		rw.Write(goutils.ToByte(err.Error()))
		return
	}
	if len(out) <= 0 {
		rw.Write(goutils.ToByte("<h3>nil</h3>"))
	} else {
		rw.Write(out)
	}
}
