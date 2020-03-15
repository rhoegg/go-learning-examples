package learning_examples

import "fmt"

const (
	NoProblem     = RuleViolation("OK")
	MoveOutOfTurn = RuleViolation("Moved out of turn")
)

type RuleViolation string
type Rule func() (RuleViolation, bool)

func (v RuleViolation) Error() string {
	return string(v)
}

func ImpossibleMove(rows, cols int) RuleViolation {
	return RuleViolation(fmt.Sprintf("Move not possible: Board is %dx%d", rows, cols))
}

func SpaceIsOccupied(x, y int) RuleViolation {
	return RuleViolation(fmt.Sprintf("Space is already taken (%d, %d)", x, y))
}

func (b Board) CheckMoveIsPossible(x, y int) Rule {
	return func() (RuleViolation, bool) {
		if x < 0 || y < 0 {
			return ImpossibleMove(b.Rows(), b.Cols()), false
		}
		if y >= b.Rows() || x >= b.Cols() {
			return ImpossibleMove(b.Rows(), b.Cols()), false
		}
		return NoProblem, true
	}
}

func (b Board) CheckXGoesFirst(state CellState) Rule {
	return func() (RuleViolation, bool) {
		if !b.inProgress && state == CellStateO {
			return MoveOutOfTurn, false
		}
		return NoProblem, true
	}
}

func (b Board) CheckTakingTurns(state CellState) Rule {
	var turnIsO bool = false
	for i := 0; i < b.Cols(); i++ {
		for j := 0; j < b.Rows(); j++ {
			if b.Cell(i, j) != CellStateEmpty {
				turnIsO = !turnIsO
			}
		}
	}
	return func() (RuleViolation, bool) {
		if turnIsO {
			if state == CellStateX {
				return MoveOutOfTurn, false
			}
		} else {
			if state == CellStateO {
				return MoveOutOfTurn, false
			}
		}
		return NoProblem, true
	}
}

func (b Board) CheckUnoccupied(x, y int) Rule {
	return func() (RuleViolation, bool) {
		if b.Cell(x, y) != CellStateEmpty {
			return SpaceIsOccupied(x, y), false
		}
		return NoProblem, true
	}
}
