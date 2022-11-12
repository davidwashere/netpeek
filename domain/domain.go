package domain

import (
	"time"
)

const (
	StatusError   = "error"
	StatusSuccess = "success"

	DirectionBoth = "both"
	DirectionSrc  = "src"
	DirectionDest = "dst"
)

// UDPWatcherService is responsible for providing insights into system UDP activity
type UDPWatcherService interface {
	// Will watch UDP traffic sent to `destPort` port for `durationSecs` seconds
	// and return the Result
	Watch(destPort string, durationSecs int, direction string) (*Result, error)
}

type Result struct {
	// Status indicates whether or not the capture was successful
	Status string `json:"status"`

	// StatusMsg holds details about a failed result
	StatusMsg string `json:"statusMsg"`

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

func NewErrorResult(errMsg string) *Result {
	return &Result{
		Status:    StatusError,
		StatusMsg: errMsg,
		Start:     time.Now().UTC(),
		End:       time.Now().UTC(),
		UDPTuples: []*NetTuple{},
	}
}
