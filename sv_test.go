package rpcsv

import (
	"testing"
)

func TestSV(t *testing.T) {
	rpc_ := new(RPC)
	in := []byte("#    Hello")
	out := make([]byte, 1)
	err := rpc_.Markdown(&in, &out)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(out))
}
