package test

import (
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
)

func TestSimpleCounter(t *testing.T) {
	page := rod.New().MustConnect().MustPage("http://165.22.91.102:8080/3373bd0b-788b-43d7-8b10-3023ef4bd267/")
	time.Sleep(50 * time.Second)
	termElem := page.MustElement("#terminal")
	assert.Equal(t, "counting: 1\ncounting: 2\ncounting: 3\ncounting: 4\ncounting: 5\n", termElem.MustText())
}
