// +build go1.10

package main

import (
	"os"
	"time"
)

func readTimeout(f *os.File, t time.Time) error {
	return f.SetReadDeadline(t)
}
