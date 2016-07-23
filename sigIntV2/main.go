package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func WaitForCtrlC() {
	var wg sync.WaitGroup
	wg.Add(1)
	
	var signalCh chan os.Signal
	signalCh = make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	// signal.Notify(signalCh, os.Interrupt, os.Kill)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGTRAP)

	go func() {
		<-signalCh

		fmt.Printf("\n processing signal")

		wg.Done()
	}()
	wg.Wait()
	os.Exit(1)
}

func main() {
	go WaitForCtrlC()
	fmt.Printf("Press Ctrl+C to end\n")
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Printf("\n Doing my work")
		}

	}
	fmt.Printf("\n end of programm")
}
