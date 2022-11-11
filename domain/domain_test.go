package domain

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResultJSONMarshall(t *testing.T) {

	r := Result{
		UDPTuples: []*NetTuple{},
	}

	nt := new(NetTuple)
	nt.SrcIP = "a"
	nt.SrcPort = "1"
	nt.DestIP = "b"
	nt.DestPort = "2"
	nt.NumPackets = 1
	r.UDPTuples = append(r.UDPTuples, nt)

	nt = new(NetTuple)
	nt.SrcIP = "c"
	nt.SrcPort = "3"
	nt.DestIP = "d"
	nt.DestPort = "4"
	nt.NumPackets = 2
	r.UDPTuples = append(r.UDPTuples, nt)

	dataB, err := json.Marshal(&r)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s\n", dataB)
}
