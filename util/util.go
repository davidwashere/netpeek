package util

import (
	"bufio"
	"log"
	"os/exec"
	"sync"
	"time"
)

// ExecuteCmd will execute the provided command, each line in stdout/stderr will be written to lineChan.
// When `timeout` seconds have elapsed the cmd is killed and method exits
//
// TODO: differentiate between stdout and stdin
//
// Implied:
// `cmdParts[0]` = command name
// `cmdParts[1:]...` = command args
func ExecuteCmd(cmdParts []string, timeout int, lineChan chan string) error {
	var err error

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.SysProcAttr = CmdProcAttrs()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			lineChan <- line
		}

		if scanner.Err() != nil {
			log.Printf("error reading stdout: %v", scanner.Err())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			lineChan <- line
		}

		if scanner.Err() != nil {
			log.Printf("error reading stderr: %v", scanner.Err())
		}
	}()

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
	for {
		select {
		case <-finishedChan:
			return nil
		case <-ticker.C:
			ticker.Stop()
			CmdKillNicely(cmd)
			return nil
		}
	}

}
