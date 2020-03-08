package learning_examples

import "fmt"

const NO_PROBLEM RuleViolation = RuleViolation("NO_PROBLEM")

type RuleViolation string

func (v RuleViolation) Error() string {
	return string(v)
}

func ImpossibleMove(rows, cols int) RuleViolation {
	return RuleViolation(fmt.Sprintf("Move not possible: Board is %dx%d", rows, cols))
}

func (b Board) CheckPossibleMove(x, y int) (RuleViolation, bool) {
	if x < 0 || y < 0 {
		return ImpossibleMove(b.Rows(), b.Cols()), false
	}
	if x >= b.Rows() || y >= b.Cols() {
		return ImpossibleMove(b.Rows(), b.Cols()), false
	}
	return NO_PROBLEM, true
}
