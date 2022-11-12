package util

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

// ExecuteCmd will execute the provided command, each line in stdout will be written to lineChan.
// When `timeout` seconds have elapsed the cmd is killed and method exits
//
// TODO: differentiate between stdout and stdin
//
// Implied:
// `cmdParts[0]` = command name
// `cmdParts[1:]...` = command args
func ExecuteCmd(cmdParts []string, timeout int, lineChan chan string) error {
	var err error

	// quitChan is used to stop writes to lineChan after cmd has been killed
	// this is necessary because the 'kill' issued to the OS is not instant
	// and there may be writes to lineChan after caller has closed lineChan
	// which triggers a panic
	quitChan := make(chan struct{})

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.SysProcAttr = CmdProcAttrs()
	stdout, _ := cmd.StdoutPipe()
	// stderr, _ := cmd.StderrPipe()

	var wg sync.WaitGroup

	wg.Add(1)
	go cap(&wg, quitChan, lineChan, stdout)

	// wg.Add(1)
	// go cap(&wg, quitChan, lineChan, stderr)

	// Run the Command
	err = cmd.Start()
	if err != nil {
		return err
	}

	finishedChan := make(chan struct{})
	go func() {
		defer close(finishedChan)

		// wait for routines to finish reading stderr/out buffers
		wg.Wait()

		// wait for the command to finish - in theory this should never actually wait
		err := cmd.Wait()
		if err != nil {
			log.Printf("error cmd.Wait: %v", err)
		}
	}()

	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	// TODO: Believe for loop here redundant
	// Prior to returning, ensure goroutines are done
	for {
		select {
		case <-finishedChan:
			return nil
		case <-ticker.C:
			ticker.Stop()
			close(quitChan) // kill the goroutines attempting to write to lineChan
			CmdKillNicely(cmd)
			return nil
		}
	}
}

func cap(wg *sync.WaitGroup, quitChan chan struct{}, lineChan chan string, stream io.ReadCloser) {
	defer wg.Done()
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()

		select {
		case _, ok := <-quitChan:
			if !ok {
				break
			}
		default:
			lineChan <- line
		}
	}

	if scanner.Err() != nil {
		log.Printf("error reading from stream: %v", scanner.Err())
	}
}

// ConvertOctalStringToUint32 interprets a string as an octal "644" and converts to decimal uin32 (420)
func ConvertOctalStringToUint32(in string) uint32 {
	r64, err := strconv.ParseUint(in, 8, 32)
	if err != nil {
		return 0
	}

	return uint32(r64)
}
