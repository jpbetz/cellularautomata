package io

import "github.com/nsf/termbox-go"

type Cell interface {
	Rune() rune
	Attribute() termbox.Attribute
}

type Position struct {
	X, Y int
}
