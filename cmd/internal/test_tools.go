package internal

import (
	"bytes"
	"os"
)

func CaptureStdout(runner func()) string {
	// Save the original stdout
	originalStdout := os.Stdout

	// Create a read/write pipe
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the code that prints
	runner()

	// Close writer and restore original stdout
	w.Close()
	os.Stdout = originalStdout

	// Read and return output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}
