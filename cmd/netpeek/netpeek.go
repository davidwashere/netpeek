package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davidwashere/netpeek/domain"
	"github.com/davidwashere/netpeek/tcpdump"
	"github.com/davidwashere/netpeek/util"

	"github.com/spf13/cobra"
)

const (
	stdoutOutput = "stdout"
)

var (
	interval  int
	duration  int
	port      string
	output    string
	perm      string
	pretty    bool
	direction string

	rootCmd = &cobra.Command{
		Use:   "netpeek",
		Short: "NetPeek captures UDP traffic periodically and summarizes captured metadata",
		Run: func(cmd *cobra.Command, args []string) {
			executeCaptureLoop()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "34197", "port to watch for traffic")
	rootCmd.PersistentFlags().StringVar(&direction, "dir", "both", "direction to watch port traffic on, valid options: 'both', 'src', or 'dst'. 'src' means capture packets on this host that came from the specified port, 'dst' means capture packets sent to the specified port")
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 5, "number of seconds to capture traffic")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 300, "number of seconds to wait between captures, set to -1 to execute once and exit")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", "where to send results, valid options: 'stdout', 'http[s]://...', or  will send 'path/to/some/file.log'")
	rootCmd.PersistentFlags().StringVar(&perm, "perm", "644", "file permissions to set when writing results to a file")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "pretty print result json when output == stdout")
}

// executeCaptureLoop will execute a UDP capture once, continuously, or repeatedly with a delay depending on config parameters
func executeCaptureLoop() {
	for {
		results := executeCapture()
		shareResults(results)

		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Second)
		} else if interval == -1 {
			return
		}

		// if here - continuos capture (no interval)
	}
}

// executeCapture will trigger the actual capture
func executeCapture() *domain.Result {
	watcher := tcpdump.NewUDPWatcherService()

	log.Printf("capture started...")
	r, err := watcher.Watch(port, duration, direction)
	if err != nil {
		r = domain.NewErrorResult(fmt.Sprintf("%v", err))
		log.Printf("error capturing packets: %v", err)
	}

	if output != stdoutOutput {
		// if we're not already writing to stdout, then log the result
		log.Printf("capture results: %v", r)
	}
	return r
}

// shareResults will print, log, or send the capture results depending on config parameters
func shareResults(results *domain.Result) {
	if output == stdoutOutput {
		shareResultsStdout(results)
		return
	}

	if strings.HasPrefix(output, "http://") || strings.HasPrefix(output, "https://") {
		shareResultsHttp(results)
		return
	}

	shareResultsFile(results)
}

func shareResultsStdout(results *domain.Result) {
	if pretty {
		fmt.Printf("%s\n---\n", results.PrettyString())
	} else {
		fmt.Printf("%s\n", results.String())
	}
}

func shareResultsHttp(results *domain.Result) {
	buf := bytes.NewBuffer(results.Bytes())

	url := output
	resp, err := http.Post(url, "application/json; charset=UTF-8", buf)
	if err != nil {
		log.Printf("error posting results to %v: %v", url, err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	log.Printf("response '%s': %s", resp.Status, body)
}

func shareResultsFile(results *domain.Result) {
	dataB := results.Bytes()
	dataB = append(dataB, '\n')

	p := util.ConvertOctalStringToUint32(perm)
	f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.FileMode(p))
	if err != nil {
		log.Printf("error opening file %v: %v", output, err)
		return
	}
	defer f.Close()

	if _, err := f.Write(dataB); err != nil {
		log.Printf("error writing data to file %v: %v", output, err)
		return
	}
}
