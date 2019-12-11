package ratelimit

import (
	"io"
	"time"
)

const bufferLength int = 1024

func computeTickRate(rate int) time.Duration {
	return time.Duration(int64(time.Second) / (int64(rate) / int64(bufferLength)))
}

// Copy will write from dst to src until EOF at a speed of no more than "rate" bytes per second.
func Copy(dst io.Writer, src io.Reader, rate int) (written int64, err error) {
	t := time.NewTicker(computeTickRate(rate))
	defer t.Stop()

	buf := make([]byte, bufferLength)

	for {
		<-t.C
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
