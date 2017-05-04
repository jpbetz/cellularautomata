package engine

import (
	"github.com/jpbetz/cellularautomata/io"
	"time"
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/grid"
)

type CellUpdate struct {
  State grid.Cell
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
	eventClock := time.NewTicker(time.Millisecond * 10)
	go func() {
		for range eventClock.C {
			e.clockEvent()
			e.UI.Draw()
		}
	}()
	return eventClock
}

func (e *Engine) clockEvent() {
	w, h := termbox.Size()
	changes := []CellUpdate {}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
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