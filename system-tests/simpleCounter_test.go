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
	"github.com/stretchr/testify/assert"
)

const MAX_SESSIONID_WAIT_TIME = 45
const TIMEOUT_TERM_RESULT = 65 * time.Second

func TestSimpleCounter(t *testing.T) {
	t.Parallel()
	var sessionId string
	go run(&sessionId, "/bin/bash", "./simple-counter.sh", "0", "5")

	waitForSessionId(t, &sessionId)

	page, err := DefaultBrowserHandler.launchBrowserTab()
	if err != nil {
		t.Fatalf("E: while launching browser-tab: %v\n", err)
	}
	defer page.MustClose()
	page.MustNavigate(genUrl(sessionId, "/"))

	// term
	termElem, err := page.Timeout(TIMEOUT_TERM_RESULT).Element("#terminal")
	assert.Nil(t, err)
	assert.Equal(t, "counting: 1\ncounting: 2\ncounting: 3\ncounting: 4\ncounting: 5\n", termElem.MustText())

	// status
	checkStatus(t, page, "SUCCESS")
}

func TestSimpleFailureHandling(t *testing.T) {
	t.Parallel()
	var sessionId string
	go run(&sessionId, "/bin/bash", "./simple-failure.sh")

	waitForSessionId(t, &sessionId)

	page, err := DefaultBrowserHandler.launchBrowserTab()
	if err != nil {
		t.Fatalf("E: while launching browser-tab: %v\n", err)
	}
	defer page.MustClose()
	page.MustNavigate(genUrl(sessionId, "/"))

	// terminal
	termElem, err := page.Timeout(TIMEOUT_TERM_RESULT).Element("#terminal")
	assert.Nil(t, err)
	assert.Equal(t, "ERR: while running script", strings.TrimSpace(termElem.MustText()))

	// status
	checkStatus(t, page, "FAILED") // TODO
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
		// fmt.Printf("got output line from command: %v\n", line)
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

func waitForSessionId(t *testing.T, sessionId *string) {
	for i := 0; i < MAX_SESSIONID_WAIT_TIME; i++ {
		if *sessionId != "" {
			t.Logf("got session-id: %s\n", *sessionId)
			break
		}
		time.Sleep(time.Second)
	}
	t.Logf("got session-id: %s\n", *sessionId)
	if *sessionId == "" {
		t.Fatal("could not figure out session-if")
	}
}

func checkStatus(t *testing.T, page *rod.Page, expectedStatus string) {
	statusElem, err := page.Timeout(2 * time.Second).Element("#status")
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, statusElem.MustText())
}
