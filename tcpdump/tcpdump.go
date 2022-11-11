package tcpdump

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davidwashere/netpeek/domain"
)

const (
	commandNameDefault      = "tcpdump"
	packetSnapLengthDefault = "96"
)

var (
	// regexFindIpPort matches #.#.#.#.# - ie: 192.168.1.1.8080
	// finds #.#.#.#.# (last # is the port) typical in a tcpdump stdout msg
	regexFindIpPort = regexp.MustCompile(`(\d){1,3}\.(\d){1,3}\.(\d){1,3}\.(\d){1,3}\.(\d){1,5}`)
)

type UDPWatcherService struct {
	// CommandName is the binary to execute that produced `tcpdump` compatible output
	// will default to `tcpdump`
	CommandName string

	// PacketSnapLength defines the number of bytes to capture of each packet, defaults to 96
	// to get the full packet header(s) while ignoring packet data
	PacketSnapLength string
}

func NewUDPWatcherService() *UDPWatcherService {

	return &UDPWatcherService{
		CommandName:      commandNameDefault,
		PacketSnapLength: packetSnapLengthDefault,
	}
}

// parseLine will parse a line of output from tcpdump and extract the source
// ip for the packet
//
// this assumes that the tcpdump filter is set so that only traffic to a specific
// port is being considered (vs. bidirectional)
func (u *UDPWatcherService) parseLine(line string) (*domain.NetTuple, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("line is empty")
	}

	matches := regexFindIpPort.FindAllString(line, 3)
	if len(matches) != 2 {
		return nil, fmt.Errorf("too few or too many ip.port matches found")
	}

	r := domain.NetTuple{}
	r.NumPackets = 1

	match := matches[0]
	i := strings.LastIndexByte(match, '.')
	r.SrcIP = match[0:i]
	r.SrcPort = match[i+1:]

	match = matches[1]
	i = strings.LastIndexByte(match, '.')
	r.DestIP = match[0:i]
	r.DestPort = match[i+1:]

	return &r, nil
}
