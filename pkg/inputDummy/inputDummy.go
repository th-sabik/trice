// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package inputDummy simulates artificial trice data for testing.
//
// On creation a byte slice with mockup data is expected.
// When end of slice is reached reading starts from the beginning.
// The time.Duration value can be used to delay the Read function.
package inputDummy

import (
	"io"
	"time"
)

// Source delivers an endless circle of bytes.
type Source struct {
	b      []byte        // internal byte buffer to read from
	t      time.Duration // tinternal delay
	offset int           // next read position
}

// New creates an instance of simulator.
// b contains bytes it reads cyclically
func New(b []byte, t time.Duration) *Source {
	s := &Source{}
	s.b = b
	s.t = t
	return s
}

// Read is the exported method.
func (p *Source) Read(b []byte) (int, error) {
	var i int
	time.Sleep(p.t)
	if nil == p.b || 0 == len(p.b) {
		return 0, io.EOF
	}
	if 0 == len(b) {
		return 0, nil
	}
	for i = range b {
		b[i] = p.b[(p.offset+i)%len(p.b)]
	}
	i++ // adjust to byte count
	p.offset += i
	return i, nil
}
