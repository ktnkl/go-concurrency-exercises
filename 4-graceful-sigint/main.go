//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// Create a process
	proc := MockProcess{}

	go func() {
		proc.Run()
		close(done)
	}()

	select {
	case <-sigChan:
		stopDone := make(chan struct{})

		go func() {
			proc.Stop()
			close(stopDone)
		}()

		select {
		case <-sigChan:
			fmt.Println("Kill process")
			os.Exit(1)
		case <-stopDone:
			fmt.Println("Process stop itself")
		}

	case <-done:
		fmt.Println("Process stop itself")
	}
}
