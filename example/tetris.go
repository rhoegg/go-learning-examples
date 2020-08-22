package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type GameState struct {
	Status   string `json:"status"`
	Score    int    `json:"score"`
	Bot      string `json:"bot"`
	Next     string `json:"next"`
	RawBoard string `json:"board"`
}

type Action func() (GameState, error)

var gameState GameState
var quietPeriod bool

func (gs GameState) BoardLines() [][]byte {
	lines := strings.Split(gs.RawBoard, "\n")
	var result [][]byte
	for _, line := range lines {
		result = append(result, []byte(line))
	}
	return result
}

var apiEndpoint string

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Usage: %s <API endpoint>", os.Args[0])
		return
	}

	apiEndpoint = os.Args[1]

	var err error
	gameState, err = getGameState()
	if err != nil {
		fmt.Printf("Error: %v")
		return
	}
	t := time.NewTicker(300 * time.Millisecond)
	for {
		select {
		case <-t.C:
			err := update()
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
		}
	}
}

func update() error {
	switch gameState.Status {
	case "Not Started":
		fmt.Println("Game is not started. Starting...")
		var err error
		gameState, err = act("S")
		return err
	case "Running":
		action := getGameState
		if pieceOnTopLine() {
			if !quietPeriod {
				quietPeriod = true
				time.AfterFunc(1*time.Second, func() {
					quietPeriod = false
				})

				action = moveRandom
			}
		}
		// if time to act,
		//	choose action and send it
		//    and return game state
		var err error
		gameState, err = action()
		return err
	default:
		return fmt.Errorf("unsupported status %s", gameState.Status)
	}
}

func getGameState() (GameState, error) {
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return gameState, fmt.Errorf("API endpoint failure %v", err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	bodyString := string(bodyBytes)

	err = json.Unmarshal(bodyBytes, &gameState)
	if err != nil {
		fmt.Printf("ERROR parsing :%s\n", bodyString)
		return gameState, err
	}
	fmt.Printf("DEBUG: board\n%s\n", gameState.RawBoard)
	return gameState, err
}

func moveRandom() (GameState, error) {
	moves := rand.Intn(5)
	rotations := rand.Intn(3)

	directionRight := rand.Float32() > 0.5
	clockwise := rand.Float32() > 0.5

	move := "<"
	if directionRight {
		move = ">"
	}
	rotation := "\\"
	if clockwise {
		rotation = "/"
	}
	actions := strings.Repeat(rotation, rotations) + strings.Repeat(move, moves)
	fmt.Printf("Moving random %s\n", actions)
	return act(actions)
}

func act(action string) (GameState, error) {
	resp, err := http.Post(apiEndpoint+"/action", "text/plain", strings.NewReader(action))
	if err != nil {
		return gameState, err
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	bodyString := string(bodyBytes)
	fmt.Println("DEBUG: \n", bodyString)

	err = json.Unmarshal(bodyBytes, &gameState)
	return gameState, err
}

func pieceOnTopLine() bool {
	topLine := gameState.BoardLines()[0]
	for _, b := range topLine {
		if b != ' ' {
			return true
		}
	}
	return false
}
