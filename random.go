package main

import (
	"crypto/sha512"
	"encoding/binary"
	"strconv"
	"unsafe"
)

// Seed is an implementation of math/rand.Source that uses SHA512 with an arbitrary-length
// "seed". It can be saved and re-loaded (using gob, json, etc.) If used from multiple
// goroutines concurrently, external locking will be needed.
type Seed struct {
	// Fields exported for encoding/*
	Text string
	Buf  [sha512.Size]byte
	Ptr  uintptr
}

func NewSeed(text string) *Seed {
	var s Seed
	s.SeedText(text)
	return &s
}

const int64Size = unsafe.Sizeof(int64(0))

func (s *Seed) Int63() int64 {
	if s.Ptr > sha512.Size-int64Size {
		s.Ptr = 0
		s.Buf = sha512.Sum512(append([]byte(s.Text), s.Buf[:]...))
	}
	n := binary.LittleEndian.Uint64(s.Buf[s.Ptr : s.Ptr+int64Size])
	s.Ptr += int64Size
	n &^= 1 << 63 // make sure the highest bit is not set
	return int64(n)
}

// Seed converts the given int64 to its decimal representation and calls SeedText. Seriously,
// just use SeedText.
func (s *Seed) Seed(seed int64) {
	s.SeedText(strconv.FormatInt(seed, 10))
}

// SeedText resets the Seed to the return value of NewSeed(seed).
func (s *Seed) SeedText(seed string) {
	s.Text = seed
	s.Buf = sha512.Sum512([]byte(seed))
	s.Ptr = 0
}
