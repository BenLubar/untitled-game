package main

import (
	"fmt"
)

func main() {
	e1 := NewEntity()
	e2 := NewEntity()
	e3 := NewEntity()
	e1.Do(func() {
		e1.Components = append(e1.Components, &CreatedComponent{ID: e2.ID, Location: e3.ID, Time: 1, Material: nil})
		e1.Components = append(e1.Components, &OwnerOfComponent{ID: e2.ID})
	})
	e2.Do(func() {
		e2.Components = append(e2.Components, &CreatedByComponent{ID: e1.ID})
		e2.Components = append(e2.Components, &OwnerComponent{ID: e1.ID, Start: 1, End: 0})
	})
	EachEntity(func(e *Entity) {
		fmt.Println(e)
		fmt.Println()
	})
}
