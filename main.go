package main

import (
	"fmt"
	"github.com/davecheney/profile"
	"github.com/nsf/termbox-go"
	"time"
)

func main() {
	defer profile.Start(&profile.Config{
		Quiet:       true,
		CPUProfile:  true,
		MemProfile:  true,
		ProfilePath: "./prof/",
	}).Stop()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	defer func() {
		if world := GetWorld(); world != nil {
			if err := world.store.Flush(); err != nil {
				panic(err)
			}
		}
	}()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	repaint := time.Tick(time.Second / 60)

	var playerX, playerY int64

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
					if e.Key == termbox.KeyArrowDown {
						playerY--
						break
					}
					if e.Key == termbox.KeyArrowUp {
						playerY++
						break
					}
					if e.Key == termbox.KeyArrowLeft {
						playerX--
						break
					}
					if e.Key == termbox.KeyArrowRight {
						playerX++
						break
					}
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
				renderWorld(playerX, playerY, w, h, world)
				// TODO: game UI
				renderBorder(w, h, world)
			}
			termbox.Flush()
		}
	}
}

func renderWorld(playerX, playerY int64, w, h int, world *World) {
	for cx := int64(-1); cx <= 1; cx++ {
		for cy := int64(-1); cy <= 1; cy++ {
			renderChunk(ChunkCoord{cx, cy}, int(cx*ChunkSize-playerX)+w/2, int(playerY-cy*ChunkSize)+h/2, world)
		}
	}
}

func renderChunk(coord ChunkCoord, startX, startY int, world *World) {
	c, err := world.RequestChunk(coord)
	if err != nil {
		panic(err)
	}
	defer world.ReleaseChunk(c)

	for x := range c.Tiles {
		for y := range c.Tiles[x] {
			var color termbox.Attribute
			var text []rune
			switch c.Tiles[x][y].Type {
			case TileAir:
				color = termbox.ColorBlack
				text = []rune(" air ")
			case TileDirt:
				color = termbox.ColorYellow
				text = []rune(" dirt ")
			case TileGrass:
				color = termbox.ColorGreen
				text = []rune(" grass ")
			case TileRock:
				color = termbox.ColorRed
				text = []rune(" rock ")
			case TileSand:
				color = termbox.ColorCyan
				text = []rune(" sand ")
			case TileWater:
				color = termbox.ColorBlue
				text = []rune(" water ")
			}
			termbox.SetCell(startX+x, startY-y, text[((coord.X*ChunkSize+int64(x)+coord.Y*ChunkSize+int64(y))%int64(len(text))+int64(len(text)))%int64(len(text))], termbox.AttrBold|color, color)
		}
	}
}

func renderBorder(w, h int, world *World) {
	t := world.Time()

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
	for year := t.Year(); year != 0 || x < 4+2; year /= 10 {
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
