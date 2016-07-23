package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func WaitForCtrlC() {
	var wg sync.WaitGroup
	wg.Add(1)
	//
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		<-signalCh

		fmt.Printf("\nProcessing signal")

		wg.Done()
	}()
	wg.Wait()
}

func main() {

	fmt.Printf("Press Ctrl+C to end\n")
	go testOracleDB()

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				fmt.Printf("\n Doing my work")
			}

		}
	}()
	WaitForCtrlC()
	fmt.Printf("\n end of programm")
}
