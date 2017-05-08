package grid

import (
	"github.com/nsf/termbox-go"
	"math"
)

type Plane interface {
	Get(position Position) Cell
	GetNeighborPositions(p Position) []Position
	GetNeighbors(p Position) []Cell
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

func (p1 Position) DistanceTo(p2 Position) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(p2.X-p1.X)), 2) + math.Pow(math.Abs(float64(p2.Y-p1.Y)), 2))
}

type Orientation int

const (
	Up Orientation = iota
	Right
	Down
	Left
)

func (o Orientation) Rotate(halfTurnsCW int) Orientation {
	applied := (int(o) + halfTurnsCW) % 4
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
	Corner1 Position
	Corner2 Position
}

func (r Rectangle) Contains(position Position) bool {
	return position.X >= r.Corner1.X && position.X <= r.Corner2.X && position.Y >= r.Corner1.Y && position.Y <= r.Corner2.Y
}
