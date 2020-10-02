package main

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestTetrisFixture(t *testing.T) {
	gunit.Run(new(TetrisFixture), t)
}

type TetrisFixture struct {
	*gunit.Fixture
}

func (this *TetrisFixture) TestDetectLeftmostLowestPoint() {
	firstBlockedSpace := map[int]int{0: 7, 1: 20, 2: 18, 3: 18, 4: 18, 5: 18, 6: 0, 7: 0, 8: 19, 9: 19}
	const pos = 4
	var targetPos, highestRowIndex int
	for c := 0; c < len(firstBlockedSpace); c++ {
		if firstBlockedSpace[c] > highestRowIndex {
			targetPos = c
			highestRowIndex = firstBlockedSpace[c]
		}
	}
	this.So(targetPos, should.Equal, 1)
}
