package engine

import (
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"time"
	"fmt"
)

type CellUpdate struct {
	State    grid.Cell
	Position grid.Position
}

type UpdateHandler interface {
	UpdateCell(plane grid.Plane, position grid.Position) []CellUpdate
}

type Engine struct {
	Plane   grid.Plane
	UI      io.Renderer
	Playing bool
	Handler UpdateHandler
}

func (e *Engine) StartClock() *time.Ticker {
	fmt.Print("Starting event clock\n")
	eventClock := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range eventClock.C {
			e.clockEvent()
		}
	}()
	fmt.Print("Event clock started\n")
	return eventClock
}

func (e *Engine) clockEvent() {
	//w, h := termbox.Size()
	//bounds := grid.Rectangle{
	//  grid.Position {0,0},
	//  grid.Position {e.UI.(*sdlui.SdlUi).Width,e.UI.(*sdlui.SdlUi).Height},
	//}
	bounds := e.Plane.Bounds()
	changes := []CellUpdate{}
	for i := bounds.Corner1.X; i < bounds.Corner2.X; i++ {
		for j := bounds.Corner1.Y; j < bounds.Corner2.Y; j++ {
			updates := e.Handler.UpdateCell(e.Plane, grid.Position{X: i, Y: j})
			for _, update := range updates {
				changes = append(changes, update)
			}
		}
	}
	for _, change := range changes {
		e.Plane.Set(change.Position, change.State)
		e.UI.Set(change.Position, change.State)
	}
	e.UI.Draw()
}

func (e *Engine) Set(position grid.Position, cell grid.Cell) {
	if !e.Plane.Bounds().Contains(position) {
		return
	}
	e.Plane.Set(position, cell)
	e.UI.Set(position, cell)
}
