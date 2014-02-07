package main

import (
	"fmt"
)

type OwnerInverseComponent struct {
	ID EntityReference
}

func init() {
	registerComponentType(&OwnerInverseComponent{})
}

func (c *OwnerInverseComponent) String() string {
	return fmt.Sprintf("OWNER_OF id[entity]=%v", c.ID)
}

type OwnerComponent struct {
	ID    EntityReference
	Start Timestamp
	End   Timestamp
}

func init() {
	registerComponentType(&OwnerComponent{})
}

func (c *OwnerComponent) String() string {
	return fmt.Sprintf("OWNER id[entity]=%v start[time]=%v end[time]=%v", c.ID, c.Start, c.End)
}
