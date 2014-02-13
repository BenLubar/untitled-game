package main

import (
	"encoding/binary"
	"log"
)

const ChunkSize = 1 << 8

type ChunkCoord struct {
	X, Y int64
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
		groundY := s.Noise2(fx, 0) * 16 / ChunkSize
		for y := range c.Tiles[x] {
			fy := float64(coord.Y) + float64(y)/float64(ChunkSize)
			if fy < groundY {
				c.Tiles[x][y].Type = TileDirt
			} else {
				c.Tiles[x][y].Type = TileAir
			}
		}
	}
	return
}
