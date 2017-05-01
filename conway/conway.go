package conway

import (
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/jpbetz/cellularautomata/engine"
)

var ALIVE = '█'
var DEAD = ' '

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
	engine.Engine
}

func asLife(cell io.Cell) Life {
	life, ok := cell.(Life)
	if !ok {
		panic("Expected Life cell")
	}
	return life
}

var Alive = Life{Alive: true}
var Off = Life{Alive: false}

func NewGameOfLife(cells [][]io.Cell, ui io.Renderer) *GameOfLife {
	game := &GameOfLife{
		Engine: engine.Engine{Cells: cells, UI: ui},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

func (this GameOfLife) initialize() {
	example := [][]io.Cell {
		{Off, Off, Off, Off, Off, Off },
		{Off, Off, Off, Off, Off, Off },
		{Off, Off, Off, Alive, Off, Off },
		{Off, Off, Off, Off, Alive, Off },
		{Off, Off, Alive, Alive, Alive, Off },
		{Off, Off, Off, Off, Off, Off },
	}

	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			this.Set(io.Position{i, j}, example[i][j])
		}
	}
}

func (this GameOfLife) UpdateCell(cells [][]io.Cell, x, y int) (state io.Cell, changed bool) {
	w, h := termbox.Size()
	cell := asLife(cells[x][y])

	neighbors := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if !(i == x && j == y) && i >= 0 && j >= 0 && i < w && j < h {
				neighbor := asLife(cells[i][j])
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

func (this GameOfLife) Toggle(cells [][]io.Cell, position io.Position) {
	cell := asLife(cells[position.X][position.Y])
	if cell.Alive {
		this.Set(position, Life{Alive: false})
	} else {
		this.Set(position, Life{Alive: true})
	}
}