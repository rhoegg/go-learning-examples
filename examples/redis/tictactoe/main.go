package main

import (
	"flag"
	"fmt"
	tictactoe "github.com/ai-battleground/codemelee/client/tictactoe/redis"
	"strings"
	"time"
)

const logTimeFormat = "2006-01-02 15:04:05.000"

func main() {
	config := parseArgs()
	driver, err := tictactoe.NewDriver(config.RedisUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	// challenge
	challenge, err := driver.Challenge(config.Bot, config.Boards, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Challenge created %s\n", challenge)
	game := waitForMatch(driver, config.Bot, challenge)
	if game == "" {
		fmt.Println("Giving up")
		return
	}
	// loop observe/act
	agent := NewAgent(config)
	for {
		o := driver.Observe(config.Bot, game)
		display(o)
		actions := agent.Act(o)
		if actions != "" {
			err = driver.Act(config.Bot, game, actions)
			if err != nil {
				fmt.Printf("Error acting: %v\n", err)
			}
		}
		if o.State == "Done" {
			break
		}
		time.Sleep(40 * time.Millisecond)
	}
}

type Config struct {
	RunSlowly bool
	RedisUrl  string
	Bot       string
	Boards    int
}

func parseArgs() Config {
	slow := flag.Bool("slow", false, "play slowly")
	redisUrl := flag.String("redis-url", "redis://localhost", "URL of redis")
	boards := flag.Int("boards", 1, "Number of tic tac toe boards to play at once")

	flag.Parse()
	return Config{
		RunSlowly: *slow,
		RedisUrl:  *redisUrl,
		Boards:    *boards,
		Bot:       strings.Join(flag.Args(), " "),
	}
}

func waitForMatch(driver tictactoe.Driver, bot, challenge string) string {
	// loop confirm
	fmt.Printf("%s Bot %s waiting for match...\n", time.Now().Format(logTimeFormat), bot)
	ticker := time.NewTicker(5 * time.Second)
	timeout := time.NewTimer(1 * time.Minute)
	var game string
	for {
		select {
		case <-ticker.C:
			game = driver.Confirm(bot, challenge)
			if game != "" {
				fmt.Printf("%s beginning match %s\n", time.Now().Format(logTimeFormat), game)
				return game
			}
			fmt.Printf("%s no match found for %s\n", time.Now().Format(logTimeFormat), challenge)
		case <-timeout.C:
			fmt.Printf("Giving up\n")
			return ""
		}
	}

}
func display(o tictactoe.Observation) {
	rows := [3]string{}
	for _, b := range o.Boards {
		if len(b) == 9 {
			rows[0] += " | " + string(b[0:3])
			rows[1] += " | " + string(b[3:6])
			rows[2] += " | " + string(b[6:9])
		}
	}
	fmt.Println(strings.Join(rows[:], "\n"))
	fmt.Printf("%s: %d / %s: %d\n", o.Bot, o.Score, o.Opponent, o.OpponentScore)
}
