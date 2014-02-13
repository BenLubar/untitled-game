package main

import (
	"fmt"
)

type LocationComponent struct {
	ID     EntityReference
	ChunkX int64
	ChunkY int64
	ChunkZ int64
	TileX  uint8
	TileY  uint8
	TileZ  uint8
}

func init() {
	registerComponentType(&LocationComponent{})
}

func (c *LocationComponent) String() string {
	return fmt.Sprintf("LOCATION id[entity]=%v chunk[ints]=(%v,%v,%v) tile[ints]=(%v,%v,%v)", c.ID, c.ChunkX, c.ChunkY, c.ChunkZ, c.TileX, c.TileY, c.TileZ)
}
