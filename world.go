package main

import (
	"fmt"
	"sync"
)

const CurrentSaveVersion = 1

type ErrSaveVersionTooNew uint64

func (err ErrSaveVersionTooNew) Error() string {
	return fmt.Sprintf("save version %d is greater than the highest supported version %d", uint64(err), CurrentSaveVersion)
}

type ErrSaveCorrupt string

func (err ErrSaveCorrupt) Error() string {
	return fmt.Sprintf("corrupt save: %s", string(err))
}

type World struct {
	Version uint64

	Seed *Seed

	// invariant: Entities[i].ID == EntityReference(i+1)
	Entities []*Entity

	Time Timestamp

	mtx sync.RWMutex
}

var (
	world     *World
	worldLock sync.Mutex
)

func GetWorld() *World {
	worldLock.Lock()
	defer worldLock.Unlock()

	return world
}

func (w *World) Do(f func()) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	f()
}

func (w *World) RDo(f func()) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	f()
}

func (w *World) AfterLoad() (err error) {
	w.Do(func() {
		if w.Version > CurrentSaveVersion {
			err = ErrSaveVersionTooNew(w.Version)
			return
		}

		if w.Seed == nil {
			err = ErrSaveCorrupt("seed is missing")
			return
		}

		switch w.Version {
		case 0:
			// TODO: generate a world
			fallthrough

		case CurrentSaveVersion:
			w.Version = CurrentSaveVersion

		default:
			panic(fmt.Sprintf("unexpected save version: %d", w.Version))
		}

		if w.Time == 0 {
			err = ErrSaveCorrupt("time is missing")
			return
		}

		for i, e := range w.Entities {
			if e.ID != EntityReference(i+1) {
				err = ErrSaveCorrupt(fmt.Sprintf("entity %d failed invariant test", i))
				return
			}
		}
	})
	return
}

func (w *World) BeforeSave() (err error) {
	w.Do(func() {
	})
	return
}
