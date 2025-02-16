//go:build system_test
// +build system_test

package test

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/stretchr/testify/assert"
)

const debug = false
const MAX_SESSIONID_WAIT_TIME = 45

func TestSimpleCounter(t *testing.T) {
	var sessionId string
	go run(&sessionId, "/bin/bash", "./simple-counter.sh", "0", "5")

	// wait for sessionId
	for i := 0; i < MAX_SESSIONID_WAIT_TIME; i++ {
		if sessionId != "" {
			t.Logf("got session-id: %s\n", sessionId)
			break
		}
		time.Sleep(time.Second)
	}
	t.Logf("got session-id: %s\n", sessionId)
	if sessionId == "" {
		t.Fatal("could not figure out session-if")
	}

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
	page := browser.MustPage(genUrl(sessionId, "/"))
	termElem, err := page.Timeout(65 * time.Second).Element("#terminal")
	assert.Nil(t, err)
	assert.Equal(t, "counting: 1\ncounting: 2\ncounting: 3\ncounting: 4\ncounting: 5\n", termElem.MustText())
}

func run(sessionId *string, command ...string) {
	// TODO instead of maxRunCount we could send INT/TERM signal to subprocess when done
	cmd := exec.Command("go", append([]string{"run", ".", "--max-run-count", "65", "--plain-text-screen", "--"}, command...)...)
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
		line := scanner.Text()
		fmt.Printf("got output line from command: %v\n", line)
		if strings.HasPrefix(strings.TrimSpace(line), "Session-ID:") {
			parsedSessionId := strings.TrimSpace(strings.TrimLeft(line, "Session-ID:"))
			*sessionId = parsedSessionId
		}
	}
	cmd.Wait()
}

func genUrl(sessionId, relPath string) string {
	return fmt.Sprintf("http://165.22.91.102:8080/%s%s", sessionId, relPath)
}
