package tcpdump

import (
	"testing"
)

func compareStr(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestParseLine(t *testing.T) {
	srv := NewUDPWatcherService()

	r, err := srv.parseLine("20:29:47.116365 IP 98.32.177.181.60973 > 10.128.0.54.34197: UDP, length 14")
	if err != nil {
		t.Error(err)
	}

	if r == nil {
		t.Error("unexpected nil result")
		return
	}

	compareStr(t, r.SrcIP, "98.32.177.181")
	compareStr(t, r.SrcPort, "60973")
	compareStr(t, r.DestIP, "10.128.0.54")
	compareStr(t, r.DestPort, "34197")

	// Expect error when only one IP/Port combo found
	r, err = srv.parseLine("20:29:47.116365 IP 98.32.177.181.60973 > 10.128.0.54: UDP, length 14")
	if err == nil {
		t.Errorf("expected error but got %v", r)
	}

	// Expect error when input empty
	r, err = srv.parseLine("")
	if err == nil {
		t.Errorf("expected error but got %v", r)
	}
}

func TestBuildFilter(t *testing.T) {
	tt := []struct {
		port string
		dir  string
		want string
	}{
		{"1", "both", "udp and port 1"},
		{"1", "src", "udp and src port 1"},
		{"1", "dst", "udp and dst port 1"},
	}

	for _, test := range tt {
		got := buildFilter(test.port, test.dir)

		if got != test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}

}
