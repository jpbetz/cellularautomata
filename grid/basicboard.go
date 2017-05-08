package grid

import (
	"fmt"
	"log"
)

type BasicBoard struct {
	Cells []Cell
	W, H  int
}

func NewBasicBoard(w, h int) *BasicBoard {
	return &BasicBoard{
		Cells: make([]Cell, w*h),
		W:     w,
		H:     h,
	}
}

func (b *BasicBoard) Initialize(cell Cell) {
	for i := 0; i < b.W*b.H; i++ {
		b.Cells[i] = cell
	}
}

func (b *BasicBoard) Get(p Position) Cell {
	if p.X >= b.W {
		panic(fmt.Sprintf("position.x out of bounds: %d >= %d", p.X, b.W))
	}
	if p.Y >= b.H {
		panic("position.y out of bounds")
	}
	return b.Cells[p.Y*b.W+p.X]
}

func (b *BasicBoard) GetNeighborPositions(p Position) []Position {
	neighbors := make([]Position, 0, 8)
	bounds := b.Bounds()
	x, y := p.X, p.Y
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if !(i == x && j == y) && i >= 0 && j >= 0 && i <= bounds.Corner2.X && j <= bounds.Corner2.Y {
				neighbors = append(neighbors, Position{i, j})
			}
		}
	}
	log.Printf("GetNeighborPositions %v: %v", p, neighbors)
	return neighbors
}

func (b *BasicBoard) GetNeighbors(p Position) []Cell {
	neighbors := make([]Cell, 0, 8)
	for _, neighborPosition := range b.GetNeighborPositions(p) {
		neighbors = append(neighbors, b.Get(neighborPosition))
	}
	return neighbors
}

func (b *BasicBoard) Set(p Position, cell Cell) {
	b.Cells[p.Y*b.W+p.X] = cell
}

func (b *BasicBoard) Bounds() Rectangle {
	return Rectangle{Corner1: Origin, Corner2: Position{X: b.W - 1, Y: b.H - 1}}
}
