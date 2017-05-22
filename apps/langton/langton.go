package langton

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

type LangtonCommand struct {
	UI io.Renderer
}

func (c *LangtonCommand) Help() string {
	return "Langton's Ants simulates an ant that walks a route that depends on state of ground cells."
}

func (c *LangtonCommand) Run(args []string) int {
	langtonMain(c.UI)
	return 0
}

func (c *LangtonCommand) Synopsis() string {
	return "Langton's Ants"
}

func langtonMain(ui io.Renderer) {
	f := setupLogging("logs/langton.log")
	defer f.Close()

	ui.Run()

	board := grid.NewBasicBoard(1000, 1000)
	view := &io.View{Plane: board, Offset: grid.Origin}
	ui.SetView(view)
	board.Initialize(Square{})
	game := NewAnts(board, ui)
	eventClock := game.StartClock()
	game.Playing = true

	done := make(chan bool)

	go func() {
		for {
			in := <-ui.Input()
			switch in.(type) {
			case io.Quit:
				done <- true
				return
			case io.Click:

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
			return '^'
		case grid.Down:
			return '_'
		case grid.Left:
			return '<'
		case grid.Right:
			return '>'
		default:
			panic("Unsupported Orientation value")
		}
	} else {
		return Empty
	}
}

func (s Square) FgAttribute() termbox.Attribute {
	if s.Ant != nil {
		return AntColor
	} else if s.White {
		return White
	} else {
		return Black
	}
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
		Engine: &engine.Engine{Plane: plane, UI: ui, ClockSpeed: time.Millisecond * 100},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

func (g *Ants) initialize() {
	position := grid.Position{20, 20}
	log.Printf("initializing ant at %d, %d\n", position.X, position.Y)
	g.Set(position, AntStart)
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
