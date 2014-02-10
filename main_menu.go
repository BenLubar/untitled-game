package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const SaveDirName = "saves_5CC9DB70-EEC5-47EA-94B6-398BFC12E4A7"

var mainMenu mainMenuUI

func init() {
	mainMenu.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 16; i++ {
		const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		mainMenu.seed = append(mainMenu.seed, rune(characters[mainMenu.rand.Intn(len(characters))]))
	}

	saves, err := os.Open(SaveDirName)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	names, err := saves.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		if strings.HasSuffix(name, ".sav") {
			mainMenu.saveNames = append(mainMenu.saveNames, name[:len(name)-len(".sav")])
		}
	}
}

const (
	menuStateMain = iota
	menuStateError
	menuStateNew
)

type mainMenuUI struct {
	rand        *rand.Rand
	saveNames   []string
	choiceIndex int
	state       uint
	saveName    []rune
	seed        []rune
	err         string
}

func (m *mainMenuUI) render(w, h int) {
	m.drawTextFlicker(w, 2, "5CC9DB70-EEC5-47EA-94B6-398BFC12E4A7", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)

	switch m.state {
	case menuStateMain:
		skip := 0

		if m.choiceIndex > h/2-5 {
			skip = m.choiceIndex - (h/2 - 5)
		}

		for i, name := range m.saveNames[skip:] {
			if m.choiceIndex == i+skip {
				m.drawText(w, 5+i, fmt.Sprintf("Load %q", name), termbox.ColorBlack, termbox.ColorWhite)
			} else {
				m.drawText(w, 5+i, fmt.Sprintf("Load %q", name), termbox.ColorWhite, termbox.ColorBlack)
			}
		}
		if m.choiceIndex == len(m.saveNames) {
			m.drawText(w, len(m.saveNames)+5-skip, "New Game", termbox.ColorBlack, termbox.ColorWhite)
		} else {
			m.drawText(w, len(m.saveNames)+5-skip, "New Game", termbox.ColorWhite, termbox.ColorBlack)
		}

	case menuStateError:
		m.drawText(w, 5, m.err, termbox.ColorRed, termbox.ColorBlack)

	case menuStateNew:
		m.drawText(w, 5, "Save Name", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
		if m.choiceIndex == 0 {
			m.drawText(w, 6, string(m.saveName)+"_", termbox.ColorBlack, termbox.ColorWhite)
		} else {
			m.drawText(w, 6, string(m.saveName), termbox.ColorWhite, termbox.ColorBlack)
		}

		m.drawText(w, 8, "Seed", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
		if m.choiceIndex == 1 {
			m.drawText(w, 9, string(m.seed)+"_", termbox.ColorBlack, termbox.ColorWhite)
		} else {
			m.drawText(w, 9, string(m.seed), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func (m *mainMenuUI) drawText(w, y int, s string, fg, bg termbox.Attribute) {
	r := []rune(s)

	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', fg, bg)
	}

	for i, ch := range r {
		termbox.SetCell((w-len(r))/2+i, y, ch, fg, bg)
	}
}

func (m *mainMenuUI) drawTextFlicker(w, y int, s string, fg, bg termbox.Attribute) {
	r := []rune(s)

	flicker := m.rand.Intn(len(r))
	if (r[flicker] < '0' || r[flicker] > '9') && (r[flicker] < 'A' || r[flicker] > 'F') {
		flicker = -1
	}

	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', fg, bg)
	}

	for i, ch := range r {
		if i == flicker {
			const hex = "0123456789ABCDEF"
			termbox.SetCell((w-len(r))/2+i, y, rune(hex[m.rand.Intn(len(hex))]), fg^termbox.AttrBold, bg)
		} else {
			termbox.SetCell((w-len(r))/2+i, y, ch, fg, bg)
		}
	}
}

func (m *mainMenuUI) inputKey(key termbox.Key, ch rune, mod termbox.Modifier) bool {
	switch m.state {
	case menuStateMain:
		switch {
		case key == termbox.KeyArrowDown:
			m.choiceIndex = (m.choiceIndex + 1) % (len(m.saveNames) + 1)
		case key == termbox.KeyArrowUp:
			m.choiceIndex = (m.choiceIndex + (len(m.saveNames) + 1 - 1)) % (len(m.saveNames) + 1)
		case key == termbox.KeyArrowLeft || key == termbox.KeyArrowRight:
			// silently ignore
		case key == termbox.KeyEnter:
			if m.choiceIndex < len(m.saveNames) {
				m.loadGame(m.saveNames[m.choiceIndex])
			} else {
				switch m.choiceIndex - len(m.saveNames) {
				case 0:
					m.newGame()
				}
			}
		case key == termbox.KeyEsc:
			return false
		default:
			panic(fmt.Sprintf("%v, %v, %v", key, ch, mod))
		}

	case menuStateError:
		switch {
		case key == termbox.KeyEsc:
			m.err = ""
			m.state = menuStateMain
		default:
			panic(fmt.Sprintf("%v, %v, %v", key, ch, mod))
		}

	case menuStateNew:
		const fieldCount = 2

		switch {
		case key == termbox.KeyEsc:
			m.state = menuStateMain
			m.choiceIndex = len(m.saveNames)
		case key == termbox.KeyEnter:
			m.choiceIndex++
			if m.choiceIndex >= fieldCount {
				m.choiceIndex = fieldCount - 1
				if len(m.saveName) == 0 {
					m.choiceIndex = 0
					fmt.Print("\a")
					return true
				}
				w := &World{
					Seed:     NewSeed(string(m.seed)),
					saveName: filepath.Join(SaveDirName, string(m.saveName)+".sav"),
				}
				err := w.save(os.O_CREATE | os.O_EXCL | os.O_WRONLY)
				if err == nil {
					err = w.AfterLoad()
				}
				if err == nil {
					worldLock.Lock()
					world = w
					worldLock.Unlock()
				} else {
					m.err = err.Error()
					m.state = menuStateError
				}
			}
		case key == termbox.KeyArrowDown:
			m.choiceIndex = (m.choiceIndex + 1) % fieldCount
		case key == termbox.KeyArrowUp:
			m.choiceIndex = (m.choiceIndex + (fieldCount - 1)) % fieldCount
		case key == termbox.KeyArrowLeft || key == termbox.KeyArrowRight:
			fmt.Print("\a")
		case m.choiceIndex == 0 && (key == termbox.KeyBackspace || key == termbox.KeyBackspace2):
			if len(m.saveName) == 0 {
				fmt.Print("\a")
			} else {
				m.saveName = m.saveName[:len(m.saveName)-1]
			}
		case m.choiceIndex == 0 && key == termbox.KeySpace:
			if len(m.saveName) == 0 {
				fmt.Print("\a")
			} else {
				m.saveName = append(m.saveName, ' ')
			}
		case m.choiceIndex == 0 && ch != 0:
			// unicode.Punctuation is not included due to characters like /
			if ch != '_' && ch != '-' && !unicode.In(ch, unicode.Letter, unicode.Number, unicode.Symbol) {
				fmt.Print("\a")
			} else {
				m.saveName = append(m.saveName, ch)
			}
		case m.choiceIndex == 1 && (key == termbox.KeyBackspace || key == termbox.KeyBackspace2):
			if len(m.seed) == 0 {
				fmt.Print("\a")
			} else {
				m.seed = m.seed[:len(m.seed)-1]
			}
		case m.choiceIndex == 1 && key == termbox.KeySpace:
			m.seed = append(m.seed, ' ')
		case m.choiceIndex == 1 && ch != 0:
			m.seed = append(m.seed, ch)
		default:
			panic(fmt.Sprintf("%v, %v, %v", key, ch, mod))
		}

	default:
		panic(fmt.Sprintf("%v, %v, %v", key, ch, mod))
	}

	return true
}

func (m *mainMenuUI) inputMouse(x, y int) {
	switch m.state {
	case menuStateMain:
		switch {
		default:
			panic(fmt.Sprintf("%v, %v", x, y))
		}
	default:
		panic(fmt.Sprintf("%v, %v", x, y))
	}
}

func (m *mainMenuUI) loadGame(name string) {
	var w World
	f, err := os.Open(filepath.Join(SaveDirName, name+".sav"))
	if err != nil {
		m.err = err.Error()
		m.state = menuStateError
		return
	}
	defer f.Close()
	g, err := gzip.NewReader(f)
	if err != nil {
		m.err = err.Error()
		m.state = menuStateError
		return
	}
	defer g.Close()
	err = gob.NewDecoder(g).Decode(&w)
	if err != nil {
		m.err = err.Error()
		m.state = menuStateError
		return
	}
	err = w.AfterLoad()
	if err != nil {
		m.err = err.Error()
		m.state = menuStateError
		return
	}
	worldLock.Lock()
	world = &w
	worldLock.Unlock()
}

func (m *mainMenuUI) newGame() {
	m.state = menuStateNew
	m.choiceIndex = 0
}
