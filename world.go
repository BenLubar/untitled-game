package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
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

const ChunkTiles = 128

type World struct {
	Version uint64

	Seed *Seed

	// invariant: Entities[i].ID == EntityReference(i+1)
	Entities []*Entity

	Time Timestamp

	Strings    []string
	revStrings map[string]int

	Overworld map[[2]int64]*[ChunkTiles * ChunkTiles]Tile
	Sites     map[EntityReference]*[ChunkTiles * ChunkTiles]Tile

	saveName string

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

func (w *World) Tick() {
	w.Do(func() {
		w.Time++
		// TODO: actual update code
	})
}

func (w *World) AfterLoad() (err error) {
	w.Do(func() {
		var (
			generatingTerrain    = []rune("Generating Terrain")
			generatingPreHistory = []rune("Generating Pre-History")
			generatingHistory    = []rune("Generating History")
		)

		if w.Version > CurrentSaveVersion {
			err = ErrSaveVersionTooNew(w.Version)
			return
		}

		if w.Seed == nil {
			err = ErrSaveCorrupt("seed is missing")
			return
		}

		w.revStrings = make(map[string]int, len(w.Strings))
		for i, s := range w.Strings {
			w.revStrings[s] = i
		}

		switch w.Version {
		case 0:
			w.Overworld = make(map[[2]int64]*[ChunkTiles * ChunkTiles]Tile)
			w.Sites = make(map[EntityReference]*[ChunkTiles * ChunkTiles]Tile)

			const worldChunks = 32
			for x := int64(-worldChunks); x <= worldChunks; x++ {
				for y := int64(-worldChunks); y <= worldChunks; y++ {
					var chunk [ChunkTiles * ChunkTiles]Tile
					for i := range chunk {
						x2 := x*ChunkTiles + int64(i)/ChunkTiles
						y2 := y*ChunkTiles + int64(i)%ChunkTiles

						if x2 == y2 || x2 == -y2 {
							chunk[i].Type = TileDirt
						} else {
							chunk[i].Type = TileGrass
						}
						// TODO: actual terrain
					}
					w.Overworld[[2]int64{x, y}] = &chunk

					termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
					w, h := termbox.Size()

					completion := int(((x+worldChunks)*(worldChunks*2+1) + y + worldChunks) * int64(w) / (worldChunks*2 + 1) / (worldChunks*2 + 1))
					for j := 0; j < w; j++ {
						ch := rune(' ')
						if j > 0 && j <= len(generatingTerrain) {
							ch = rune(generatingTerrain[j-1])
						}
						if j < completion {
							termbox.SetCell(j, h/2, ch, termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
						} else {
							termbox.SetCell(j, h/2, ch, termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
						}
					}
					termbox.Flush()
				}
			}

			for i := Timestamp(0); i < ts_days_per_year*100; i++ {
				for j := Timestamp(0); j < ts_ticks_per_day; j++ {
					// TODO
				}

				termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
				w, h := termbox.Size()

				for j, ch := range generatingTerrain {
					termbox.SetCell(j+1, h/2-1, ch, termbox.ColorWhite, termbox.ColorBlack)
				}
				completion := int(i * Timestamp(w) / (ts_days_per_year * 100))
				for j := 0; j < w; j++ {
					ch := rune(' ')
					if j > 0 && j <= len(generatingPreHistory) {
						ch = rune(generatingPreHistory[j-1])
					}
					if j < completion {
						termbox.SetCell(j, h/2, ch, termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
					} else {
						termbox.SetCell(j, h/2, ch, termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
					}
				}

				termbox.Flush()
			}

			w.Time = 0
			for i := Timestamp(0); i < ts_days_per_year*100; i++ {
				for j := Timestamp(0); j < ts_ticks_per_day; j++ {
					w.Time++
					// TODO
				}

				termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
				w, h := termbox.Size()

				for j, ch := range generatingTerrain {
					termbox.SetCell(j+1, h/2-2, ch, termbox.ColorWhite, termbox.ColorBlack)
				}
				for j, ch := range generatingPreHistory {
					termbox.SetCell(j+1, h/2-1, ch, termbox.ColorWhite, termbox.ColorBlack)
				}
				completion := int(i * Timestamp(w) / (ts_days_per_year * 100))
				for j := 0; j < w; j++ {
					ch := rune(' ')
					if j > 0 && j <= len(generatingHistory) {
						ch = rune(generatingHistory[j-1])
					}
					if j < completion {
						termbox.SetCell(j, h/2, ch, termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
					} else {
						termbox.SetCell(j, h/2, ch, termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
					}
				}

				termbox.Flush()
			}

			err = fmt.Errorf("TODO: world generator (sorry)")
			return

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

func (w *World) StringID(s string) (i int) {
	w.Do(func() {
		var ok bool
		if i, ok = w.revStrings[s]; ok {
			return
		}
		i = len(w.Strings)
		w.Strings = append(w.Strings, s)
		w.revStrings[s] = i
	})
	return
}

func (w *World) StringForID(i int) (s string, ok bool) {
	w.RDo(func() {
		if i < 0 || i >= len(w.Strings) {
			return
		}
		s, ok = w.Strings[i], true
	})
	return
}

func (w *World) Save() (err error) {
	w.Do(func() {
		err = w.save(os.O_WRONLY | os.O_CREATE | os.O_TRUNC)
	})
	return
}

func (w *World) save(flag int) error {
	_ = os.Mkdir(SaveDirName, 0777)
	f, err := os.OpenFile(w.saveName, flag, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer g.Close()

	return gob.NewEncoder(g).Encode(w)
}
