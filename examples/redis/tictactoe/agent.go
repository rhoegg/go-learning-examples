package main

import (
	"fmt"
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

func (a *Agent) Act(o Observation) string {
	if !o.MyTurn {
		return ""
	}
	actions := ""
	for i := range o.Boards {
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
