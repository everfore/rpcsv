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
	RPC_Client = rpcsv.RPCClient("182.254.132.59:32769")
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
		return
	}
	tpl.Execute(rw, nil)
}

func markdown(rw http.ResponseWriter, req *http.Request) {
	/*tpl, err := template.New("index.html").ParseFiles("index.html")
	if goutils.CheckErr(err) {
		return
	}
	tpl.Execute(rw, nil)*/
	req.ParseForm()
	rawContent := req.Form.Get("rawContent")
	fmt.Println(rawContent)
	out := make([]byte, 10)
	in := goutils.ToByte(rawContent)
	err := rpcsv.Markdown(RPC_Client, &in, &out)
	if goutils.CheckErr(err) {
		return
	}
	mdbs := goutils.ToByte(template.HTMLEscapeString(goutils.ToString(out)))
	rw.Write(mdbs)
}
