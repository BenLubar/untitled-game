package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

func main() {
	for t := Timestamp(1); t <= ts_ticks_per_year + 1; t += ts_ticks_per_day / 5 {
		fmt.Printf("It is %v in %v, the %v day of %v.\n", t.TimeOfDay(), t.Season(), humanize.Ordinal(int(t.Day())), t.Year())
	}
}
