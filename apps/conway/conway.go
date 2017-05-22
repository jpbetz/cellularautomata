package conway

import (
	"fmt"
	"github.com/jpbetz/cellularautomata/engine"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/nsf/termbox-go"
	"log"
	"os"
	"time"
)

const ALIVE = '█'
const DEAD = ' '

type ConwayCommand struct {
	UI io.Renderer
}

func (c *ConwayCommand) Help() string {
	return "Conway's Game of Life"
}

func (c *ConwayCommand) Run(args []string) int {
	conwayMain(c.UI)
	return 0
}

func (c *ConwayCommand) Synopsis() string {
	return "Conway's Game of Life"
}

func conwayMain(ui io.Renderer) {
	f := setupLogging("logs/conway.log")
	defer f.Close()

	ui.Run()

	board := grid.NewBasicBoard(80, 80)
	view := &io.View{Plane: board, Offset: grid.Origin}
	ui.SetView(view)
	board.Initialize(Life{Alive: false})
	game := NewGameOfLife(board, ui)
	eventClock := game.StartClock()
	game.Playing = true

	done := make(chan bool)

	go func() {
		for {
			in := <-ui.Input()
			switch event := in.(type) {
			case io.Quit:
				done <- true
				return
			case io.Click:
				cell := game.Toggle(game.Plane, event.Position)
				if cell != nil {
					ui.Draw()
				}
			case io.Pause:
				if game.Playing {
					eventClock.Stop()
					game.Playing = false
				} else {
					eventClock = game.StartClock()
					game.Playing = true
				}
			}
		}
	}()

	ui.Loop(done)
}

func setupLogging(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	log.SetOutput(f)
	return f
}

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
	if life.Alive {
		return termbox.ColorBlue
	} else {
		return termbox.ColorDefault
	}
}

func (life Life) BgAttribute() termbox.Attribute {
	return termbox.ColorDefault
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
		Engine: &engine.Engine{Plane: plane, UI: ui, ClockSpeed: time.Millisecond * 250},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

func (g *GameOfLife) initialize() {
	example := [][]grid.Cell{
		{Off, Off, Off, Off, Off, Off},
		{Off, Off, Off, Off, Off, Off},
		{Off, Off, Off, Alive, Off, Off},
		{Off, Off, Off, Off, Alive, Off},
		{Off, Off, Alive, Alive, Alive, Off},
		{Off, Off, Off, Off, Off, Off},
	}

	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			g.Set(grid.Position{i, j}, example[i][j])
		}
	}
	g.UI.SetStatus("Conway's game of life")
}

func (g *GameOfLife) UpdateCell(plane grid.Plane, position grid.Position) []engine.CellUpdate {

	bounds := plane.Bounds()
	if !bounds.Contains(position) {
		return []engine.CellUpdate{}
	}

	cell := asLife(plane.Get(position))

	neighbors := 0
	for _, neighbor := range plane.GetNeighbors(position) {
		neighbor := asLife(neighbor)
		if neighbor.Alive {
			neighbors += 1
		}
	}
	if cell.Alive {
		if neighbors >= 2 && neighbors <= 3 {
			return []engine.CellUpdate{{Alive, position}}
		} else {
			return []engine.CellUpdate{{Off, position}}
		}
	} else if !cell.Alive && neighbors == 3 {
		return []engine.CellUpdate{{Alive, position}}
	}
	return []engine.CellUpdate{}
}

func (g *GameOfLife) Toggle(plane grid.Plane, position grid.Position) grid.Cell {
	if !plane.Bounds().Contains(position) {
		return nil
	}
	cell := asLife(plane.Get(position))
	if cell.Alive {
		g.Set(position, Life{Alive: false})
	} else {
		g.Set(position, Life{Alive: true})
	}
	return cell
}
