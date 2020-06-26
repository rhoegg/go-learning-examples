package learning_examples

import (
	"strings"
)

type Board struct {
	cells      [][]CellState
	inProgress bool
}

func EmptyBoard(rows, cols int) *Board {
	var board Board
	for r := 0; r < rows; r++ {
		var row []CellState
		for c := 0; c < cols; c++ {
			row = append(row, CellStateEmpty)
		}
		board.cells = append(board.cells, row)
	}
	return &board
}

func (b Board) Height() int {
	return len(b.cells)
}

func (b Board) Width() int {
	return len(b.cells[0])
}

func (b Board) Cell(x, y int) CellState {
	return b.cells[y][x]
}

func (b Board) String() string {
	var lines []string
	for _, row := range b.cells {
		var rowstrings []string
		for c := range row {
			rowstrings = append(rowstrings, string(c))
		}
		lines = append(lines, strings.Join(rowstrings, " | "))
	}
	return strings.Join(lines, "\n---------\n")
}

func (b *Board) X(x, y int) error {
	return b.move(x, y, CellStateX)
}

func (b *Board) O(x, y int) error {
	return b.move(x, y, CellStateO)
}

func (b *Board) move(x, y int, move CellState) error {
	rules := []Rule{
		b.CheckMoveIsPossible(x, y),
		b.CheckXGoesFirst(move),
		b.CheckTakingTurns(move),
		b.CheckUnoccupied(x, y)}

	for _, rule := range rules {
		v, ok := rule()
		if !ok {
			return v
		}
	}
	b.cells[y][x] = move
	b.inProgress = true
	return nil
}
