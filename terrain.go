package main

type Tile struct {
	Type     TileType
	Entities []EntityReference
}

type TileType uint8

const (
	TileRock TileType = iota
	TileSand
	TileDirt
	TileGrass
	TileWater
)
