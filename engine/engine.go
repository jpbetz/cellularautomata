package engine

import (
	"github.com/jpbetz/cellularautomata/io"
	"time"
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/grid"
)

type UpdateHandler interface {
	UpdateCell(plane grid.Plane, position grid.Position) (state grid.Cell, changed bool)
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
			e.UI.Draw()
		}
	}()
	return eventClock
}

type CellChange struct {
	Position grid.Position
	Cell grid.Cell
}

func (e *Engine) clockEvent() {
	w, h := termbox.Size()
	changes := []CellChange {}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			state, changed := e.Handler.UpdateCell(e.Plane, grid.Position{i, j})
			if changed {
				change := CellChange{grid.Position{i, j}, state}
				e.UI.Set(change.Position, change.Cell)
				changes = append(changes, change)
			}
		}
	}
	for _, change := range changes {
		e.Plane.Set(change.Position, change.Cell)
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