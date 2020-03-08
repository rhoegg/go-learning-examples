package learning_examples

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestTicTacToeFixture(t *testing.T) {
	gunit.Run(new(TicTacToeFixture), t)
}

type TicTacToeFixture struct {
	*gunit.Fixture
	normalBoard, giantBoard Board
}

func (this *TicTacToeFixture) SetupBoards() {
	this.normalBoard = EmptyBoard(3, 3)
	this.giantBoard = EmptyBoard(13, 17)
}

func (this *TicTacToeFixture) TestBoardRows() {
	this.So(this.normalBoard.Rows(), should.Equal, 3)
	this.So(this.giantBoard.Rows(), should.Equal, 13)
}

func (this *TicTacToeFixture) TestBoardCols() {
	this.So(this.normalBoard.Cols(), should.Equal, 3)
	this.So(this.giantBoard.Cols(), should.Equal, 17)
}

func (this *TicTacToeFixture) TestEmptyBoardCellsAreEmpty() {
	for r := 0; r < this.normalBoard.Rows(); r++ {
		for c := 0; c < this.normalBoard.Cols(); c++ {
			this.So(this.normalBoard.Cell(r, c), should.Equal, CellStateEmpty)
		}
	}
}

func (this *TicTacToeFixture) TestMoveXWorks() {
	errNormal := this.normalBoard.MoveX(1, 2)
	errGiant := this.giantBoard.MoveX(5, 6)
	this.So(errNormal, should.BeNil)
	this.So(this.normalBoard.Cell(1, 2), should.Equal, CellStateX)
	this.So(errGiant, should.BeNil)
	this.So(this.giantBoard.Cell(5, 6), should.Equal, CellStateX)
}

func (this *TicTacToeFixture) TestMoveOWorks() {
	errNormal := this.normalBoard.MoveO(2, 1)
	errGiant := this.giantBoard.MoveO(6, 5)
	this.So(errNormal, should.BeNil)
	this.So(this.normalBoard.Cell(2, 1), should.Equal, CellStateO)
	this.So(errGiant, should.BeNil)
	this.So(this.giantBoard.Cell(6, 5), should.Equal, CellStateO)
}

func (this *TicTacToeFixture) TestInvalidMoveGivesError() {
	normalBoardNonsense := [][2]int{{-1, 1}, {1, -1}, {0, 3}, {3, 0}, {1000, 1000}}
	giantBoardNonsense := [][2]int{{-1, 0}, {0, -1}, {0, 17}, {13, 0}, {1000, 1000}}
	for _, testdata := range normalBoardNonsense {
		err := this.normalBoard.MoveX(testdata[0], testdata[1])
		this.So(err, should.BeError, ImpossibleMove(3, 3))
	}
	for _, testdata := range giantBoardNonsense {
		err := this.giantBoard.MoveO(testdata[0], testdata[1])
		this.So(err, should.BeError, ImpossibleMove(13, 17))
	}
}
