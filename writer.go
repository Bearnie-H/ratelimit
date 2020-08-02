package ratelimit

import (
	"io"
	"time"
)

// RateWriter implements an extension to an io.Writer such that it will write
// at no more than a given maximum rate.
type RateWriter struct {
	w         io.Writer
	rate      int
	lastWrite time.Time
}

// NewRateWriter will wrap an io.Writer so that it will write no more than perSecond bytes
// per second to the underlying writer.
func NewRateWriter(w io.Writer, perSecond int) *RateWriter {

	return &RateWriter{
		w:         w,
		rate:      perSecond,
		lastWrite: time.Now(),
	}
}

// Write implements the io.Writer interface for the RateWriter struct.
func (W *RateWriter) Write(p []byte) (n int, err error) {

	// Update the last write time
	defer func() { W.lastWrite = time.Now() }()

	// Wait until there's "room" to write
	W.wait(len(p))

	// Perform the actual write
	return W.w.Write(p)
}

// Close extends the RateWriter to also implement the io.Closer interface.
// This will release any internal resources, as well as calling the Close()
// function of the underlying io.Writer if it implements io.Closer.
func (W *RateWriter) Close() error {

	// If the underlying writer implements io.Closer, close it
	if c, ok := W.w.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (W *RateWriter) wait(n int) {

	// Wait until the time since the last write is equal to or greater than
	// what would be required to write n bytes to the writer at a rate no greater
	// than W.rate.
	time.Sleep(
		time.Until(
			W.lastWrite.Add(
				time.Duration(float64(time.Second) * (float64(n) / float64(W.rate))),
			),
		),
	)
}
