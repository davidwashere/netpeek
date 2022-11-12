package tcpdump

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davidwashere/netpeek/domain"
)

const (
	// commandNameDefault default binary to execute to capture packets
	commandNameDefault = "tcpdump"

	// packetSnapLengthDefault the default number of bytes to capture per tcpdump packet, ideally this is the
	// minimum bytes necessary to extract relevant packet headers
	packetSnapLengthDefault = "96"
)

var (
	// regexFindIpPort matches #.#.#.#.# (the last # is the port) - ie: 192.168.1.1.8080
	regexFindIpPort = regexp.MustCompile(`(\d){1,3}\.(\d){1,3}\.(\d){1,3}\.(\d){1,3}\.(\d){1,5}`)
)

type UDPWatcherService struct {
	// CommandName is the binary to execute that produces `tcpdump` compatible output
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

// parseLine will extract relevant fields from `tcpdump` stdout
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

func buildFilter(port string, direction string) string {
	dir := ""
	if direction == domain.DirectionSrc {
		dir = "src "
	} else if direction == domain.DirectionDest {
		dir = "dst "
	}

	filter := fmt.Sprintf("udp and %sport %s", dir, port)
	return filter
}
