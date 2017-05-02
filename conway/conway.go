package conway

import (
	"github.com/jpbetz/cellularautomata/io"
	"github.com/jpbetz/cellularautomata/engine"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/nsf/termbox-go"
)

const ALIVE = '█'
const DEAD = ' '

type Life struct {
	Alive bool
}

func (life Life) Rune() rune {
	if life.Alive {
		return ALIVE
	} else {
		return DEAD
	}
}

func (life Life) FgAttribute() termbox.Attribute {
	return termbox.ColorBlue
}

type GameOfLife struct {
	*engine.Engine
}

func asLife(cell grid.Cell) Life {
	life, ok := cell.(Life)
	if !ok {
		panic("Expected Life cell")
	}
	return life
}

var Alive = Life{Alive: true}
var Off = Life{Alive: false}

func NewGameOfLife(plane grid.Plane, ui io.Renderer) *GameOfLife {
	game := &GameOfLife{
		Engine: &engine.Engine{Plane: plane, UI: ui},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

func (g *GameOfLife) initialize() {
	example := [][]grid.Cell {
		{Off, Off, Off, Off, Off, Off },
		{Off, Off, Off, Off, Off, Off },
		{Off, Off, Off, Alive, Off, Off },
		{Off, Off, Off, Off, Alive, Off },
		{Off, Off, Alive, Alive, Alive, Off },
		{Off, Off, Off, Off, Off, Off },
	}

	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			g.Set(grid.Position{i, j}, example[i][j])
		}
	}
}

func (g *GameOfLife) UpdateCell(plane grid.Plane, position grid.Position) (state grid.Cell, changed bool) {

	x, y := position.X, position.Y
	bounds := plane.Bounds()
	if !bounds.Contains(position) {
		return
	}

	cell := asLife(plane.Get(position))

	neighbors := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if !(i == x && j == y) && i >= 0 && j >= 0 && i <= bounds.X2Y2.X && j <= bounds.X2Y2.Y {
				neighbor := asLife(plane.Get(grid.Position{i, j}))
				if neighbor.Alive {
					neighbors += 1
				}
			}
		}
	}
	if cell.Alive {
		if neighbors >= 2 && neighbors <= 3 {
			return Alive, true
		} else {
			return Off, true
		}
	} else if !cell.Alive && neighbors == 3 {
		return Alive, true
	}
	return Off, false
}

func (g *GameOfLife) Toggle(plane grid.Plane, position grid.Position) {
	if !plane.Bounds().Contains(position) {
		return
	}
	cell := asLife(plane.Get(position))
	if cell.Alive {
		g.Set(position, Life{Alive: false})
	} else {
		g.Set(position, Life{Alive: true})
	}
}