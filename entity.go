package main

import (
	"encoding/binary"
	"strconv"
	"sync"
)

type EntityReference uint64

func (id EntityReference) bytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[:], uint64(id))
	return b
}

type Entity struct {
	ID         EntityReference
	Components []Component
	references uint
	mtx        sync.RWMutex
}

func (e *Entity) String() string {
	var buf []byte
	buf = append(buf, "ENTITY id[entity]="...)
	buf = strconv.AppendUint(buf, uint64(e.ID), 10)

	e.RDo(func() {
		for _, c := range e.Components {
			buf = append(buf, "\n\t"...)
			buf = append(buf, c.String()...)
		}
	})

	return string(buf)
}

func (e *Entity) Do(f func()) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	f()
}

func (e *Entity) RDo(f func()) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	f()
}
