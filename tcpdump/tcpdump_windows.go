package tcpdump

import (
	"log"
	"time"

	"github.com/davidwashere/netpeek/domain"
)

func (u *UDPWatcherService) Watch(destPort string, duration int, direction string) (*domain.Result, error) {
	log.Printf("WARNING: faking tcpdump response for windows testing, NOT REAL DATA")
	// TODO: consider using mock instead of platform specific here

	result := domain.Result{}

	result.Start = time.Now().UTC().Add(-time.Duration(duration) * time.Second)
	result.End = time.Now().UTC()
	result.UDPTuples = append(result.UDPTuples,
		&domain.NetTuple{
			SrcIP:      "10.128.0.54",
			SrcPort:    "34197",
			DestIP:     "98.32.177.181",
			DestPort:   "54166",
			NumPackets: 59411,
		},
		&domain.NetTuple{
			SrcIP:      "98.32.177.181",
			SrcPort:    "54166",
			DestIP:     "10.128.0.54",
			DestPort:   "34197",
			NumPackets: 59319,
		},
	)

	return &result, nil

}
