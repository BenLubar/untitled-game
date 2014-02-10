package main

import (
	"fmt"
)

type LocationComponent struct {
	ID   EntityReference
	X, Y int64
}

func init() {
	registerComponentType(&LocationComponent{})
}

func (c *LocationComponent) String() string {
	return fmt.Sprintf("LOCATION id[entity]=%v x[int64]=%v y[int64]=%v", c.ID, c.X, c.Y)
}
