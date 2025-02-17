//go:build system_test
// +build system_test

package test

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	fmt.Println("closing browser ...")
	err := DefaultBrowserHandler.browser.Close()
	if err != nil {
		fmt.Printf("E: while closing browser: %v\n", err)
	}
}
