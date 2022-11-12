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
      --dir string      direction to watch port traffic on, valid options: 'both', 'src', or 'dst'. 'src' means capture packets on this host that came from the specified port, 'dst' means capture packets sent to the specified port (default "both")
  -d, --duration int    number of seconds to capture traffic (default 5)
  -h, --help            help for netpeek
  -i, --interval int    number of seconds to wait between captures, set to -1 to execute once and exit (default 300)
  -o, --output string   where to send results, valid options: 'stdout', 'http[s]://...', or  will send 'path/to/some/file.log' (default "stdout")
      --perm string     file permissions to set when writing results to a file (default "644")
  -p, --port string     port to watch for traffic (default "34197")
      --pretty          pretty print result json when output == stdout
```

## Improvements
- [ ] Migrate to `https://github.com/google/gopacket` from `tcpdump`
  - enables windows support
- [ ] Add ability to capture established TCP connections
    - Consider: https://github.com/weaveworks/procspy
- [ ] Add Auth to http/https dst's (basic auth, bearer, etc.)
- [ ] Test keeping handle to file output open instead of closing