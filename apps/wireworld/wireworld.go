package wireworld

import (
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/engine"
	"github.com/jpbetz/cellularautomata/io"
)

const ElectronHeadColor = termbox.ColorBlue
const ElectronTailColor = termbox.ColorRed
const ConductorColor = termbox.ColorYellow
const EmptyColor = termbox.ColorDefault

const ActiveRune = '█'
const EmptyRune = ' '

type State int
const (
	Empty State = iota
	ElectronHead
	ElectronTail
	Conductor
)

type Cell struct {
	State State
}

func (s Cell) Rune() rune {
	switch s.State {
	case Empty:
		return EmptyRune
	default:
		return ActiveRune
	}
}

func (s Cell) FgAttribute() termbox.Attribute {
	switch s.State {
	case Empty:
		return EmptyColor
	case ElectronHead:
		return ElectronHeadColor
	case ElectronTail:
		return ElectronTailColor
	case Conductor:
		return ConductorColor
	default:
		panic("Unsupported state")
	}
}

func (s Cell) BgAttribute() termbox.Attribute {
	return EmptyColor
}

type Wireworld struct {
	*engine.Engine
}

func asCell(cell grid.Cell) Cell {
	life, ok := cell.(Cell)
	if !ok {
		panic("Expected Wireworld cell")
	}
	return life
}

var Default = Cell{}

func NewWireworld(plane grid.Plane, ui io.Renderer) *Wireworld {
	game := &Wireworld{
		Engine: &engine.Engine{Plane: plane, UI: ui},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

var O = Cell{}
var C = Cell{State: Conductor}
var H = Cell{State: ElectronHead}
var T = Cell{State: ElectronTail}

func (g *Wireworld) initialize() {
	example := [][]grid.Cell {
		{O, C, C, T, H, C, C, C, C, O, O, O, O, O, O, O, O, O, O, O, O, O },
		{C, O, O, O, O, O, O, O, O, C, C, C, C, C, C, O, O, O, O, O, O, O },
		{O, C, H, T, C, C, C, C, C, O, O, O, O, O, O, C, O, O, O, O, O, O },
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, C, C, C, C, O, O, O, O },
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, C, O, O, C, C, C, C, C },
		{O, O, O, O, O, O, O, O, O, O, O, O, O, O, C, C, C, C, O, O, O, O },
		{O, C, C, C, C, C, C, C, C, O, O, O, O, O, O, C, O, O, O, O, O, O },
		{C, O, O, O, O, O, O, O, O, T, C, C, C, C, C, O, O, O, O, O, O, O },
		{O, C, H, T, C, C, C, C, H, O, O, O, O, O, O, O, O, O, O, O, O, O },
	}

	for i := 0; i < len(example); i++ {
		for j := 0; j < len(example[0]); j++ {
			g.Set(grid.Position{j, i}, example[i][j])
		}
	}
	g.UI.SetStatus("WireWorld")
}

func (g *Wireworld) UpdateCell(plane grid.Plane, position grid.Position) []engine.CellUpdate {

	if !plane.Bounds().Contains(position) {
		return []engine.CellUpdate {}
	}

	cell := asCell(plane.Get(position))
	switch cell.State {
	case ElectronHead:
		cell.State = ElectronTail
		return []engine.CellUpdate {{cell, position}}
	case ElectronTail:
		cell.State = Conductor
		return []engine.CellUpdate {{cell, position}}
	case Conductor:
		neighboringElectronHeads := 0
		for _, neighbor := range plane.GetNeighbors(position) {
			if asCell(neighbor).State == ElectronHead {
				neighboringElectronHeads++
			}
		}
		if neighboringElectronHeads > 0 && neighboringElectronHeads < 3 {
			cell.State = ElectronHead
			return []engine.CellUpdate {{cell, position}}
		} else {
			return []engine.CellUpdate {}
		}
	case Empty:
		return []engine.CellUpdate {}
	default:
		panic("Unsupported State")
	}
}
