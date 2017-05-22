package io

import "github.com/jpbetz/cellularautomata/grid"

type InputEvent interface {
	EventName() string
}

type Click struct {
	Position grid.Position
}

func (Click) EventName() string {
	return "Click"
}

type Quit struct{}

func (Quit) EventName() string {
	return "Quit"
}

type Pause struct{}

func (Pause) EventName() string {
	return "Pause"
}

type Save struct{}

func (Save) EventName() string {
	return "Save"
}
