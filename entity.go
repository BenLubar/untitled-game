package main

import (
	"strconv"
	"sync"
)

type Entity struct {
	ID         EntityReference
	Components []Component
	sync.RWMutex
}

// invariant: entities.list[i].ID == EntityReference(i+1)
var entities struct {
	list []*Entity
	sync.RWMutex
}

func NewEntity() *Entity {
	entities.Lock()
	defer entities.Unlock()

	ent := &Entity{ID: EntityReference(len(entities.list) + 1)}
	entities.list = append(entities.list, ent)
	return ent
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
	e.Lock()
	defer e.Unlock()

	f()
}

func (e *Entity) RDo(f func()) {
	e.RLock()
	defer e.RUnlock()

	f()
}

type EntityReference uint64

func (ref EntityReference) Get() *Entity {
	if ref == 0 {
		return nil
	}

	entities.RLock()
	defer entities.RUnlock()

	return entities.list[ref-1]
}

func EachEntity(f func(*Entity)) {
	entities.RLock()
	defer entities.RUnlock()

	for _, e := range entities.list {
		func() {
			// unlock the list so new entities can be created during iteration
			entities.RUnlock()
			defer entities.RLock()

			f(e)
		}()
	}
}
