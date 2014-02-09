package main

import (
	"strconv"
	"sync"
)

type Entity struct {
	ID         EntityReference
	Components []Component
	mtx        sync.RWMutex
}

func (w *World) NewEntity() (ent *Entity) {
	w.Do(func() {
		ent = &Entity{ID: EntityReference(len(w.Entities) + 1)}
		w.Entities = append(w.Entities, ent)
	})
	return
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

type EntityReference uint64

func (ref EntityReference) Get(w *World) (ent *Entity) {
	if ref != 0 {
		w.RDo(func() {
			ent = w.Entities[ref-1]
		})
	}
	return
}

func (w *World) EachEntity(f func(*Entity)) {
	var entities []*Entity

	w.RDo(func() {
		// This is safe because entities are never removed, so the only thing we have
		// to worry about is grabbing the entity list in between an update to the length
		// and to the pointer.
		entities = w.Entities
	})

	for _, e := range entities {
		f(e)
	}
}
