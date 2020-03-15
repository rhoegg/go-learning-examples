package learning_examples

type CellState int

const (
	CellStateEmpty CellState = iota
	CellStateX
	CellStateO
)

func (c CellState) String() string {
	switch c {
	case CellStateX:
		return "X"
	case CellStateO:
		return "O"
	case CellStateEmpty:
		return " "
	}
	return ""
}

type GameOutcome int

const (
	Undetermined GameOutcome = iota
	WonByX
)

func (g GameOutcome) String() string {
	switch g {
	case WonByX:
		return "X wins"
	case Undetermined:
		return "Game in progress"
	}
	return "Invalid game state"
}

func (b Board) GameOutcome() GameOutcome {
	if b.score(CellStateX) {
		return WonByX
	}
	return Undetermined
}

func (b Board) score(state CellState) bool {
	for _, row := range b.cells {
		if contiguousCells(state, 3, row...) {
			return true
		}
	}
	return false
}

func contiguousCells(state CellState, needed int, cells ...CellState) bool {
	if len(cells) < needed {
		return false
	}
	if needed < 1 {
		return true
	}
	if cells[0] != state {
		return contiguousCells(state, needed, cells[1:]...)
	}
	for i := 0; i < needed; i++ {
		if cells[i] != state {
			return false // TODO: maybe recurse for the remainder of the row here?
		}
	}
	return true
}
