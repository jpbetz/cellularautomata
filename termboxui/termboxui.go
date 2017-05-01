package termboxui

import (
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/io"
)

type TermboxUI struct {
  backbuf []termbox.Cell
	refreshCh chan bool
	input *io.Input
}

func NewTermboxUI(input *io.Input) *TermboxUI {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	w, h := termbox.Size()

	ui := &TermboxUI{
		refreshCh: make(chan bool, 5),
		input: input,
		backbuf: make([]termbox.Cell, w*h),
	}

	return ui
}

func (ui TermboxUI) Run() {
	go ui.refresh()
	go ui.handleInput()
}

func (ui TermboxUI) Close() {
	defer termbox.Close()
	defer close(ui.refreshCh)
}

func (ui TermboxUI) reallocBackBuffer(w, h int) {
	ui.backbuf = make([]termbox.Cell, w*h)
}

func pos(x int, y int) int {
	w, _ := termbox.Size()
	return y * w + x
}

func (ui TermboxUI) Set(position io.Position, cell io.Cell) {
	ui.backbuf[pos(position.X, position.Y)] = termbox.Cell{Ch: cell.Rune(), Fg: cell.Attribute()}
}

func (ui TermboxUI) Draw() {
	ui.refreshCh <- true
}

func (ui TermboxUI) handleInput() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				ui.input.Quit <- true
				return
			} else if ev.Key == termbox.KeySpace {
				ui.input.PausePlay <- true
			}
		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft {
				ui.input.Click <- io.Position{ev.MouseX, ev.MouseY}
			}
		case termbox.EventResize:
			ui.reallocBackBuffer(ev.Width, ev.Height)
		}

		ui.Draw()
	}
}

func (ui TermboxUI) refresh() {
	for range ui.refreshCh {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		copy(termbox.CellBuffer(), ui.backbuf)
		termbox.Flush()
	}
}