package main

import (
	"fmt"
	"github.com/ai-battleground/codemelee/client/tetris/redis"
	"math/rand"
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

func (a *Agent) Act(o tetris.Observation) string {
	if time.Now().Sub(a.lastAction) < 700*time.Millisecond {
		return ""
	}
	piece, surface := readBoard(o.Board)
	if piece == "" {
		return ""
	}
	const pos = 4
	var targetPos, highestRowIndex int
	adjustment := ""
	switch piece {
	case "I": // just put at lowest point
		for c := 0; c < len(surface); c++ {
			if surface[c] > highestRowIndex {
				targetPos = c
				highestRowIndex = surface[c]
			}
		}
		var candidates []int
		for c := range surface {
			if surface[c] == surface[targetPos] {
				candidates = append(candidates, c)
			}
		}
		targetPos = candidates[rand.Intn(len(candidates))]
	case "O": // lowest flat spot
		flatSpots := []int{}
		for c := 0; c < len(surface); c++ {
			if c+1 < len(surface) && surface[c+1] == surface[c] {
				flatSpots = append(flatSpots, c)
			}
		}
		if len(flatSpots) == 0 {
			// just slam it to the left
			targetPos = 0
		} else {
			for _, c := range flatSpots {
				if surface[c] > highestRowIndex {
					targetPos = c
					highestRowIndex = surface[c]
				}
			}
		}
	case "T":
		for c := 0; c < len(surface); c++ {
			if surface[c] > highestRowIndex {
				targetPos = c
				highestRowIndex = surface[c]
			}
		}
		var candidates []int
		for c := range surface {
			if surface[c] == surface[targetPos] {
				candidates = append(candidates, c)
			}
		}
		targetPos = candidates[rand.Intn(len(candidates))]
		switch {
		case targetPos == 0:
			adjustment = "/"
		case targetPos == len(surface)-1:
			adjustment = "\\>"
		default:
			leftRowIndex, rightRowIndex := surface[targetPos-1], surface[targetPos+1]
			switch {
			case leftRowIndex < rightRowIndex:
				adjustment = "/"
			case rightRowIndex < leftRowIndex:
				adjustment = "\\<"
			default: // they match
				if leftRowIndex != highestRowIndex {
					adjustment = "//"
				} else {
					// otherwise, they are all the same and we should leave the bottom flat
				}
				adjustment = adjustment + "<"
			}
		}
	case "L":
		candidates := map[int][]string{}
		for c := 0; c < len(surface); c++ {
			if c < len(surface)-2 &&
				surface[c]-1 == surface[c+1] &&
				surface[c]-1 == surface[c+2] {
				candidates[c] = append(candidates[c], "//")
			} else if c > 0 &&
				surface[c]-surface[c-1] > 1 {
				candidates[c] = append(candidates[c], "\\<")
			} else if c < len(surface)-1 &&
				surface[c] == surface[c+1] {
				candidates[c] = append(candidates[c], "/")
			} else if c < len(surface)-2 &&
				surface[c] == surface[c+1] &&
				surface[c] == surface[c+2] {
				candidates[c] = append(candidates[c], "")
			}
		}
		if len(candidates) > 0 {
			for c := range candidates {
				if surface[c] > highestRowIndex {
					highestRowIndex = surface[c]
				}
			}
			var finalists []int
			for c := range candidates {
				if surface[c] == highestRowIndex {
					finalists = append(finalists, c)
				}
			}
			fmt.Printf("L Finalists %v", finalists)
			targetPos = finalists[rand.Intn(len(finalists))]
			adjustment = candidates[targetPos][rand.Intn(len(candidates[targetPos]))]
		}
	case "J":
		candidates := map[int][]string{}
		for c := 0; c < len(surface); c++ {
			if c < len(surface)-2 &&
				surface[c] == surface[c+1] &&
				surface[c]+1 == surface[c+2] {
				candidates[c] = append(candidates[c], "//")
			} else if c > 0 &&
				surface[c]-surface[c+1] > 1 {
				candidates[c] = append(candidates[c], "/")
			} else if c < len(surface)-1 &&
				surface[c] == surface[c+1] {
				candidates[c] = append(candidates[c], "\\")
			} else if c < len(surface)-2 &&
				surface[c] == surface[c+1] &&
				surface[c] == surface[c+2] {
				candidates[c] = append(candidates[c], "")
			}
		}
		if len(candidates) > 0 {
			for c := range candidates {
				if surface[c] > highestRowIndex {
					highestRowIndex = surface[c]
				}
			}
			var finalists []int
			for c := range candidates {
				if surface[c] == highestRowIndex {
					finalists = append(finalists, c)
				}
			}
			fmt.Printf("J Finalists %v", finalists)
			targetPos = finalists[rand.Intn(len(finalists))]
			adjustment = candidates[targetPos][rand.Intn(len(candidates[targetPos]))]
		}

	}
	var action string
	if targetPos < pos {
		for i := pos; i > targetPos; i-- {
			action = action + "<"
		}
	} else {
		for i := pos; i < targetPos; i++ {
			action = action + ">"
		}
	}
	action = adjustment + action
	if a.Slow {
		action += "___"
	} else {
		action += "^"
	}
	a.lastAction = time.Now()
	return action
}

func readBoard(board [][]byte) (string, map[int]int) {
	surface := map[int]int{}
	var piece string
	for r := len(board) - 1; r >= 0; r-- { // bottom to top
		for c := 0; c < len(board[r]); c++ {
			if c > 3 && c < 7 && r < 5 { // probably the active piece
				if piece == "" && board[r][c] != ' ' {
					piece = string(board[r][c])
				}
				continue
			}
			if board[r][c] == ' ' {
				if r == len(board)-1 { // bottom row
					surface[c] = r + 1
				}
			} else {
				surface[c] = r
			}
		}
	}
	return piece, surface
}
