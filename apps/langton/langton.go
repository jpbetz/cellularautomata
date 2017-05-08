package langton

import (
	"github.com/jpbetz/cellularautomata/engine"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/nsf/termbox-go"
)

const AntColor = termbox.ColorBlue
const Empty = ' '
const White = termbox.ColorWhite
const Black = termbox.ColorDefault

type Ant struct {
	orientation grid.Orientation
}

type Square struct {
	White bool
	Ant   *Ant
}

func (s Square) Rune() rune {
	if s.Ant != nil {
		switch s.Ant.orientation {
		case grid.Up:
			return ''
		case grid.Down:
			return ''
		case grid.Left:
			return ''
		case grid.Right:
			return ''
		default:
			panic("Unsupported Orientation value")
		}
	} else {
		return Empty
	}
}

func (s Square) FgAttribute() termbox.Attribute {
	return AntColor
}

func (s Square) BgAttribute() termbox.Attribute {
	if s.White {
		return White
	} else {
		return Black
	}
}

type Ants struct {
	*engine.Engine
}

func asSquare(cell grid.Cell) Square {
	life, ok := cell.(Square)
	if !ok {
		panic("Expected Langton's Ant Square cell")
	}
	return life
}

var AntStart = Square{Ant: &Ant{}}
var Default = Square{}

func NewAnts(plane grid.Plane, ui io.Renderer) *Ants {
	game := &Ants{
		Engine: &engine.Engine{Plane: plane, UI: ui},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

func (g *Ants) initialize() {
	w, h := termbox.Size()
	g.Set(grid.Position{w / 2, h / 2}, AntStart)
	g.UI.SetStatus("Langton's Ants")
}

func (g *Ants) UpdateCell(plane grid.Plane, position grid.Position) []engine.CellUpdate {

	if !plane.Bounds().Contains(position) {
		return []engine.CellUpdate{}
	}

	cell := asSquare(plane.Get(position))

	if cell.Ant != nil {
		var nextOrientation grid.Orientation
		if cell.White {
			nextOrientation = cell.Ant.orientation.Rotate(1)
		} else {
			nextOrientation = cell.Ant.orientation.Rotate(-1)
		}

		updatedCell := Square{Ant: nil, White: !cell.White}

		nextPosition := position.Translate(nextOrientation, 1)
		if !plane.Bounds().Contains(nextPosition) {
			return []engine.CellUpdate{}
		}

		nextCell := asSquare(plane.Get(nextPosition))
		updatedNextCell := Square{Ant: &Ant{nextOrientation}, White: nextCell.White}

		return []engine.CellUpdate{{updatedCell, position}, {updatedNextCell, nextPosition}}
	} else {
		return []engine.CellUpdate{}
	}
}
