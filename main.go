package main

import (
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/jpbetz/cellularautomata/termboxui"
	"github.com/jpbetz/cellularautomata/conway"
)

func main() {
	input := make(chan io.InputEvent, 100)

	ui := termboxui.NewTermboxUI(input)
	defer ui.Close()
	ui.Run()

	w, h := termbox.Size()

	cells := make([][]io.Cell, w)
	for row := range cells {
		cells[row] = make([]io.Cell, h)
		for i := 0; i < h; i++ {
			cells[row][i] = conway.Off
		}
	}
	game := conway.NewGameOfLife(cells, ui)
	eventClock := game.StartClock()
	game.Playing = true

	for {
		select {
		case in := <-input:
			switch in.(type) {
			case io.Quit:
				return
			case io.Click:
				game.Toggle(game.Cells, in.(io.Click).Position)
				ui.Draw()
			case io.Pause:
				if game.Playing {
					eventClock.Stop()
					game.Playing = false
				} else {
					eventClock = game.StartClock()
					game.Playing = true
				}
			}
		}
	}
}