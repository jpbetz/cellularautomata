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

func (b *BasicBoard) Set(p Position, cell Cell) {
	b.cells[p.Y*b.w + p.X] = cell
}

func (b *BasicBoard) Bounds() Rectangle {
	return Rectangle{ X1Y1: Origin, X2Y2: Position{ X: b.w-1, Y: b.h-1} }
}
