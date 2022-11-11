package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davidwashere/netpeek/domain"
	"github.com/davidwashere/netpeek/tcpdump"

	"github.com/spf13/cobra"
)

const (
	stdoutOutput = "stdout"
)

var (
	interval int
	duration int
	port     string
	output   string

	rootCmd = &cobra.Command{
		Use:   "netpeek",
		Short: "NetPeek captures UDP traffic periodically and summarizes captured metadata",
		Run: func(cmd *cobra.Command, args []string) {
			execute()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "34197", "Port to watch for traffic")
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 5, "Number of seconds to capture traffic")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 300, "Number of seconds to wait between captures, set to -1 to execute once and exit")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", "Defines where to send results, valid options: 'stdout', 'http[s]://...', or  will send 'path/to/some/file.log'")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute() {
	for {
		executeCapture()

		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Second)
		} else if interval == -1 {
			return
		}

		// otherwise continually
	}
}

func executeCapture() {
	watcher := tcpdump.NewUDPWatcherService()

	log.Printf("capture started...")
	r, err := watcher.Watch(port, duration)
	log.Printf("capture finished")
	if err != nil {
		panic(err)
	}

	shareResults(r)
}

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
	dataB, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", dataB)
}

func shareResultsHttp(results *domain.Result) {
	dataB, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(dataB)

	url := output
	resp, err := http.Post(url, "application/json; charset=UTF-8", buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	log.Printf("response status: %s", resp.Status)
	log.Printf("response body: %s", body)
}

func shareResultsFile(results *domain.Result) {
	dataB, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	dataB = append(dataB, '\n')

	f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(dataB); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
}
