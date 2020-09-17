package main

import "strings"

type Action func() (GameState, error)

type GameState struct {
	Status   string `json:"status"`
	Score    int    `json:"score"`
	Bot      string `json:"bot"`
	Next     string `json:"next"`
	RawBoard string `json:"board"`
}

func (gs GameState) BoardLines() [][]byte {
	lines := strings.Split(gs.RawBoard, "\n")
	var result [][]byte
	for _, line := range lines {
		result = append(result, []byte(line))
	}
	return result
}
