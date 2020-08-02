package ratelimit

import (
	"errors"
	"io"
	"sync"
)

// Package errors for the ReadPool and WritePool structs
var (
	// The pool could not allocate enough bandwidth for the request.
	ErrPoolExhausted error = errors.New("ratelimit pool error: Insufficient capacity to allocate new reader or writer")
)

// Don't allow allocations of smaller than 1024 bytes per second.
// If arbitrarily small allocations are allowed, the overhead to
// manage these pools outweighs the benefits.
//
// This must be some power of 2.
const (
	minAllocationSize int = 1 << 10
)

// ReadPool allows for a set amount of IO bandwidth to be explicitly allocatable and shareable
// between io.Readers as required by an application.
type ReadPool struct {
	maxAllocation int

	mu                *sync.Mutex
	currentAllocation int
	readers           map[*RateReader]int
}

// NewReadPool will...
func NewReadPool(TotalRate int) *ReadPool {

	TotalRate = normalizeRate(TotalRate)

	return &ReadPool{
		maxAllocation:     TotalRate,
		mu:                &sync.Mutex{},
		currentAllocation: 0,
		readers:           make(map[*RateReader]int),
	}
}

// NewReader will attempt to allocate a new RateReader from the pool with the requested maximum rate.
func (R *ReadPool) NewReader(r io.Reader, Rate int) (*RateReader, error) {
	R.mu.Lock()
	defer R.mu.Unlock()

	Rate = normalizeRate(Rate)

	if R.maxAllocation-R.currentAllocation >= Rate {
		RR := NewRateReader(r, Rate)
		R.readers[RR] = Rate
		R.currentAllocation += Rate
		return RR, nil
	}

	return nil, ErrPoolExhausted
}

// ReleaseReader will release the bandwidth associated with the given RateReader back to the pool.
// This should be called as soon as the reader is no longer required.
func (R *ReadPool) ReleaseReader(r *RateReader) {
	R.mu.Lock()
	defer R.mu.Unlock()

	if _, exist := R.readers[r]; exist {
		R.currentAllocation -= r.rate
		delete(R.readers, r)
	}
}

// WritePool will...
type WritePool struct {
	maxAllocation int

	mu                *sync.Mutex
	currentAllocation int
	writers           map[io.Writer]int
}

// NewWritePool will...
func NewWritePool(TotalRate int) *WritePool {

	TotalRate = normalizeRate(TotalRate)

	return &WritePool{
		maxAllocation:     TotalRate,
		mu:                &sync.Mutex{},
		currentAllocation: 0,
		writers:           make(map[io.Writer]int),
	}
}

// NewWriter will attempt to allocate a new RateWriter from the pool with the requested maximum rate.
func (W *WritePool) NewWriter(w io.Writer, Rate int) (*RateWriter, error) {
	W.mu.Lock()
	defer W.mu.Unlock()

	Rate = normalizeRate(Rate)

	if W.maxAllocation-W.currentAllocation >= Rate {
		RW := NewRateWriter(w, Rate)
		W.writers[RW] = Rate
		W.currentAllocation += Rate
		return RW, nil
	}

	return nil, ErrPoolExhausted
}

// ReleaseWriter will release the bandwidth associated with the given RateWriter back to the pool.
// This should be called as soon as the writer is no longer required.
func (W *WritePool) ReleaseWriter(w *RateWriter) {
	W.mu.Lock()
	defer W.mu.Unlock()

	if _, exist := W.writers[w]; exist {
		W.currentAllocation -= w.rate
		delete(W.writers, w)
	}
}

func normalizeRate(Rate int) int {
	switch {
	case Rate <= 0:
		return minAllocationSize
	case Rate&(minAllocationSize-1) != 0:
		return ((Rate & ^(minAllocationSize - 1)) + 1) * minAllocationSize
	default:
		return Rate
	}
}
