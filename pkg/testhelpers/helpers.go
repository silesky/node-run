package testhelpers

import (
	"bytes"
	"io"
	"os"
)

// CaptureOutput captures stdout output
func CaptureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	outC := make(chan string)
	// Copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Run the function
	f()

	// Restore original stdout
	w.Close()
	os.Stdout = old
	return <-outC
}

// SimulateInput simulates stdin input
func SimulateInput(input string, f func()) {
	r, w, _ := os.Pipe()
	old := os.Stdin
	os.Stdin = r

	// Write the input in a separate goroutine
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Run the function
	f()

	// Restore original stdin
	os.Stdin = old
}
