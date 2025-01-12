package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	screen := WebRTCScreen{}

	go run(&screen, "/bin/bash", []string{"./simple-counter.sh", "0", "5"})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// TODO handle shutdown
	// go func() {
	// 	for sig := range c {
	// 	}
	// }()
	<-c
	fmt.Println("good bye!")
}
