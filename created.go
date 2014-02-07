package main

import (
	"fmt"
)

type CreatedByComponent struct {
	ID EntityReference
}

func init() {
	registerComponentType(&CreatedByComponent{})
}

func (c *CreatedByComponent) String() string {
	return fmt.Sprintf("CREATED_BY id[entity]=%v", c.ID)
}

type CreatedComponent struct {
	ID       EntityReference
	Location EntityReference
	Time     Timestamp
	Material []EntityReference
}

func init() {
	registerComponentType(&CreatedComponent{})
}

func (c *CreatedComponent) String() string {
	return fmt.Sprintf("CREATED id[entity]=%v location[entity]=%v time[time]=%v material[entities]=%v", c.ID, c.Location, c.Time, c.Material)
}
