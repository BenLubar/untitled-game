package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/cznic/kv"
	"math/rand"
	"sync"
)

const CurrentSaveVersion = 1

type World struct {
	entities map[EntityReference]*Entity
	db       *kv.DB
	sync.Mutex
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

var kTime = []byte("ctime")

func (w *World) Time() Timestamp {
	t, err := w.db.Inc(kTime, 0)
	if err != nil {
		panic(err)
	}
	return Timestamp(t)
}

var kSeed = []byte("cseed")

func (w *World) Rand(f func(*rand.Rand)) (err error) {
	w.Lock()
	defer w.Unlock()

	if err = w.db.BeginTransaction(); err != nil {
		return
	}
	defer func() {
		if err_ := w.db.Commit(); err == nil {
			err = err_
		}
	}()

	b, err := w.db.Get(nil, kSeed)
	if err != nil {
		return
	}

	var seed Seed
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&seed)
	if err != nil {
		return
	}

	f(rand.New(&seed))

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&seed)
	if err != nil {
		return
	}

	err = w.db.Set(kSeed, buf.Bytes())
	return
}

func (w *World) setSeed(seed *Seed) (err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&seed)
	if err != nil {
		return
	}

	err = w.db.Set(kSeed, buf.Bytes())
	return
}

var kNextEntityID = []byte("centid")

func (w *World) NewEntity() (ent *Entity, err error) {
	w.Lock()
	defer w.Unlock()

	id_, err := w.db.Inc(kNextEntityID, 1)
	if err != nil {
		return
	}
	id := EntityReference(id_)
	ent = &Entity{ID: id}
	ent.references++
	if w.entities == nil {
		w.entities = make(map[EntityReference]*Entity)
	}
	w.entities[id] = ent
	return
}

func (w *World) RequestEntity(id EntityReference) (ent *Entity, err error) {
	w.Lock()
	defer w.Unlock()

	if ent = w.entities[id]; ent != nil {
		ent.references++
		return ent, nil
	}

	v, err := w.db.Get(nil, id.bytes())
	if err != nil {
		return
	}

	err = gob.NewDecoder(bytes.NewReader(v)).Decode(&ent)
	if err != nil {
		ent = nil
	} else {
		ent.references++
		if w.entities == nil {
			w.entities = make(map[EntityReference]*Entity)
		}
		w.entities[id] = ent
	}
	return
}

func (w *World) ReleaseEntity(ent *Entity) {
	w.Lock()
	defer w.Unlock()

	if ent == nil {
		panic("release of nil entity")
	}

	if w.entities[ent.ID] != ent {
		panic(fmt.Sprintf("Entity %d released, but a different object was in the cache:\n\n%v\n\n%v", ent.ID, ent, w.entities[ent.ID]))
	}

	if ent.references == 0 {
		panic(fmt.Sprintf("release of unreferenced entity:\n\n%v", ent))
	}
	ent.references--
	if ent.references == 0 {
		if err := w.db.BeginTransaction(); err != nil {
			panic(err)
		}
		defer func() {
			if err := w.db.Commit(); err != nil {
				panic(err)
			}
		}()
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(ent); err != nil {
			panic(err)
		}
		if err := w.db.Set(ent.ID.bytes(), buf.Bytes()); err != nil {
			panic(err)
		}
		delete(w.entities, ent.ID)
	}
}

var kVersion = []byte("cver")

func (w *World) init() (err error) {
	if err = w.db.BeginTransaction(); err != nil {
		return
	}
	defer func() {
		if err_ := w.db.Commit(); err == nil {
			err = err_
		}
	}()

	version, err := w.db.Inc(kVersion, 0)
	if err != nil {
		return
	}

	return fmt.Errorf("TODO: world generation (%d)", version)
}

func (w *World) Tick() {
	panic("TODO: game ticks")
}
