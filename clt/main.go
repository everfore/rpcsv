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
	// RPC_Client = rpcsv.RPCClient("182.254.132.59:8800")
	RPC_Client = rpcsv.RPCClient("127.0.0.1:8800")
	if RPC_Client == nil {
		panic("nil RPC_Client")
	}
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
	fmt.Println(req.RemoteAddr, req.Referer())
	// fmt.Println(rawContent)
	out := make([]byte, 0, 100)
	in := goutils.ToByte(rawContent)
	err := rpcsv.Markdown(RPC_Client, &in, &out)
	if goutils.CheckErr(err) {
		rw.Write(goutils.ToByte(err.Error()))
		return
	}
	if len(out) <= 0 {
		rw.Write(goutils.ToByte("{response:nil}"))
		return
	}
	writeCrossDomainHeaders(rw, req)
	rw.Write(out)
}

func writeCrossDomainHeaders(w http.ResponseWriter, req *http.Request) {
	// Cross domain headers
	if acrh, ok := req.Header["Access-Control-Request-Headers"]; ok {
		w.Header().Set("Access-Control-Allow-Headers", acrh[0])
	}
	w.Header().Set("Access-Control-Allow-Credentials", "True")
	if acao, ok := req.Header["Access-Control-Allow-Origin"]; ok {
		w.Header().Set("Access-Control-Allow-Origin", acao[0])
	} else {
		if _, oko := req.Header["Origin"]; oko {
			w.Header().Set("Access-Control-Allow-Origin", req.Header["Origin"][0])
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Connection", "Close")

}
