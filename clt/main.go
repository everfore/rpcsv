package main

import (
	"encoding/json"
	"fmt"
	"github.com/everfore/rpcsv"
	"github.com/shaalx/goutils"
	"html/template"
	"net/http"
	"net/rpc"
	"time"
)

var (
	RPC_Client *rpc.Client
)

func connect() {
	RPC_Client = rpcsv.RPCClient("182.254.132.59:8800")
	// RPC_Client = rpcsv.RPCClient("127.0.0.1:8800")
	go func() {
		time.Sleep(2e9)
		RPC_Client.Close()
	}()
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/markdown", markdown)
	http.HandleFunc("/markdownCB", markdownCB)
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
	times := 0
	connect()
retry:
	times++
	err := rpcsv.Markdown(RPC_Client, &in, &out)
	if goutils.CheckErr(err) {
		connect()
		if times < 6 {
			goto retry
		}
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
	// w.Header().Set("Connection", "Close")
}

func markdownCB(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	rawContent := req.Form.Get("rawContent")
	fmt.Println(req.RemoteAddr, req.Referer())
	// fmt.Println(rawContent)
	out := make([]byte, 0, 100)
	in := goutils.ToByte(rawContent)
	RPC_Client = rpcsv.RPCClient("182.254.132.59:8800")
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
	fmt.Println(req.RemoteAddr)
	CallbackFunc := fmt.Sprintf("CallbackFunc(%v);", string(Json(goutils.ToString(out))))
	fmt.Fprint(rw, CallbackFunc)
}

type CallbackData struct {
	Mddata interface{} `json:"mddata"`
}

func Json(data interface{}) []byte {
	bs, err := json.Marshal(CallbackData{Mddata: data})
	if goutils.CheckErr(err) {
		return nil
	}
	return bs
}
