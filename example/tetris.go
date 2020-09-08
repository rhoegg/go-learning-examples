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
		fmt.Printf("Error: %v", err)
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
		if !quietPeriod {
			if pieceOnTopLine() || pieceOnLine(1) {
				quietPeriod = true
				time.AfterFunc(500*time.Millisecond, func() {
					quietPeriod = false
				})

				action = moveLowestThatFits
			}
		}
		// if time to act,
		//	choose action and send it
		//    and return game state
		var err error
		gameState, err = action()
		return err
	case "Game Over":
		return fmt.Errorf("game over with score %d", gameState.Score)
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

func moveLowestThatFits() (GameState, error) {
	piece, surface := readBoard()
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
		time.Sleep(500 * time.Millisecond)
		for i := pos; i > targetPos; i-- {
			action = action + "<"
		}
	} else {
		time.Sleep(500 * time.Millisecond)
		for i := pos; i < targetPos; i++ {
			action = action + ">"
		}
	}
	action = adjustment + action
	return act(action + "^")
}

func readBoard() (string, map[int]int) {
	surface := map[int]int{}
	var piece string
	for r := len(gameState.BoardLines()) - 1; r >= 0; r-- { // bottom to top
		for c := 0; c < len(gameState.BoardLines()[r]); c++ {
			if c > 3 && c < 7 && r < 4 { // probably the active piece
				if piece == "" && gameState.BoardLines()[r][c] != ' ' {
					piece = string(gameState.BoardLines()[r][c])
				}
				continue
			}
			if gameState.BoardLines()[r][c] == ' ' {
				if r == len(gameState.BoardLines())-1 { // bottom row
					surface[c] = r + 1
				}
			} else {
				surface[c] = r
			}
		}
	}
	return piece, surface
}

var oMove = 0

func moveO() (GameState, error) {
	moves := []string{"<<<<", "<<", "", ">>", ">>>>"}
	move := moves[oMove]
	oMove = (oMove + 1) % 5
	return act(move)
}

func drop() (GameState, error) {
	return act("_")
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

	err = json.Unmarshal(bodyBytes, &gameState)
	return gameState, err
}

func pieceOnTopLine() bool {
	return pieceOnLine(0)
}
func pieceOnLine(i int) bool {
	topLine := gameState.BoardLines()[i]
	for _, b := range topLine {
		if b != ' ' {
			return true
		}
	}
	return false
}
