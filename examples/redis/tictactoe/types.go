package main

import "time"

type Observation struct {
	Boards        [][]byte
	MyTurn        bool
	State         string
	Error         error
	Score         int
	OpponentScore int
	Bot           string
	Opponent      string
	Round         int
	MoveTimeout   time.Duration
	MatchTimeout  time.Duration
}

type GameTime struct {
	Ticks int
	Time  time.Time
}
