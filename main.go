package main

import (
	"github.com/jpbetz/cellularautomata/io"
	"github.com/jpbetz/cellularautomata/termboxui"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/apps/langton"
	"github.com/jpbetz/cellularautomata/apps/conway"
	"github.com/jpbetz/cellularautomata/apps/wireworld"
)

func main() {
	//conwayMain()
	//langtonMain()
	wireworldMain()
}


func wireworldMain() {
	input := make(chan io.InputEvent, 100)

	ui := termboxui.NewTermboxUI(input)
	defer ui.Close()
	ui.Run()

	board := grid.NewBasicBoard(1000, 1000)
	view := &io.View{Plane: board, Offset: grid.Origin}
	ui.SetView(view)
	board.Initialize(wireworld.Cell{})
	game := wireworld.NewWireworld(board, ui)
	eventClock := game.StartClock()
	game.Playing = true

	for {
		select {
		case in := <-input:
			switch in.(type) {
			case io.Quit:
				return
			case io.Click:

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

func langtonMain() {
	input := make(chan io.InputEvent, 100)

	ui := termboxui.NewTermboxUI(input)
	defer ui.Close()
	ui.Run()

	board := grid.NewBasicBoard(1000, 1000)
	view := &io.View{Plane: board, Offset: grid.Origin}
	ui.SetView(view)
	board.Initialize(langton.Square{})
	game := langton.NewAnts(board, ui)
	eventClock := game.StartClock()
	game.Playing = true

	for {
		select {
		case in := <-input:
			switch in.(type) {
			case io.Quit:
				return
			case io.Click:

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

func conwayMain() {
	input := make(chan io.InputEvent, 100)

	ui := termboxui.NewTermboxUI(input)
	defer ui.Close()
	ui.Run()

	board := grid.NewBasicBoard(1000, 1000)
	view := &io.View{Plane: board, Offset: grid.Origin}
	ui.SetView(view)
	board.Initialize(conway.Life{Alive: false})
	game := conway.NewGameOfLife(board, ui)
	eventClock := game.StartClock()
	game.Playing = true

	for {
		select {
		case in := <-input:
			switch in.(type) {
			case io.Quit:
				return
			case io.Click:
				game.Toggle(game.Plane, in.(io.Click).Position)
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