package domain

import (
	"fmt"
	"time"
)

// UDPWatcherService
type UDPWatcherService interface {
	// Will watch UDP traffic sent to `destPort` port for `durationSecs` seconds
	// and return the Result
	Watch(destPort string, durationSecs int) (*Result, error)
}

type Result struct {
	// Start is the timestamp when the observation started in UTC
	Start time.Time `json:"start"`

	// Start is the timestamp when the observation ended in UTC
	End time.Time `json:"end"`

	// Tuples will contain the number of packets observed during duration
	// for each src ip:port and dst ip:port combination
	UDPTuples []*NetTuple `json:"udpTuples"`
}

type NetTuple struct {
	SrcIP      string `json:"srcIP"`
	SrcPort    string `json:"srcPort"`
	DestIP     string `json:"destIP"`
	DestPort   string `json:"destPort"`
	NumPackets int    `json:"numPackets"`
}

func (t *NetTuple) Key() string {
	return fmt.Sprintf("%s:%s->%s:%s", t.SrcIP, t.SrcPort, t.DestIP, t.DestPort)
}
