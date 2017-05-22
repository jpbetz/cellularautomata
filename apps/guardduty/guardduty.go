package guardduty

import (
	"fmt"
	"github.com/jpbetz/cellularautomata/engine"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/nsf/termbox-go"
	"os"
	"log"
)


type GuardDutyCommand struct {
	UI io.Renderer
}

func (c *GuardDutyCommand) Help() string {
	return "Guard Duty creates a simple waypoint circle that a guard walks around, using A* to navigate."
}

func (c *GuardDutyCommand) Run(args []string) int {
	guardDutyMain(c.UI)
	return 0
}

func (c *GuardDutyCommand) Synopsis() string {
	return "Guard Duty"
}


func guardDutyMain(ui io.Renderer) {
	f := setupLogging("logs/guardduty.log")
	defer f.Close()

	ui.Run()

	board := grid.NewBasicBoard(100, 100)
	view := &io.View{Plane: board, Offset: grid.Origin}
	ui.SetView(view)

	for x := board.Bounds().Corner1.X; x <= board.Bounds().Corner2.X; x++ {
		for y := board.Bounds().Corner1.Y; y <= board.Bounds().Corner2.Y; y++ {
			p := grid.Position{x, y}
			board.Set(p, Cell{
				State:    Empty,
				Position: p,
				Plane:    board,
			})
		}
	}

	game := NewGuardDuty(board, ui)
	eventClock := game.StartClock()
	game.Playing = true

	done := make(chan bool)
	go func() {
		for {
			in := <-ui.Input()
			switch event := in.(type) {
			case io.Quit:
				done <- true
				return
			case io.Click:
				cell := board.Get(event.Position).(Cell)
				if cell.State == Barrier {
					cell.State = Empty
				} else {
					cell.State = Barrier
				}
				game.Set(event.Position, cell)
				game.UI.Draw()
			case io.Pause:
				if game.Playing {
					eventClock.Stop()
					game.Playing = false
				} else {
					eventClock = game.StartClock()
					game.Playing = true
				}
			}
		}
	}()

	// main thread game loop
	ui.Loop(done)
}

func setupLogging(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	log.SetOutput(f)
	return f
}

type Unit interface{}

type Waypoint struct {
	position grid.Position
	next     *Waypoint
}

type Guard struct {
	nextWaypoint      *Waypoint
	nextWaypointRoute *grid.Path
}

type CellState int

const (
	Empty CellState = iota
	Barrier
)

type Cell struct {
	Plane    grid.Plane
	Position grid.Position
	Unit     Unit
	State    CellState
}

type CellNeighbor struct {
	Cell Cell
}

func (c CellNeighbor) GetNode() grid.Node {
	return c.Cell
}

func (c CellNeighbor) GetDistance() float64 {
	return 1
}

func (s Cell) Id() grid.NodeId {
	return s.Position
}

func (s Cell) GetNeighbors() []grid.Neighbor {
	if s.Plane == nil {
		panic("Cell has no plane.")
	}
	neighbors := s.Plane.GetNeighbors(s.Position)
	results := make([]grid.Neighbor, 0, len(neighbors))
	for _, neighbor := range neighbors {
		neighborCell := asCell(neighbor)
		if neighborCell.State != Barrier && (neighborCell.Position.X == s.Position.X || neighborCell.Position.Y == s.Position.Y) {
			results = append(results, CellNeighbor{neighborCell})
		}
	}
	return results
}

func (s Cell) Rune() rune {
	if s.Unit != nil {
		return ''
	}
	switch s.State {
	case Empty:
		return ' '
	default:
		return '█'
	}
}

func (s Cell) FgAttribute() termbox.Attribute {
	if s.Unit != nil {
		return termbox.ColorRed
	}
	switch s.State {
	case Empty:
		return termbox.ColorDefault
	case Barrier:
		return termbox.ColorBlue
	default:
		panic("Unsupported state")
	}
}

func (s Cell) BgAttribute() termbox.Attribute {
	return termbox.ColorDefault
}

func asCell(cell grid.Cell) Cell {
	life, ok := cell.(Cell)
	if !ok {
		panic("Expected Guardduty cell")
	}
	return life
}

type GuardDuty struct {
	*engine.Engine
}

func NewGuardDuty(plane grid.Plane, ui io.Renderer) *GuardDuty {
	game := &GuardDuty{
		Engine: &engine.Engine{Plane: plane, UI: ui},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

var Waypoint1 = &Waypoint{position: grid.Position{6, 1}}
var Waypoint2 = &Waypoint{position: grid.Position{6, 6}}
var Waypoint3 = &Waypoint{position: grid.Position{1, 1}}

var Guard1 = &Guard{nextWaypoint: Waypoint1}

var O = Cell{State: Empty}
var B = Cell{State: Barrier}
var G = Cell{State: Empty, Unit: Guard1}

func (g *GuardDuty) initialize() {
	Waypoint1.next = Waypoint2
	Waypoint2.next = Waypoint3
	Waypoint3.next = Waypoint1

	example := [][]Cell{
		{B, B, B, B, B, B, B, B},
		{B, G, O, O, O, O, O, B},
		{B, O, B, O, O, B, B, B},
		{B, B, B, O, O, O, B, B},
		{B, O, B, O, B, O, O, B},
		{B, O, O, O, B, B, O, B},
		{B, O, O, O, O, O, O, B},
		{B, B, B, B, B, B, B, B},
	}

	for x := 0; x < len(example); x++ {
		for y := 0; y < len(example[x]); y++ {
			example := example[x][y]
			p := grid.Position{x, y}
			current := asCell(g.Plane.Get(p))
			current.State = example.State
			current.Unit = example.Unit
			if p != current.Position {
				panic("p != current.Position")
			}
			g.Set(p, current)
		}
	}
	g.UI.SetStatus("GuardDuty")
}

func (g *GuardDuty) UpdateCell(plane grid.Plane, position grid.Position) []engine.CellUpdate {
	if !plane.Bounds().Contains(position) {
		return []engine.CellUpdate{}
	}

	cell := asCell(plane.Get(position))
	if cell.Unit != nil {
		switch unit := cell.Unit.(type) {
		case *Guard:
			guard := unit
			if guard.nextWaypointRoute == nil && guard.nextWaypoint != nil {
				g.UI.SetStatus(fmt.Sprintf("Next Waypoint: %v", guard.nextWaypoint.position))
				path, ok := findPath(cell, asCell(plane.Get(guard.nextWaypoint.position)))
				if ok {
					guard.nextWaypointRoute = path
				}
			}
			if guard.nextWaypointRoute != nil {
				var route = guard.nextWaypointRoute.Nodes
				if len(route) > 0 {
					var tail grid.Node
					tail, guard.nextWaypointRoute.Nodes = route[len(route)-1], route[:len(route)-1]
					nextPosition := tail.(Cell).Position
					nextCell := asCell(plane.Get(nextPosition))

					if nextCell.State == Barrier {
						guard.nextWaypointRoute = nil
					} else {
						nextCell.Unit = cell.Unit
						cell.Unit = nil
						return []engine.CellUpdate{
							{cell, cell.Position},
							{nextCell, nextCell.Position},
						}
					}
				} else {
					guard.nextWaypointRoute = nil

					guard.nextWaypoint = guard.nextWaypoint.next
				}
			}
		}
	}
	return []engine.CellUpdate{}
}

func costHuristic(p1, p2 grid.Node) float64 {
	return p1.(Cell).Position.DistanceTo(p2.(Cell).Position)
}

func findPath(start, goal Cell) (*grid.Path, bool) {
	return grid.FindPath(start, goal, costHuristic)
}
