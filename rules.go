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

func (b Board) CheckMoveIsPossible(x, y int) Rule {
	return func() (RuleViolation, bool) {
		if x < 0 || y < 0 {
			return ImpossibleMove(b.Rows(), b.Cols()), false
		}
		if x >= b.Rows() || y >= b.Cols() {
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
