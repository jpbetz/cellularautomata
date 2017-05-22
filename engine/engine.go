package engine

import (
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"time"
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
	eventClock := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range eventClock.C {
			e.clockEvent()
		}
	}()
	return eventClock
}

func (e *Engine) clockEvent() {
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
		e.Set(change.Position, change.State)
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
