package io

type InputEvent interface {
	EventName() string
}

type Click struct {
	Position Position
}
func (Click) EventName() string {
	return "Click"
}

type Quit struct {}
func (Quit) EventName() string {
	return "Quit"
}

type Pause struct {}
func (Pause) EventName() string {
	return "Pause"
}
