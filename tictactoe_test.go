package learning_examples

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestTicTacToeFixture(t *testing.T) {
	gunit.RunSequential(new(TicTacToeFixture), t)
}

type TicTacToeFixture struct {
	*gunit.Fixture
	normalBoard, giantBoard *Board
}

func (this *TicTacToeFixture) SetupBoards() {
	this.normalBoard = EmptyBoard(3, 3)
	this.giantBoard = EmptyBoard(13, 17)
}

func (this *TicTacToeFixture) TestBoardRows() {
	this.So(this.normalBoard.Height(), should.Equal, 3)
	this.So(this.giantBoard.Height(), should.Equal, 13)
}

func (this *TicTacToeFixture) TestBoardCols() {
	this.So(this.normalBoard.Width(), should.Equal, 3)
	this.So(this.giantBoard.Width(), should.Equal, 17)
}

func (this *TicTacToeFixture) TestEmptyBoardCellsAreEmpty() {
	for r := 0; r < this.normalBoard.Height(); r++ {
		for c := 0; c < this.normalBoard.Width(); c++ {
			this.So(this.normalBoard.Cell(r, c), should.Equal, CellStateEmpty)
		}
	}
}

func (this *TicTacToeFixture) TestNextMoveIsXWhenEmpty() {
	this.So(this.normalBoard.NextMove(), should.Equal, CellStateX)
}

func (this *TicTacToeFixture) TestMoveXWorks() {
	errNormal := this.normalBoard.X(1, 2)
	errGiant := this.giantBoard.X(5, 6)
	this.So(errNormal, should.BeNil)
	this.So(this.normalBoard.Cell(1, 2), should.Equal, CellStateX)
	this.So(errGiant, should.BeNil)
	this.So(this.giantBoard.Cell(5, 6), should.Equal, CellStateX)
}

func (this *TicTacToeFixture) TestMoveOWorks() {
	// X goes first
	this.normalBoard.X(0, 0)
	this.giantBoard.X(0, 0)
	errNormal := this.normalBoard.O(2, 1)
	errGiant := this.giantBoard.O(6, 5)
	this.So(errNormal, should.BeNil)
	this.So(this.normalBoard.Cell(2, 1), should.Equal, CellStateO)
	this.So(errGiant, should.BeNil)
	this.So(this.giantBoard.Cell(6, 5), should.Equal, CellStateO)
}

func (this *TicTacToeFixture) TestInvalidMoveGivesError() {
	normalBoardNonsense := [][2]int{{-1, 1}, {1, -1}, {0, 3}, {3, 0}, {1000, 1000}}
	giantBoardNonsense := [][2]int{{-1, 0}, {0, -1}, {17, 0}, {0, 13}, {1000, 1000}}
	for _, testdata := range normalBoardNonsense {
		err := this.normalBoard.X(testdata[0], testdata[1])
		this.So(err, should.BeError, ImpossibleMove(3, 3))
	}
	for _, testdata := range giantBoardNonsense {
		err := this.giantBoard.O(testdata[0], testdata[1])
		this.So(err, should.BeError, ImpossibleMove(13, 17))
	}
}

func (this *TicTacToeFixture) TestXGoesFirst() {
	err := this.normalBoard.O(0, 0)
	this.So(err, should.BeError, MoveOutOfTurn)
	err = this.normalBoard.X(2, 2)
	this.So(err, should.BeNil)
	err = this.normalBoard.O(0, 0)
	this.So(err, should.BeNil)
}

func (this *TicTacToeFixture) TestOneTurnAtATime() {
	err := this.normalBoard.X(0, 0)
	err = this.normalBoard.X(1, 1)
	this.So(err, should.BeError, MoveOutOfTurn)
	err = this.normalBoard.O(1, 1)
	err = this.normalBoard.O(0, 2)
	this.So(err, should.BeError, MoveOutOfTurn)
}

func (this *TicTacToeFixture) TestMoveMustNotBOccupied() {
	this.normalBoard.X(1, 1)
	err := this.normalBoard.O(1, 1)
	this.So(err, should.BeError, SpaceIsOccupied(1, 1))
	this.normalBoard.O(0, 0)
	err = this.normalBoard.X(0, 0)
	this.So(err, should.BeError, SpaceIsOccupied(0, 0))
}

func (this *TicTacToeFixture) TestXWins_3x3_Horizontally() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(1, 0)
	this.normalBoard.O(0, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.X(2, 0) // three in a row
	this.So(this.normalBoard.GameOutcome(), should.Equal, WonByX)
}

func (this *TicTacToeFixture) TestOWins_3x3_Horizontally() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(1, 0)
	this.normalBoard.O(1, 1)
	this.normalBoard.X(0, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.O(2, 1) // three in a row
	this.So(this.normalBoard.GameOutcome(), should.Equal, WonByO)
}

func (this *TicTacToeFixture) TestXWins_3x3_Vertically() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(2, 2)
	this.normalBoard.X(0, 1)
	this.normalBoard.O(2, 1)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.X(0, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, WonByX)
}

func (this *TicTacToeFixture) TestOWins_3x3_Vertically() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(2, 2)
	this.normalBoard.X(0, 1)
	this.normalBoard.O(2, 1)
	this.normalBoard.X(1, 1)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.O(2, 0)
	this.So(this.normalBoard.GameOutcome(), should.Equal, WonByO)
}

func (this *TicTacToeFixture) TestXWins_3x3_DescendingDiagonal() {
	this.normalBoard.X(0, 2)
	this.normalBoard.O(0, 0)
	this.normalBoard.X(1, 1)
	this.normalBoard.O(2, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.X(2, 0)
	this.So(this.normalBoard.GameOutcome(), should.Equal, WonByX)
}

func (this *TicTacToeFixture) TestOWins_3x3_AscendingDiagonal() {
	this.normalBoard.X(0, 2)
	this.normalBoard.O(0, 0)
	this.normalBoard.X(2, 0)
	this.normalBoard.O(1, 1)
	this.normalBoard.X(1, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.O(2, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, WonByO)
}

func (this *TicTacToeFixture) TestDraw_3x3() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(1, 1)
	this.normalBoard.X(2, 2)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(2, 1) // block forced from here on
	this.normalBoard.O(2, 0)
	this.normalBoard.X(0, 2)
	this.normalBoard.O(1, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Undetermined)
	this.normalBoard.X(1, 0)
	this.So(this.normalBoard.GameOutcome(), should.Equal, Draw)
}
