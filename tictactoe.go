package learning_examples

type CellState int

const (
	CellStateEmpty CellState = iota
	CellStateX
	CellStateO
)

type Board [][]CellState

func EmptyBoard(rows, cols int) Board {
	var board Board
	for r := 0; r < rows; r++ {
		var row []CellState
		for c := 0; c < cols; c++ {
			row = append(row, CellStateEmpty)
		}
		board = append(board, row)
	}
	return board
}

func (b Board) Rows() int {
	return len(b)
}

func (b Board) Cols() int {
	return len(b[0])
}

func (b Board) Cell(x, y int) CellState {
	return b[x][y]
}

func (b Board) MoveX(x, y int) error {
	v, ok := b.CheckPossibleMove(x, y)
	if !ok {
		return v
	}
	b[x][y] = CellStateX
	return nil
}

func (b Board) MoveO(x, y int) error {
	v, ok := b.CheckPossibleMove(x, y)
	if !ok {
		return v
	}
	b[x][y] = CellStateO
	return nil
}
