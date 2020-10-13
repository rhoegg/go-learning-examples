package main

import (
	"fmt"
	tictactoe "github.com/ai-battleground/codemelee/client/tictactoe/redis"
	"math/rand"
	"strings"
	"time"
)

type Agent struct {
	Slow       bool
	lastAction time.Time
}

func NewAgent(config Config) *Agent {
	return &Agent{
		Slow: config.RunSlowly,
	}
}

func (a *Agent) Act(o tictactoe.Observation) string {
	if !o.MyTurn {
		//fmt.Printf("Not my turn\n")
		return ""
	}
	var activeBoards []int
	for i, _ := range o.Boards {
		// don't put it in there if someone won
		if !gameOver(o.Boards[i]) {
			activeBoards = append(activeBoards, i)
		}
	}
	var maxSpaces int
	for _, i := range activeBoards {
		spaces := strings.Count(string(o.Boards[i]), " ")
		if spaces > maxSpaces {
			maxSpaces = spaces
		}
	}
	var actionBoards []int
	for _, i := range activeBoards {
		b := o.Boards[i]
		if strings.Count(string(b), " ") == maxSpaces {
			actionBoards = append(actionBoards, i)
		}
	}
	actions := ""
	for _, i := range actionBoards {
		var available []int
		for j, space := range o.Boards[i] {
			if space == ' ' {
				available = append(available, j)
			}
		}
		// choose a random available space
		if len(available) > 0 {
			chosen := rand.Intn(len(available))
			actions += fmt.Sprintf("%d%d", i, available[chosen])
		}
	}
	var tmp []string
	for _, b := range o.Boards {
		tmp = append(tmp, string(b))
	}
	//fmt.Printf("%s My turn: %s: %s\n", time.Now().Format(logTimeFormat), strings.Join(tmp, " "), actions)
	return actions
}

var wins = [][]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

func gameOver(board []byte) bool {
	for _, win := range wins {
		if board[win[0]] == ' ' || board[win[1]] == ' ' || board[win[2]] == ' ' {
			continue
		}
		if board[win[0]] == board[win[1]] &&
			board[win[0]] == board[win[2]] {
			return true
		}
	}
	return false
}
