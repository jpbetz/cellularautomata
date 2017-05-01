package io

import "github.com/nsf/termbox-go"

type Cell interface {
	Rune() rune
	FgAttribute() termbox.Attribute
}

type Position struct {
	X, Y int
}
