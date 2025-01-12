package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	screen := WebRTCScreen{}

	runnerDone := make(chan bool, 1)
	go run(&screen, runnerDone, "/bin/bash", []string{"./simple-counter.sh", "0", "5"})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// TODO handle shutdown
	// go func() {
	// 	for sig := range c {
	// 	}
	// }()
	select {
	case <-c:
	case <-runnerDone:
	}
	fmt.Println("good bye!")
}
