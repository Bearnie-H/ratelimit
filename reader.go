package ratelimit

import (
	"io"
	"time"
)

// RateReader extends an io.Reader to read at no more than rate bytes per second.
type RateReader struct {
	r        io.Reader
	rate     int
	lastRead time.Time
}

// NewRateReader will create an initialize a new RateReader to read from r at no more than
// perSecond bytes per second.
func NewRateReader(r io.Reader, perSecond int) *RateReader {
	return &RateReader{
		r:        r,
		rate:     perSecond,
		lastRead: time.Now(),
	}
}

// Read implements the io.Reader interface for RateReader.
func (R *RateReader) Read(p []byte) (n int, err error) {

	// Check how much room there is to read
	l := int((float64(time.Now().Sub(R.lastRead)) / float64(time.Second)) * float64(R.rate))

	// If there's no room, do nothing
	if l <= 0 {
		return 0, nil
	}

	// If there is room, defer a call to update when the last read occurred
	defer func() { R.lastRead = time.Now() }()

	// Read as many bytes as there's room for
	return R.r.Read(p[:l])
}

// Close implements the io.Closer interface for a RateReader,
// closing any held resources and calling the Close() method
// of the underlying io.Reader if it implements io.Closer.
func (R *RateReader) Close() error {

	// If the underlying writer implements io.Closer, close it
	if c, ok := R.r.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
