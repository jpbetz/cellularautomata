package io

import "github.com/jpbetz/cellularautomata/grid"

type Renderer interface {
	Run()
	Close()
	Set(position grid.Position, change grid.Cell)
	Draw()
	SetStatus(msg string)
}

type View struct {
	Plane  grid.Plane
	Offset grid.Position
}
