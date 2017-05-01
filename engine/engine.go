package engine

import (
	"github.com/jpbetz/cellularautomata/io"
	"time"
	"github.com/nsf/termbox-go"
)

type UpdateHandler interface {
	UpdateCell(cells [][]io.Cell, x, y int) (state io.Cell, changed bool)
}

type Engine struct {
	Cells   [][]io.Cell
	UI      io.Renderer
	Playing bool
	Handler UpdateHandler
}

func (this Engine) StartClock() *time.Ticker {
	eventClock := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range eventClock.C {
			this.clockEvent()
			this.UI.Draw()
		}
	}()
	return eventClock
}

type CellChange struct {
	Position io.Position
	Cell io.Cell
}

func (this Engine) clockEvent() {
	w, h := termbox.Size()
	changes := []CellChange {}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			state, changed := this.Handler.UpdateCell(this.Cells, i, j)
			if changed {
				change := CellChange{io.Position{i, j}, state}
				this.UI.Set(change.Position, change.Cell)
				changes = append(changes, change)
			}
		}
	}
	for _, change := range changes {
		this.Cells[change.Position.X][change.Position.Y] = change.Cell
	}
	this.UI.Draw()
}

func (this Engine) Set(position io.Position, cell io.Cell) {
	this.Cells[position.X][position.Y] = cell
	this.UI.Set(position, cell)
}