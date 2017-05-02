package grid

import "github.com/nsf/termbox-go"

type Plane interface {
	Get(position Position) Cell
	Set(position Position, cell Cell)
	Bounds() Rectangle
}

type Cell interface {
	Rune() rune
	FgAttribute() termbox.Attribute
}

type Position struct {
	X, Y int
}

var Origin = Position{0, 0}

type Rectangle struct {
	X1Y1 Position
  X2Y2 Position
}

func (r Rectangle) Contains(position Position) bool {
	return position.X >= r.X1Y1.X && position.X <= r.X2Y2.X && position.Y >= r.X1Y1.Y && position.Y <= r.X2Y2.Y
}