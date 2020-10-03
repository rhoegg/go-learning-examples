package main

import (
	"flag"
	"fmt"
	tetris "github.com/ai-battleground/codemelee/client/tetris/redis"
	"strconv"
	"strings"
	"time"
)

const logTimeFormat = "2006-01-02 15:04:05.999"

func main() {
	config := parseArgs()
	driver, err := tetris.NewDriver(config.RedisUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	// for each game in args
	for _, g := range config.Games {
		fmt.Printf("Playing game %s\n", g)
		agent := NewAgent(config)
		//  if game is not started, start it
		o := driver.Observe(g)
		if o.Error != nil {
			fmt.Printf("Error starting game %s %v\n", g, o.Error)
			return
		}
		state := o.State
		fmt.Printf("Found game %s (%s)\n", g, o.State)
		if o.State == "Not Started" {
			fmt.Printf("Starting %s\n", g)
			driver.Act(g, "Start")
		}

		//  loop until game is over
		ticker := time.NewTicker(500 * time.Millisecond)
		for state != "Game Over" {
			select {
			case <-ticker.C:
				o := driver.Observe(g)
				display(o)
				state = o.State
				//    determine reward for last action and learn
				//    compute action
				action := agent.Act(o)
				//    act
				if action != "" {
					driver.Act(g, action)
				}
			}
		}
	}
}

type Config struct {
	RunSlowly bool
	RedisUrl  string
	Games     []string
}

func parseArgs() Config {
	slow := flag.Bool("slow", false, "play slowly")
	redisUrl := flag.String("redis-url", "redis://localhost", "URL of redis")

	flag.Parse()
	return Config{
		RunSlowly: *slow,
		RedisUrl:  *redisUrl,
		Games:     flag.Args(),
	}
}

func display(o tetris.Observation) {
	if len(o.Board) > 0 {
		fmt.Println(strings.Repeat("=", len(o.Board[0])+2+13))
		for i, r := range o.Board {
			fmt.Print("=")
			fmt.Print(string(r))
			switch i {
			case 0, 3:
				fmt.Print(strings.Repeat("=", 13))
			case 1:
				fmt.Print(displayScore(o))
			case 2:
				fmt.Print(displayLines(o))
			}
			fmt.Println("=")
		}
		fmt.Println(strings.Repeat("=", len(o.Board[0])+2))
	}
}

func displayScore(o tetris.Observation) string {
	score := strconv.Itoa(o.Score)
	pad := strings.Repeat(" ", 4-len(score))
	return "= Score: " + score + pad
}

func displayLines(o tetris.Observation) string {
	lines := strconv.Itoa(o.Lines)
	pad := strings.Repeat(" ", 4-len(lines))
	return "= Lines: " + lines + pad
}
