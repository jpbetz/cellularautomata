package termboxui

import (
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/io"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/jpbetz/cellularautomata/grid"
)

type TermboxUI struct {
	// rendering internals
  backbuf []termbox.Cell
	refreshCh chan bool

	// UI
	View *io.View

	// IO
	input chan io.InputEvent
}

func NewTermboxUI(input chan io.InputEvent) *TermboxUI {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	w, h := termbox.Size()

	return &TermboxUI{
		refreshCh: make(chan bool, 5),
		input: input,
		backbuf: make([]termbox.Cell, w*h),
	}
}

func (ui *TermboxUI) SetView(view *io.View) {
	ui.View = view
}

func (ui *TermboxUI) Run() {
	go ui.refresh()
	go ui.handleInput()
}

func (ui *TermboxUI) Close() {
	defer termbox.Close()
	defer close(ui.refreshCh)
}

func (ui *TermboxUI) reallocBackBuffer(w, h int) {
	ui.backbuf = make([]termbox.Cell, w*h)
}

func pos(x int, y int) int {
	w, _ := termbox.Size()
	return y * w + x
}

func (ui *TermboxUI) Set(position grid.Position, cell grid.Cell) {
	p := pos(position.X - ui.View.Offset.X, position.Y - ui.View.Offset.Y)

	if p >= 0 && p < len(ui.backbuf) {
		ui.backbuf[p] = termbox.Cell{Ch: cell.Rune(), Fg: cell.FgAttribute()}
	}
}

func (ui *TermboxUI) FullRefresh() {
	v := ui.View
	viewBound := v.Plane.Bounds().X2Y2
	p := v.Plane
	x, y := v.Offset.X, v.Offset.Y
	w, h := termbox.Size()
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if x + i <= viewBound.X && y + j < viewBound.Y {
				cell := p.Get(grid.Position{x + i, y + j})
				ui.backbuf[pos(i, j)] = termbox.Cell{Ch: cell.Rune(), Fg: cell.FgAttribute()}
			} else {
				ui.backbuf[pos(i, j)] = termbox.Cell{}
			}
		}
	}
}

func (ui *TermboxUI) Draw() {
	ui.refreshCh <- true
}

func (ui *TermboxUI) handleInput() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 'q' {
				ui.input <- io.Quit{}
				return
			} else if ev.Key == termbox.KeySpace {
				ui.input <- io.Pause{}
			}
		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft {
				ui.input <- io.Click {Position: grid.Position{ev.MouseX, ev.MouseY}}
				ui.Draw()
			}
		case termbox.EventResize:
			ui.reallocBackBuffer(ev.Width, ev.Height)
			//ui.FullRefresh()
			ui.Draw()
		default:
			ui.warn(fmt.Sprintf("Unexpected input: %v", ev))
			ui.Draw()
		}
	}
}

func (ui *TermboxUI) warn(msg string) {
	_, h := termbox.Size()
	ui.tbprint(0,h-1, termbox.ColorBlue, termbox.ColorBlack, msg)
}

func (ui *TermboxUI) tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	w, _ := termbox.Size()
	for _, c := range msg {
		if x >= w {
			return
		}
		ui.backbuf[pos(x, y)] = termbox.Cell{Ch: c, Fg: fg, Bg: bg}
		x += runewidth.RuneWidth(c)
	}
}

func (ui *TermboxUI) refresh() {
	for range ui.refreshCh {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		copy(termbox.CellBuffer(), ui.backbuf)
		termbox.Flush()
	}
}