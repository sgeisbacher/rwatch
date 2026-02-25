//go:build system_test
// +build system_test

package test

import (
	"bufio"
	"context"
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

	// setup
	page, done, err := createTestSetup("/bin/bash", "./simple-counter.sh", "0", "5")
	if err != nil {
		t.Fatalf("E: while creating test setup: %w", err)
	}

	// term
	termElem, err := page.Timeout(TIMEOUT_TERM_RESULT).Element("#terminal")
	assert.Nil(t, err)
	assert.Equal(t, "counting: 1\ncounting: 2\ncounting: 3\ncounting: 4\ncounting: 5\n", termElem.MustText())

	// status
	checkStatus(t, page, "SUCCESS")
	done()
}

func TestSimpleFailureHandling(t *testing.T) {
	t.Parallel()

	// setup
	page, done, err := createTestSetup("/bin/bash", "./simple-failure.sh")
	if err != nil {
		t.Fatalf("E: while creating test setup: %w", err)
	}

	// terminal
	termElem, err := page.Timeout(TIMEOUT_TERM_RESULT).Element("#terminal")
	assert.Nil(t, err)
	assert.Equal(t, "ERR: while running script", strings.TrimSpace(termElem.MustText()))

	// status
	checkStatus(t, page, "FAILED")
	done()
}

func createTestSetup(command ...string) (*rod.Page, func(), error) {
	var sessionId string
	ctx, cancelRun := context.WithCancel(context.Background())

	// run rwatch in the background
	go run(ctx, &sessionId, command...)

	err := waitForSessionId(&sessionId)
	if err != nil {
		return nil, cancelRun, fmt.Errorf("E: while waiting for session-id: %w", err)
	}

	page, err := DefaultBrowserHandler.launchBrowserTab()
	if err != nil {
		err = fmt.Errorf("E: while launching browser-tab: %w\n", err)
		return page, cancelRun, err
	}
	page.MustNavigate(genUrl(sessionId, "/"))
	doneFn := func() {
		cancelRun()
		page.MustClose()
	}
	return page, doneFn, err
}

func run(ctx context.Context, sessionId *string, command ...string) {
	// TODO instead of maxRunCount we could send INT/TERM signal to subprocess when done
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", ".", "--max-run-count", "65", "--plain-text-screen", "--"}, command...)...)
	cmd.Dir = ".."
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("E: getting stderr for rwatch: %v --\n", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Printf("E: starting rwatch: %v --\n", err)
		return
	}
	scanner := bufio.NewScanner(stdoutReader)
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

func waitForSessionId(sessionId *string) error {
	for i := 0; i < MAX_SESSIONID_WAIT_TIME; i++ {
		if *sessionId != "" {
			t.Logf("got session-id: %s\n", *sessionId)
			break
		}
		time.Sleep(time.Second)
	}
	if *sessionId == "" {
		return fmt.Errorf("could not figure out session-if")
	}
	return nil
}

func checkStatus(t *testing.T, page *rod.Page, expectedStatus string) {
	statusElem, err := page.Timeout(2 * time.Second).Element("#status")
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, statusElem.MustText())
}
