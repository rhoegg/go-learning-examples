package main

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"time"
)

const logTimeFormat = "2006-01-02 15:04:05.999"

func main() {
	config := parseArgs()
	driver, err := NewDriver(config.RedisUrl)
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
		time.Sleep(300 * time.Millisecond)
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

func waitForMatch(driver Driver, bot, challenge string) string {
	// loop confirm
	fmt.Printf("Bot %s waiting for match...\n", bot)
	ticker := time.NewTicker(3 * time.Second)
	timeout := time.NewTimer(2 * time.Minute)
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
			fmt.Printf("Giving up")
			return ""
		}
	}

}
func display(o Observation) {
	rows := [3]string{}
	X, O := 0, 0
	for _, b := range o.Boards {
		for _, s := range b {
			switch s {
			case 'X':
				X++
			case 'Y':
				O++
			}
		}
		if math.Abs(float64(X-O)) > 1 {
			fmt.Println("Out of whack")
		}
		rows[0] += " " + string(b[0:3])
		rows[1] += " " + string(b[3:6])
		rows[2] += " " + string(b[6:9])
	}
	fmt.Println(strings.Join(rows[:], "\n"))
	fmt.Printf("%s: %d / %s: %d\n", o.Bot, o.Score, o.Opponent, o.OpponentScore)
}
