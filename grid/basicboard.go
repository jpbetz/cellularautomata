package grid

type BasicBoard struct {
	cells []Cell
	w, h int
}

func NewBasicBoard(w, h int) *BasicBoard {
	return &BasicBoard{
		cells: make([]Cell, w * h),
		w: w,
		h: h,
	}
}

func (b *BasicBoard) Initialize(cell Cell) {
	for i := 0; i < b.w * b.h; i++ {
		b.cells[i] = cell
	}
}

func (b *BasicBoard) Get(p Position) Cell {
	if p.X >= b.w {
		panic("position.x out of bounds")
	}
	if p.Y >= b.h {
		panic("position.y out of bounds")
	}
	return b.cells[p.Y*b.w + p.X]
}

func (b *BasicBoard) GetNeighbors(p Position) []Cell {
	neighbors := make([]Cell, 0, 8)
	bounds := b.Bounds()
	x, y := p.X, p.Y
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if !(i == x && j == y) && i >= 0 && j >= 0 && i <= bounds.X2Y2.X && j <= bounds.X2Y2.Y {
				neighbors = append(neighbors, b.Get(Position {i, j }))
			}
		}
	}
	return neighbors
}

func (b *BasicBoard) Set(p Position, cell Cell) {
	b.cells[p.Y*b.w + p.X] = cell
}

func (b *BasicBoard) Bounds() Rectangle {
	return Rectangle{ X1Y1: Origin, X2Y2: Position{ X: b.w-1, Y: b.h-1} }
}
