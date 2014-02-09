package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	repaint := time.Tick(time.Second / 60)

	events := make(chan termbox.Event)
	go pollEvents(events)
	for {
		select {
		case e := <-events:
			switch e.Type {
			case termbox.EventError:
				panic(e.Err)
			case termbox.EventKey:
				if world := GetWorld(); world == nil {
					if !mainMenu.inputKey(e.Key, e.Ch, e.Mod) {
						return
					}
				} else {
					// TODO: game UI
					panic(fmt.Sprintf("%v, %v, %v", e.Key, e.Ch, e.Mod))
				}
			case termbox.EventMouse:
				if world := GetWorld(); world == nil {
					mainMenu.inputMouse(e.MouseX, e.MouseY)
				} else {
					// TODO: game UI
					panic(fmt.Sprintf("%v, %v", e.MouseX, e.MouseY))
				}
			case termbox.EventResize:
				// ignore
			}

		case <-repaint:
			termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)

			w, h := termbox.Size()

			if world := GetWorld(); world == nil {
				mainMenu.render(w, h)
			} else {
				world.Tick()
				// TODO: game UI
				renderBorder(w, h, world)
			}
			termbox.Flush()
		}
	}
}

func renderBorder(w, h int, world *World) {
	var t Timestamp
	world.RDo(func() {
		t = world.Time
	})

	x := 0

	x++
	termbox.SetCell(w-x, 0, '╗', termbox.ColorBlack, termbox.ColorWhite)
	x++
	termbox.SetCell(w-x, 0, '╞', termbox.ColorBlack, termbox.ColorWhite)

	divider := func() {
		x++
		termbox.SetCell(w-x, 0, '╡', termbox.ColorBlack, termbox.ColorWhite)
		x++
		termbox.SetCell(w-x, 0, '═', termbox.ColorBlack, termbox.ColorWhite)
		x++
		termbox.SetCell(w-x, 0, '╞', termbox.ColorBlack, termbox.ColorWhite)
	}

	// always use at least four digits for the year
	for year := t.Year(); year != 0 || x < 4; year /= 10 {
		x++
		termbox.SetCell(w-x, 0, '0'+rune(year%10), termbox.ColorBlack, termbox.ColorWhite)
	}
	divider()
	season := t.Season().String()
	for i, ch := range season {
		termbox.SetCell(w-x-len(season)+i, 0, ch, termbox.ColorBlack, termbox.ColorWhite)
	}
	x += len(season)
	divider()
	tod := t.TimeOfDay().String()
	for i, ch := range tod {
		termbox.SetCell(w-x-len(tod)+i, 0, ch, termbox.ColorBlack, termbox.ColorWhite)
	}
	x += len(tod)

	x++
	termbox.SetCell(w-x, 0, '╡', termbox.ColorBlack, termbox.ColorWhite)
	for x < w-1 {
		x++
		termbox.SetCell(w-x, 0, '═', termbox.ColorBlack, termbox.ColorWhite)
	}
	termbox.SetCell(0, 0, '╔', termbox.ColorBlack, termbox.ColorWhite)
	for y := 1; y < h-1; y++ {
		termbox.SetCell(0, y, '║', termbox.ColorBlack, termbox.ColorWhite)
		termbox.SetCell(w-1, y, '║', termbox.ColorBlack, termbox.ColorWhite)
	}
	termbox.SetCell(0, h-1, '╚', termbox.ColorBlack, termbox.ColorWhite)
	for x = 1; x < w-1; x++ {
		termbox.SetCell(x, h-1, '═', termbox.ColorBlack, termbox.ColorWhite)
	}
	termbox.SetCell(w-1, h-1, '╝', termbox.ColorBlack, termbox.ColorWhite)
}

func pollEvents(ch chan<- termbox.Event) {
	for {
		ch <- termbox.PollEvent()
	}
}
