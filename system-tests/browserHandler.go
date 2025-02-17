package test

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const debug = false

type BrowserHandler struct {
	mu      sync.Mutex
	browser *rod.Browser
}

var DefaultBrowserHandler BrowserHandler = BrowserHandler{}

func (bh *BrowserHandler) launchBrowserTab() (*rod.Page, error) {
	headless := os.Getenv("CI") == "true"
	bh.mu.Lock()
	defer bh.mu.Unlock()
	if bh.browser == nil {
		fmt.Println("launching browser ...")
		if debug {
			l := launcher.New().
				Headless(false).
				Devtools(headless)

			// defer l.Cleanup()

			url := l.MustLaunch()

			// Trace shows verbose debug information for each action executed
			// SlowMotion is a debug related function that waits 2 seconds between
			// each action, making it easier to inspect what your code is doing.
			bh.browser = rod.New().ControlURL(url).Trace(true).SlowMotion(2 * time.Second).MustConnect()

			// ServeMonitor plays screenshots of each tab. This feature is extremely
			// useful when debugging with headless mode.
			// You can also enable it with flag "-rod=monitor"
			launcher.Open(bh.browser.ServeMonitor(""))
		} else {
			l := launcher.New().
				NoSandbox(true).
				Headless(headless)
			// defer l.Cleanup()

			url := l.MustLaunch()
			bh.browser = rod.New().ControlURL(url).MustConnect()
			// browser = rod.New().MustConnect()
		}
	} else {
		fmt.Println("reusing browser ...")
	}
	return bh.browser.MustPage(), nil
}
