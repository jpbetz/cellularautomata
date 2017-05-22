package guardduty

import (
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/jpbetz/cellularautomata/engine"
	"github.com/jpbetz/cellularautomata/flatbuffers/region"
	"github.com/jpbetz/cellularautomata/grid"
	"github.com/jpbetz/cellularautomata/io"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"log"
	"os"
	"time"
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

	board := grid.NewBasicBoard(40, 40)
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
			case io.Save:
				log.Printf("Writing log file: %s\n", saveDataFile)
				buf := game.Save(board.Cells, board.W, board.H)
				if err := ioutil.WriteFile(saveDataFile, buf, 0664); err != nil {
					log.Printf("Failed to write file %v\n", err)
				}
				log.Printf("Wrote log file: %s\n", saveDataFile)
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
		Engine: &engine.Engine{Plane: plane, UI: ui, ClockSpeed: time.Millisecond * 100},
	}
	game.Engine.Handler = game
	game.initialize()
	return game
}

var O = Cell{State: Empty}
var B = Cell{State: Barrier}

var initialDataFile = "data/guardduty/initial.dat"
var saveDataFile = "data/guardduty/save.dat"

func (g *GuardDuty) initialize() {

	file := saveDataFile
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Println("Save file not found. Loading initial file.")
		file = initialDataFile
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("Unable to read file: %s: %#v", file, err))
	}

	example, w, h := g.Load(buf)
	for i := 0; i < len(example); i++ {
		example := example[i]
		p := grid.Position{i % w, i / w}
		current := asCell(g.Plane.Get(p))
		current.State = example.State
		current.Unit = example.Unit
		if p != current.Position {
			panic("p != current.Position")
		}
		g.Set(p, current)
	}
	g.UI.SetStatus(fmt.Sprintf("GuardDuty (%d, %d)", w, h))
}

func (g *GuardDuty) Load(buf []byte) ([]Cell, int, int) {
	basicBoard := region.GetRootAsBasicBoard(buf, 0)
	log.Print("Loaded board\n")
	plane := basicBoard.Plane(&region.Plane{})

	planeW, planeH := int(plane.W()), int(plane.H())
	log.Printf("Loaded plane (w: %d, h: %d, total tiles: %d)\n", planeH, planeW, plane.TilesLength())

	cells := make([]Cell, plane.TilesLength())
	for i := 0; i < plane.TilesLength(); i++ {
		tile := &region.Tile{}
		plane.Tiles(tile, i)
		//log.Printf("Loaded tile at %d (%d, %d), type: %d\n", i, (i % planeW), (i / planeW), tile.TileType())
		if int(tile.TileType()) == region.TileTypeBarrier {
			cells[i] = B
		} else {
			cells[i] = O
		}
	}
	log.Printf("Loaded %d tiles\n", plane.TilesLength())

	guard := &region.GuardUnit{}
	basicBoard.Guard(guard)
	guardPosition := &region.Position{}
	guard.Position(guardPosition)
	log.Printf("Loaded guard at (%d, %d)\n", guardPosition.X(), guardPosition.Y())
	waypoints := make([]*Waypoint, guard.WaypointsLength())
	for i := 0; i < guard.WaypointsLength(); i++ {
		position := &region.Position{}
		guard.Waypoints(position, i)
		log.Printf("Loaded waypoint at (%d, %d)\n", position.X(), position.Y())
		waypoint := &Waypoint{
			position: grid.Position{int(position.X()), int(position.Y())},
		}
		if i > 0 {
			waypoints[i-1].next = waypoint
		}
		if i == guard.WaypointsLength()-1 {
			waypoint.next = waypoints[0]
		}
		waypoints[i] = waypoint
	}
	log.Printf("Loaded waypoints: %#v\n", waypoints)

	guardIdx := int(guardPosition.X()) + (planeW * int(guardPosition.Y()))
	log.Printf("Loading guard at idx: %d, x,y = %d, %d", guardIdx, guardPosition.X(), guardPosition.Y())
	cells[guardIdx] = Cell{
		State: Empty,
		Unit:  &Guard{nextWaypoint: waypoints[0]},
	}
	return cells, planeW, planeH
}

func (g *GuardDuty) Save(cells []grid.Cell, w, h int) []byte {
	builder := flatbuffers.NewBuilder(0)

	// tiles
	log.Printf("Saving %d cells", len(cells))
	tileEnds := make([]flatbuffers.UOffsetT, len(cells))
	var guard *Guard
	var guardPosition grid.Position
	for i, gridCell := range cells {
		cell := gridCell.(Cell)
		var state int
		if cell.State == Barrier {
			state = region.TileTypeBarrier
		} else {
			state = region.TileTypeEmpty
		}
		region.TileStart(builder)
		region.TileAddTileType(builder, int32(state))
		tileEnds[i] = region.TileEnd(builder)
		var ok bool
		if cell.Unit != nil {
			guard, ok = cell.Unit.(*Guard)
			if !ok {
				panic(fmt.Sprintf("Failed to convert cell unit to guard: %v", cell.Unit))
			}
			guardPosition = cell.Position
		}
	}

	// plane tiles vector
	log.Printf("Writing %d tiles", len(tileEnds))
	region.PlaneStartTilesVector(builder, len(tileEnds))
	for i := len(tileEnds) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(tileEnds[i])
	}
	tilesVectorEnd := builder.EndVector(len(tileEnds))

	// plane
	region.PlaneStart(builder)
	region.PlaneAddW(builder, int32(w))
	region.PlaneAddH(builder, int32(h))
	region.PlaneAddTiles(builder, tilesVectorEnd)
	planeEnd := region.PlaneEnd(builder)

	// waypoint positions
	start := guard.nextWaypoint

	waypoints := make([]grid.Position, 0)
	current := start
	for ; current.next != start; current = current.next {
		waypoints = append(waypoints, current.position)
	}
	waypoints = append(waypoints, current.position)

	// waypoint positions vector
	region.GuardUnitStartWaypointsVector(builder, len(waypoints))
	for i := len(waypoints) - 1; i >= 0; i-- {
		log.Printf("Writing waypoint %d (%d, %d)", i, int32(waypoints[i].X), int32(waypoints[i].Y))
		region.CreatePosition(builder, int32(waypoints[i].X), int32(waypoints[i].Y))
	}
	waypointsVectorEnd := builder.EndVector(len(waypoints))

	// guard unit
	log.Printf("Writing guard position (%d, %d)", int32(guardPosition.X), int32(guardPosition.Y))
	guardUnitPosition := region.CreatePosition(builder, int32(guardPosition.X), int32(guardPosition.Y))

	region.GuardUnitStart(builder)
	region.GuardUnitAddPosition(builder, guardUnitPosition)
	region.GuardUnitAddWaypoints(builder, waypointsVectorEnd)
	guardEnd := region.GuardUnitEnd(builder)

	// basic board
	region.BasicBoardStart(builder)
	region.BasicBoardAddPlane(builder, planeEnd)
	region.BasicBoardAddGuard(builder, guardEnd)
	basicBoardEnd := region.BasicBoardEnd(builder)

	builder.Finish(basicBoardEnd)
	return builder.Bytes[builder.Head():]
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
