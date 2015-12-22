package main

import (
	"github.com/everfore/rpcsv"
	"github.com/shaalx/goutils"
)

func main() {
	lis, err := rpcsv.RPCServe("8800")
	if goutils.CheckErr(err) {
		return
	}
	defer lis.Close()
	<-make(chan bool)
}
