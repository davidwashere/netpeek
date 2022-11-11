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
$ netpeek -h
NetPeek captures UDP traffic periodically and summarizes captured metadata

Usage:
  netpeek [flags]

Flags:
  -d, --duration int    Number of seconds to capture traffic (default 5)
  -h, --help            help for netpeek
  -i, --interval int    Number of seconds to wait between captures, set to -1 to execute once and exit (default 300)
  -o, --output string   Defines where to send results, valid options: 'stdout', 'http[s]://...', or  will send 'path/to/some/file.log' (default "stdout")
  -p, --port string     Port to watch for traffic (default "34197")
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
