package main

import (
	"encoding/binary"
	"log"
)

const chunkShift = 8
const ChunkSize = 1 << chunkShift

type ChunkCoord struct {
	X, Y int64
}

func ChunkForTile(x, y int64) ChunkCoord {
	return ChunkCoord{x >> chunkShift, y >> chunkShift}
}

func (coord ChunkCoord) bytes() []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[0:8], uint64(coord.X))
	binary.BigEndian.PutUint64(b[8:16], uint64(coord.Y))
	return b
}

type Chunk struct {
	ChunkCoord
	Tiles      [ChunkSize][ChunkSize]Tile
	references uint
}

type Tile struct {
	Type TileType
}

type TileType uint8

const (
	TileAir TileType = iota
	TileRock
	TileSand
	TileDirt
	TileGrass
	TileWater
)

func (w *World) generateChunk(coord ChunkCoord) (c *Chunk, err error) {
	s, err := w.getSimplex()
	if err != nil {
		log.Printf("error getting simplex: %v", err)
		return
	}

	c = &Chunk{ChunkCoord: coord}
	// TODO: more interesting worldgen than "flat ground with lumps"
	for x := range c.Tiles {
		fx := float64(coord.X) + float64(x)/float64(ChunkSize)
		groundY := s.Noise2(fx, 0.)*16./ChunkSize + 4./ChunkSize
		groundY += s.Noise2(fx*4., 1.) * 2. / ChunkSize
		rockY := s.Noise2(fx, 2.)*16./ChunkSize - 6./ChunkSize
		grassY := groundY - 1./ChunkSize
		waterY := 0. / ChunkSize
		if groundY < waterY-1./ChunkSize {
			grassY = groundY + 10000. // no grass underwater.
		}
		for y := range c.Tiles[x] {
			fy := float64(coord.Y) + float64(y)/float64(ChunkSize)
			if fy < rockY {
				c.Tiles[x][y].Type = TileRock
			} else if fy < groundY {
				if fy >= grassY {
					c.Tiles[x][y].Type = TileGrass
				} else {
					c.Tiles[x][y].Type = TileDirt
				}
			} else if fy < waterY {
				c.Tiles[x][y].Type = TileWater
			} else {
				c.Tiles[x][y].Type = TileAir
			}
		}
	}
	return
}
