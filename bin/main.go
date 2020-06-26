package main

import (
	"bufio"
	"fmt"
	learning_examples "github.com/rhoegg/learning-examples-go"
	"os"
)

func main() {
	b := learning_examples.EmptyBoard(3, 3)

	scanner := bufio.NewScanner(os.Stdin)
	playerX := true
	for {
		fmt.Printf("%v", b)
		fmt.Print("\nMove: ")
		scanner.Scan()
		x, y, err := parseMove(scanner.Text())
		if err != nil {
			fmt.Printf("Error %v", err)
			return
		}
		if playerX {
			err = b.X(x, y)
		} else {
			err = b.O(x, y)
		}
		if err != nil {
			fmt.Printf("Error %v", err)
			return
		}
		playerX = !playerX
	}
}

func parseMove(text string) (int, int, error) {
	switch text {
	case "1":
		return 0, 0, nil
	case "2":
		return 1, 0, nil
	case "3":
		return 2, 0, nil
	case "4":
		return 0, 1, nil
	case "5":
		return 1, 1, nil
	case "6":
		return 2, 1, nil
	case "7":
		return 0, 2, nil
	case "8":
		return 1, 2, nil
	case "9":
		return 2, 2, nil
	default:
		return -1, -1, fmt.Errorf("Unexpected move %v: Expected 1-9", text)
	}
}
