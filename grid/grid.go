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
	BgAttribute() termbox.Attribute
}

type Position struct {
	X, Y int
}

type Orientation int

const (
	Up Orientation = iota
	Right
	Down
	Left
)

func (o Orientation) Rotate(halfTurns int) Orientation {
	applied := (int(o) + halfTurns) % 4
	if applied < 0 {
		applied += 4
	}
	return Orientation(applied)
}

func (p Position) Translate(orientation Orientation, distance int) Position {
	switch orientation {
	case Up:
		return Position{p.X, p.Y - distance}
	case Down:
		return Position{p.X, p.Y + distance}
	case Left:
		return Position{p.X - distance, p.Y}
	case Right:
		return Position{p.X + distance, p.Y}
	default:
		panic("Unsupported Orientation value")
	}
}

var Origin = Position{0, 0}

type Rectangle struct {
	X1Y1 Position
  X2Y2 Position
}

func (r Rectangle) Contains(position Position) bool {
	return position.X >= r.X1Y1.X && position.X <= r.X2Y2.X && position.Y >= r.X1Y1.Y && position.Y <= r.X2Y2.Y
}