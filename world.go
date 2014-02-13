package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/BenLubar/untitled-game/simplex"
	"github.com/steveyen/gkvlite"
	"log"
	"math/rand"
	"os"
	"sync"
)

const CurrentSaveVersion = 1

type World struct {
	chunks   map[ChunkCoord]*Chunk
	entities map[EntityReference]*Entity

	storeFile *os.File

	store  *gkvlite.Store
	global *gkvlite.Collection
	chunk  *gkvlite.Collection
	entity *gkvlite.Collection

	simplex *simplex.Simplex

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

var kTime = []byte("time")

func (w *World) Time() Timestamp {
	t, err := w.global.Get(kTime)
	if err != nil {
		panic(err)
	}
	if len(t) != 8 {
		return Timestamp(0)
	}
	return Timestamp(binary.BigEndian.Uint64(t))
}

var kSeed = []byte("seed")

func (w *World) Rand(f func(*rand.Rand)) (err error) {
	w.Lock()
	defer w.Unlock()

	return w.rand(f)
}

func (w *World) rand(f func(*rand.Rand)) (err error) {
	b, err := w.global.Get(kSeed)
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

	err = w.global.Set(kSeed, buf.Bytes())
	return
}

func (w *World) setSeed(seed *Seed) (err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&seed)
	if err != nil {
		return
	}

	// This is ok here, but nowhere else. This function is the only one that
	// can be called on an uninitialized world.
	if w.global == nil {
		w.global = w.store.SetCollection("global", nil)
	}

	err = w.global.Set(kSeed, buf.Bytes())
	return
}

var kSimplex = []byte("simplex")

func (w *World) Simplex() (s *simplex.Simplex, err error) {
	w.Lock()
	defer w.Unlock()

	return w.getSimplex()
}

func (w *World) getSimplex() (s *simplex.Simplex, err error) {
	if w.simplex != nil {
		return w.simplex, nil
	}

	b, err := w.global.Get(kSimplex)
	if err != nil {
		return
	}

	if len(b) == 0 {
		err = w.rand(func(r *rand.Rand) {
			s = simplex.New(r)
		})
		if err != nil {
			s = nil
			return
		}

		var buf bytes.Buffer
		err = gob.NewEncoder(&buf).Encode(&s)
		if err != nil {
			s = nil
			return
		}

		err = w.global.Set(kSimplex, buf.Bytes())
		w.simplex = s

		return
	}

	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&s)
	if err != nil {
		s = nil
	}
	w.simplex = s
	return
}

var kNextEntityID = []byte("entid")

func (w *World) NewEntity() (ent *Entity, err error) {
	w.Lock()
	defer w.Unlock()

	id_, err := w.global.Get(kNextEntityID)
	if err != nil {
		return
	}

	id := EntityReference(1)
	if len(id_) == 8 {
		id = EntityReference(binary.BigEndian.Uint64(id_) + 1)
	}

	id_ = make([]byte, 8)
	binary.BigEndian.PutUint64(id_, uint64(id))
	err = w.global.Set(kNextEntityID, id_)
	if err != nil {
		return
	}

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

	v, err := w.entity.Get(id.bytes())
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
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(ent); err != nil {
			panic(err)
		}
		if err := w.entity.Set(ent.ID.bytes(), buf.Bytes()); err != nil {
			panic(err)
		}
		delete(w.entities, ent.ID)
	}
}

func (w *World) RequestChunk(coord ChunkCoord) (c *Chunk, err error) {
	w.Lock()
	defer w.Unlock()

	if c = w.chunks[coord]; c != nil {
		c.references++
		return c, nil
	}

	v, err := w.chunk.Get(coord.bytes())
	if err != nil {
		log.Printf("error reading chunk (%d, %d): %v", coord.X, coord.Y, err)
		return
	}

	if len(v) == 0 {
		c, err = w.generateChunk(coord)
		if err != nil {
			log.Printf("error generating chunk (%d, %d): %v", coord.X, coord.Y, err)
			c = nil
			return
		}

		c.references++
		if w.chunks == nil {
			w.chunks = make(map[ChunkCoord]*Chunk)
		}
		w.chunks[coord] = c
		return
	}

	err = gob.NewDecoder(bytes.NewReader(v)).Decode(&c)
	if err != nil {
		log.Printf("error decoding chunk (%d, %d): %v", coord.X, coord.Y, err)
		c = nil
		return
	}
	c.references++
	if w.chunks == nil {
		w.chunks = make(map[ChunkCoord]*Chunk)
	}
	w.chunks[coord] = c
	return
}

func (w *World) ReleaseChunk(c *Chunk) {
	w.Lock()
	defer w.Unlock()

	if c == nil {
		panic("release of nil chunk")
	}

	if w.chunks[c.ChunkCoord] != c {
		panic(fmt.Sprintf("Chunk %v released, but a different object was in the cache:\n\n%v\n\n%v", c.ChunkCoord, c, w.chunks[c.ChunkCoord]))
	}

	if c.references == 0 {
		panic(fmt.Sprintf("release of unreferenced chunk:\n\n%v", c))
	}
	c.references--
	if c.references == 0 {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(c); err != nil {
			panic(err)
		}
		if err := w.chunk.Set(c.ChunkCoord.bytes(), buf.Bytes()); err != nil {
			panic(err)
		}
		delete(w.chunks, c.ChunkCoord)
	}
}

var kVersion = []byte("version")

func (w *World) init() (err error) {
	w.global = w.store.SetCollection("global", nil)
	w.chunk = w.store.SetCollection("chunk", nil)
	w.entity = w.store.SetCollection("entity", nil)

	versionBuf, err := w.global.Get(kVersion)
	if err != nil {
		log.Printf("error getting version: %v", err)
		return
	}

	{
		tmp := make([]byte, 8)
		copy(tmp, versionBuf)
		versionBuf = tmp
	}

	version := binary.BigEndian.Uint64(versionBuf)

	switch version {
	case 0:
		for cx := int64(-16); cx <= 16; cx++ {
			for cy := int64(-16); cy <= 16; cy++ {
				c, err := w.RequestChunk(ChunkCoord{cx, cy})
				if err != nil {
					log.Printf("error getting chunk (%d, %d): %v", cx, cy, err)
					return err
				}
				w.ReleaseChunk(c)
			}
		}

		binary.BigEndian.PutUint64(versionBuf, CurrentSaveVersion)
		err = w.global.Set(kVersion, versionBuf)
		if err != nil {
			return err
		}

	case CurrentSaveVersion:
		// no updates

	default:
		return fmt.Errorf("unexpected version: %d", version)
	}

	for i := int64(-1); i <= int64(1); i++ {
		for j := int64(-1); j <= int64(1); j++ {
			w.RequestChunk(ChunkCoord{i, j})
		}
	}

	return w.store.Flush()
}

func (w *World) Tick() {
	// TODO: game ticks
}
