package learning_examples

type CellState int

const (
	CellStateEmpty CellState = iota
	CellStateX
	CellStateO
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

func (b Board) Rows() int {
	return len(b.cells)
}

func (b Board) Cols() int {
	return len(b.cells[0])
}

func (b *Board) Cell(x, y int) CellState {
	return b.cells[x][y]
}

func (b *Board) X(x, y int) error {
	return b.move(x, y, CellStateX)
}

func (b *Board) O(x, y int) error {
	return b.move(x, y, CellStateO)
}

func (b *Board) move(x, y int, move CellState) error {
	for _, rule := range []Rule{b.CheckMoveIsPossible(x, y), b.CheckXGoesFirst(move)} {
		v, ok := rule()
		if !ok {
			return v
		}
	}
	b.cells[x][y] = move
	b.inProgress = true
	return nil
}
