# NetPeek

NetPeek will capture UDP packets at intervals, sending summary of packets observed to a destination (stdout, file, http, etc.) for further processing

_Linux Only at this time_

## Motivation

Detecting idle servers (ie: game servers) via stdout/logs is not always feasible. A more consistent way to determine if clients are connected or communicated was desired

## Prerequisites

- `tcpdump` installed and available on `PATH`
- executed as a user with privileges necessary to run `tcpdump` (typically `root`)

## Usage

```
sudo ./netpeek <flags>
```


## Scratch
```sh
tcpdump -s 96 -nn "udp and dst port 34197"
tcpdump -s 96 -nn "udp and port 34197"


{
    "start":"val",
    "end":"val",
    "udpTuples": [
        { "srcIP":"" , "srcPort":"" , "destIP":"" , "destPort":"", "numPackets": 13 },
        { "srcIP":"" , "srcPort":"" , "destIP":"" , "destPort":"", "numPackets": 7 },
    ]
}
```

## Improvements
- [ ] Add windows support
- [ ] Add ability to capture established TCP connections
    - Consider: https://github.com/weaveworks/procspy
- [ ] Migrate to `https://github.com/google/gopacket` so not dependant on `tcpdump` directly
    - map help with cross platform support (ie: windows)
- [ ] Add Auth to http/https dst's (basic auth, bearer, etc.)
- [ ] Add ability to customize capture filter / parameters
- [ ] Add flag to stream raw captured packet data
