package termboxui

import (
	"github.com/nsf/termbox-go"
	"github.com/jpbetz/cellularautomata/io"
	"fmt"
	"github.com/mattn/go-runewidth"
)

type TermboxUI struct {
  backbuf []termbox.Cell
	refreshCh chan bool
	input chan io.InputEvent
}

func NewTermboxUI(input chan io.InputEvent) *TermboxUI {
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
	ui.backbuf[pos(position.X, position.Y)] = termbox.Cell{Ch: cell.Rune(), Fg: cell.FgAttribute()}
}

func (ui TermboxUI) Draw() {
	ui.refreshCh <- true
}

func (ui TermboxUI) handleInput() {
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
				ui.input <- io.Click {Position: io.Position{ev.MouseX, ev.MouseY}}
				ui.Draw()
			}
		case termbox.EventResize:
			ui.reallocBackBuffer(ev.Width, ev.Height)
			ui.Draw()
		default:
			ui.warn(fmt.Sprintf("Unexpected input: %v", ev))
			ui.Draw()
		}
	}
}

func (ui TermboxUI) warn(msg string) {
	_, h := termbox.Size()
	ui.tbprint(0,h-1, termbox.ColorBlue, termbox.ColorBlack, msg)
}

func (ui TermboxUI) tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	w, _ := termbox.Size()
	for _, c := range msg {
		if x >= w {
			return
		}
		ui.backbuf[pos(x, y)] = termbox.Cell{Ch: c, Fg: fg, Bg: bg}
		x += runewidth.RuneWidth(c)
	}
}

func (ui TermboxUI) refresh() {
	for range ui.refreshCh {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		copy(termbox.CellBuffer(), ui.backbuf)
		termbox.Flush()
	}
}