package learning_examples

type Game struct {
	*Board
}

func NewGame(x, y int) Game {
	return Game{
		Board: EmptyBoard(x, y),
	}
}
