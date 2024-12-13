package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func run(screen *Screen, commandName string, args []string) {
	cmd := exec.Command(commandName, args...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		screen.text = fmt.Sprintf("could not connect to command-stdout: %v\n", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		screen.text = fmt.Sprintf("could not start command: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("got output line from command: %v\n", scanner.Text())
		screen.text += scanner.Text() + "\n"
	}
	// chunk := make([]byte, 150)
	// for
	// n, err := reader.Read(chunk)
	// if err != nil {
	// 	screen.text = fmt.Sprintf("could not read command output: %v\n", err)
	// 	return
	// }
	// line := chunk[0:n]
	// screen.text += string(line)
}
