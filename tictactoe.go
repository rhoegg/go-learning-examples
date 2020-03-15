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
	Draw
	WonByX
	WonByO
)

func (g GameOutcome) String() string {
	switch g {
	case Draw:
		return "Draw"
	case WonByO:
		return "O wins"
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
	} else if b.score(CellStateO) {
		return WonByO
	}
	for _, row := range b.rows() {
		for _, cell := range row {
			if cell == CellStateEmpty { // at least one available cell remains
				return Undetermined
			}
		}
	}
	return Draw
}

func (b Board) score(state CellState) bool {
	for _, r := range b.rows() {
		if contiguousCells(state, 3, r) {
			return true
		}
	}
	for _, c := range b.columns() {
		if contiguousCells(state, 3, c) {
			return true
		}
	}
	for _, d := range b.descents() {
		if contiguousCells(state, 3, d) {
			return true
		}
	}
	for _, a := range b.ascents() {
		if contiguousCells(state, 3, a) {
			return true
		}
	}
	return false
}

func (b Board) rows() [][]CellState {
	return b.cells
}

func (b Board) columns() [][]CellState {
	var result [][]CellState
	for x := 0; x < b.Width(); x++ {
		var cells []CellState
		for y := 0; y < b.Height(); y++ {
			cells = append(cells, b.Cell(x, y))
		}
		result = append(result, cells)
	}
	return result
}

func (b Board) descents() [][]CellState {
	var result [][]CellState
	for x := 0; x < b.Width(); x++ {
		for y := b.Height(); y >= x; y-- {
			// diagonal down to x
			var descent []CellState
			for i, j := x, y; i < y; i, j = i+1, j-1 {
				descent = append(descent, b.Cell(i, j-1))
			}
			result = append(result, descent)
		}
	}
	return result
}

func (b Board) ascents() [][]CellState {
	var result [][]CellState
	for x := 0; x < b.Width(); x++ {
		for y := 0; y < b.Height(); y++ {
			var ascent []CellState
			for i, j := x, y; i < b.Width() && j < b.Height(); i, j = i+1, j+1 {
				ascent = append(ascent, b.Cell(i, j))
			}
			result = append(result, ascent)
		}
	}
	return result
}

func contiguousCells(state CellState, needed int, cells []CellState) bool {
	if len(cells) < needed {
		return false
	}
	if needed < 1 {
		return true
	}
	if cells[0] != state {
		return contiguousCells(state, needed, cells[1:])
	}
	for i := 0; i < needed; i++ {
		if cells[i] != state {
			return false // TODO: maybe recurse for the remainder of the cells here?
		}
	}
	return true
}
