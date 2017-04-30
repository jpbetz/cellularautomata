package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

var backbuf []termbox.Cell
var nextbuf []termbox.Cell

var on = '█'
var off = ' '

type attrFunc func(int) (rune, termbox.Attribute, termbox.Attribute)

func reallocBackBuffer(w, h int) {
	backbuf = make([]termbox.Cell, w*h)
	nextbuf = make([]termbox.Cell, w*h)
}

func clockEvent() {
	w, h := termbox.Size()
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			updateCell(i, j)
		}
	}
	copy(backbuf, nextbuf)
}

func StartClock(eventChannel chan<- bool) *time.Ticker {
	eventClock := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range eventClock.C {
			clockEvent()
			eventChannel <- true
		}
	}()
	return eventClock
}

func updateCell(x, y int) {
	w, h := termbox.Size()
	cell := backbuf[pos(x, y)]
	neighbors := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if !(i == x && j == y) && i >= 0 && j >= 0 && i < w && j < h {
				neighbor := backbuf[pos(i, j)]
				if neighbor.Ch == on {
					neighbors += 1
				}
			}
		}
	}
	isAlive := cell.Ch == on
	if isAlive {
		if neighbors >= 2 && neighbors <= 3 {
			nextbuf[pos(x, y)] = termbox.Cell{Ch: on, Fg: termbox.ColorWhite}
		} else {
			nextbuf[pos(x, y)] = termbox.Cell{Ch: off, Fg: termbox.ColorWhite}
		}
	} else if !isAlive && neighbors == 3 {
		nextbuf[pos(x, y)] = termbox.Cell{Ch: on, Fg: termbox.ColorWhite}
	}
}

func toggle(x, y int) {
	cell := backbuf[pos(x, y)]
	if cell.Ch == on {
		backbuf[pos(x, y)] = termbox.Cell{Ch: off, Fg: termbox.ColorWhite}
	} else {
		backbuf[pos(x, y)] = termbox.Cell{Ch: on, Fg: termbox.ColorWhite}
	}
}

func pos(x int, y int) int {
	w, _ := termbox.Size()
	return y * w + x
}

func initialize() {
	example := [][]rune {
		{ off, off, off, off, off, off },
		{ off, on, on, off, off, off },
		{ off, on, off, off, off, off },
		{ off, off, off, off, on, off },
		{ off, off, off, on, on, off },
		{ off, off, off, off, off, off },
	}
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			backbuf[pos(i, j)] = termbox.Cell{Ch: example[i][j], Fg: termbox.ColorWhite}
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	reallocBackBuffer(termbox.Size())

	initialize()

	refresh := make(chan bool, 5)
	defer close(refresh)
	go renderer(refresh)

	eventClock := StartClock(refresh)
	running := true

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				break mainloop
			} else if ev.Key == termbox.KeySpace {
				if running {
					eventClock.Stop()
					running = false
				} else {
					eventClock = StartClock(refresh)
					running = true
				}

			}
		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft {
				toggle(ev.MouseX, ev.MouseY)
			}
		case termbox.EventResize:
			reallocBackBuffer(ev.Width, ev.Height)
		}

		refresh <- true
	}
}

func renderer(refresh chan bool) {
	for range refresh {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		copy(termbox.CellBuffer(), backbuf)
		termbox.Flush()
	}
}