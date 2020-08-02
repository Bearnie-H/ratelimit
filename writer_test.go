package ratelimit

import (
	"bytes"
	"os"
	"testing"
)

func TestWriter(t *testing.T) {

	// Wrap the file to only write a small number of bytes per cycle
	W := NewRateWriter(os.Stdout, 8)

	buf := bytes.NewBuffer(nil)
	buf.WriteString("Writer Test!\n")

	// Perform 10 iterations to make sure the Write is waiting enough time.
	for i := 0; i < 10; i++ {
		_, err := W.Write(buf.Bytes())
		if err != nil {
			t.Fatalf("Copy failure - %s", err)
		}
	}
}
