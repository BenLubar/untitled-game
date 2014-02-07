package main

import (
	"encoding/gob"
	"fmt"
)

type Component interface {
	fmt.Stringer
}

func registerComponentType(v Component) {
	gob.Register(v)
}
