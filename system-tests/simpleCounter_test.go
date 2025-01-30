//go:build system_test
// +build system_test

package test

import (
	"bufio"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/stretchr/testify/assert"
)

const debug = false

func TestSimpleCounter(t *testing.T) {
	go run("/bin/bash", "./simple-counter.sh", "0", "5")
	time.Sleep(2 * time.Second)
	var browser *rod.Browser
	if debug {
		l := launcher.New().
			Headless(false).
			Devtools(true)

		defer l.Cleanup()

		url := l.MustLaunch()

		// Trace shows verbose debug information for each action executed
		// SlowMotion is a debug related function that waits 2 seconds between
		// each action, making it easier to inspect what your code is doing.
		browser = rod.New().ControlURL(url).Trace(true).SlowMotion(2 * time.Second).MustConnect()

		// ServeMonitor plays screenshots of each tab. This feature is extremely
		// useful when debugging with headless mode.
		// You can also enable it with flag "-rod=monitor"
		launcher.Open(browser.ServeMonitor(""))

		defer browser.MustClose()
	} else {
		l := launcher.New().
			NoSandbox(true).
			Headless(true)
		defer l.Cleanup()

		url := l.MustLaunch()
		browser = rod.New().ControlURL(url).MustConnect()
		// browser = rod.New().MustConnect()
		defer browser.MustClose()
	}
	page := browser.MustPage(genUrl("/"))
	termElem, _ := page.Timeout(3 * time.Minute).Element("#terminal")
	html, err := page.HTML()
	assert.Nil(t, err)
	assert.Equal(t, "somehtml", html)
	assert.Equal(t, "counting: 2\ncounting: 2\ncounting: 3\ncounting: 4\ncounting: 5\n", termElem.MustText())
}

func run(command ...string) {
	// TODO instead of maxRunCount we could send INT/TERM signal to subprocess when done
	cmd := exec.Command("go", append([]string{"run", ".", "-maxRunCount", "40", "--"}, command...)...)
	cmd.Dir = ".."
	reader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("E: getting stderr for rwatch: %v --\n", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Printf("E: starting rwatch: %v --\n", err)
		return
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("got output line from command: %v\n", scanner.Text())
	}
	cmd.Wait()
}

func genUrl(relPath string) string {
	return fmt.Sprintf("http://165.22.91.102:8080/d43981bd-3822-4127-8cec-662f9a4d54f0%s", relPath)
}
