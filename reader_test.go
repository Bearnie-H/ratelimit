package ratelimit

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

func TestReader(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	buf.WriteString("Reader Test! ")
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString(buf.String())
	buf.WriteString("\n")

	R := NewRateReader(buf, 1<<10)

	time.Sleep(1.5e9)

	if _, err := io.Copy(ioutil.Discard, R); err != nil {
		t.Fatalf("Error while copying - %s", err)
	}
}
