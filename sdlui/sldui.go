package sdlui

import (
	"fmt"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/nsf/termbox-go"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

type UIUpdate struct {
	Position grid.Position
	Color    uint32
}

type UIRefresh struct {
}

type SdlUi struct {
	UpdateCh chan interface{}
	window   *sdl.Window
	surface  *sdl.Surface

	cells []*sdl.Rect

	// UI
	View *io.View

	// IO
	input chan io.InputEvent

	// number of cells wide and high
	Width  int32
	Height int32

	// width and height of each cell
	CellWidth  int32
	CellHeight int32
}

func NewSdlUi(input chan io.InputEvent, w, h int32, cellW, cellH int32) *SdlUi {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(w*cellW), int(h*cellH), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(fmt.Sprintf("Error creating sdl window: %v", err))
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	s := &SdlUi{
		UpdateCh:   make(chan interface{}, 5000),
		window:     window,
		surface:    surface,
		input:      input,
		Width:      w,
		Height:     h,
		CellWidth:  cellW,
		CellHeight: cellH,
	}

	s.cells = make([]*sdl.Rect, w*h)
	for i := int32(0); i < w; i++ {
		for j := int32(0); j < h; j++ {
			rect := &sdl.Rect{i * cellW, j * cellH, cellW, cellH}
			s.cells[s.pos(int(i), int(j))] = rect
			surface.FillRect(rect, toHex(termbox.ColorDefault))
		}
	}

	window.UpdateSurface()

	return s
}

func (s *SdlUi) SetView(view *io.View) {
	s.View = view
}

func (s *SdlUi) Run() {
}

func (s *SdlUi) Input() chan io.InputEvent {
	return s.input
}

func (s *SdlUi) Loop(done <-chan bool) {

	var lastMousePosition *grid.Position = nil
	for {
		if event := sdl.PollEvent(); event != nil {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				//log.Println("Quit event.")
				s.input <- io.Quit{}
				return
			case *sdl.MouseButtonEvent:
				//log.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
				//	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if t.Button == 1 && t.State&sdl.BUTTON_LEFT > 0 {
					s.input <- io.Click{Position: grid.Position{int(t.X) / int(s.CellWidth), int(t.Y) / int(s.CellHeight)}}
				}
			case *sdl.MouseMotionEvent:
				//log.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\ttxrel:%d\ttyrel:%d\tstate:%d\n",
				//	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel, t.State)

				newPosition := grid.Position{int(t.X) / int(s.CellWidth), int(t.Y) / int(s.CellHeight)}
				if t.State&sdl.BUTTON_LEFT > 0 && newPosition != *lastMousePosition {
					s.input <- io.Click{Position: newPosition}
				}
				lastMousePosition = &newPosition
			case *sdl.KeyUpEvent:
				//log.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				switch t.Keysym.Sym {
				case ' ':
					s.input <- io.Pause{}
				case 's':
					s.input <- io.Save{}
				case 'q':
					s.input <- io.Quit{}
				default:
					// do nothing
				}
			default:
				// do nothing
			}
		}

		select {
		case update := <-s.UpdateCh:
			switch update := update.(type) {
			case UIRefresh:
				s.Refresh()
			case UIUpdate:
				s.UpdateCell(update.Position, update.Color)
			}
		case <-done:
			log.Println("Done event recieved. Exiting Loop.")
			return
		default:
			// do nothing
		}
	}
}

func (s *SdlUi) Refresh() {
	s.window.UpdateSurface()
}

func (s *SdlUi) UpdateCell(position grid.Position, hexColor uint32) {
	if s.pos(position.X, position.Y) >= len(s.cells) {
		return
	}
	rect := s.cells[s.pos(position.X, position.Y)]
	s.surface.FillRect(rect, hexColor)
}

func (s *SdlUi) handleInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			s.input <- io.Quit{}
			return
		case *sdl.MouseButtonEvent:
			log.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
				t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
		case *sdl.KeyUpEvent:
			log.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			//if t.Keysym.Sym ==
		default:
			log.Println("Some event")
		}
		log.Println("Exiting handleInput.")
	}

}

func (s *SdlUi) Close() {
	sdl.Quit()
	s.window.Destroy()
}

func (s *SdlUi) Set(position grid.Position, change grid.Cell) {
	hexColor := toHex(change.FgAttribute())
	s.UpdateCh <- UIUpdate{position, hexColor}
}

func toHex(attribute termbox.Attribute) uint32 {
	switch attribute {
	case termbox.ColorBlue:
		return 0x00333fff
	case termbox.ColorRed:
		return 0x00ff3358
	case termbox.ColorYellow:
		return 0x00fff933
	case termbox.ColorWhite:
		return 0x00ffffff
	case termbox.ColorDefault:
		return 0x000e0e0e
	default:
		return 0x000e0e0e
	}
}

func (s *SdlUi) pos(x int, y int) int {
	return y*int(s.Width) + x
}

func (ui *SdlUi) Draw() {
	ui.UpdateCh <- UIRefresh{}
}

func (s *SdlUi) SetStatus(msg string) {
	fmt.Sprintf(msg)
}
