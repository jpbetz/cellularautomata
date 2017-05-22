package sdlui

import (
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
        "github.com/veandco/go-sdl2/sdl"
	"github.com/nsf/termbox-go"
	"fmt"
	"log"
)

type UIUpdate struct {
	Position grid.Position
	Color uint32
}

type UIRefresh struct {

}

type SdlUi struct {
	UpdateCh chan interface{}
	window *sdl.Window
	surface *sdl.Surface

	cells []*sdl.Rect

	// UI
	View *io.View

	// IO
	input chan io.InputEvent

	Width int
	Height int
}

var (
  w, h int32 = 800, 600
  cellW, cellH int32 = 20, 20
)

func NewSdlUi(input chan io.InputEvent) *SdlUi {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(w), int(h), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(fmt.Sprintf("Error creating sdl window: %v", err))
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	cells := make([]*sdl.Rect, w/cellW*h/cellH)
	for i := int32(0); i < w/cellW; i++ {
		for j := int32(0); j < h/cellH; j++ {
			rect := &sdl.Rect{i*cellW, j*cellH, cellW, cellH}
			cells[pos(int(i), int(j))] = rect
			surface.FillRect(rect, toHex(termbox.ColorDefault))
		}
	}

	window.UpdateSurface()

	return &SdlUi{
		UpdateCh: make(chan interface{}, 5000),
		window: window,
		surface: surface,
		cells: cells,
		input: input,
		Width: int(w/cellW),
		Height: int(h/cellH),
	}
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

	for {
		if  event := sdl.PollEvent(); event != nil {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				log.Println("Quit event.")
				s.input <- io.Quit{}
				return
			case *sdl.MouseButtonEvent:
				log.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if t.Button == 1 && t.State == 0 {
					s.input <- io.Click{Position: grid.Position{int(t.X)/int(cellW), int(t.Y)/int(cellH)}}
				}
			case *sdl.KeyUpEvent:
				log.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				switch t.Keysym.Sym {
				case ' ':
					s.input <- io.Pause{}
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
	if pos(position.X, position.Y) >= len(s.cells) {
		return
	}
	rect := s.cells[pos(position.X, position.Y)]
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

func pos(x int, y int) int {
	return y*(int(w)/int(cellW)) + x
}

func (ui *SdlUi) Draw() {
	ui.UpdateCh <- UIRefresh{}
}

func (s *SdlUi) SetStatus(msg string) {
	fmt.Sprintf(msg)
}
