package main

import (
	"fmt"
	"github.com/mediocregopher/radix/v3"
	"github.com/rs/xid"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Driver struct {
	pool               *radix.Pool
	token              map[string]string
	started            map[string]GameTime
	lastActionResultId radix.StreamEntryID
}

func NewDriver(redisUrl string) (Driver, error) {
	d := Driver{
		token:   make(map[string]string),
		started: make(map[string]GameTime),
	}
	u, err := url.Parse(redisUrl)
	if err != nil {
		return d, fmt.Errorf("Invalid redis URL %s: %w", u, err)
	}
	if u.Scheme != "redis" {
		return d, fmt.Errorf("Invalid redis URL: scheme %s not recognized", u.Scheme)
	}
	if len(u.Port()) == 0 {
		u.Host = u.Hostname() + ":6379"
	}
	d.pool, err = radix.NewPool("tcp", u.Host, 3)
	return d, err
}

func (d Driver) Challenge(bot string, boards int, opponent string) (string, error) {
	if bot == "" {
		return "", fmt.Errorf("bot name is required")
	}
	if boards < 1 {
		return "", fmt.Errorf("boards must be positive")
	}
	// generate challenge id
	id := xid.New()
	challengeId := fmt.Sprintf("%s:%s", bot, id.String())
	// create challenge
	key := fmt.Sprintf("challenge:tictactoe:%s", challengeId)
	hsetArgs := []string{key, "boards", strconv.Itoa(boards)}
	if opponent != "" {
		hsetArgs = append(hsetArgs, "opponent", opponent)
	}
	err := d.pool.Do(radix.Cmd(nil, "HSET", hsetArgs...))
	if err != nil {
		return "", err
	}
	return challengeId, nil
}

func (d Driver) Confirm(bot, challenge string) string {
	var game string
	// look for opportunity
	err := d.pool.Do(radix.Cmd(&game, "GET", fmt.Sprintf("opportunity:tictactoe:%s", challenge)))
	if err != nil || game == "" {
		return ""
	}
	// if it's there, put the game ID under match key and generate a token
	token := xid.New().String()
	err = d.pool.Do(radix.Cmd(nil, "HSET", fmt.Sprintf("challenge:tictactoe:%s", challenge),
		"match", game,
		"token", token))
	if err != nil {
		return ""
	}
	d.token[fmt.Sprintf("%s:%s", bot, game)] = token
	return game
}

func (d Driver) Observe(bot, game string) Observation {
	key := fmt.Sprintf("observe:tictactoe:%s", game)
	var exists bool
	fact := make(map[string]string)
	err := d.pool.Do(radix.Cmd(&exists, "EXISTS", key))
	if err != nil {
		return Observation{Error: err}
	}
	if !exists {
		return Observation{Error: fmt.Errorf("no game %s", game)}
	}
	err = d.pool.Do(radix.Cmd(&fact, "HGETALL", key))
	if err != nil {
		return Observation{Error: err}
	}
	o := Observation{
		State: fact["state"],
		Bot:   bot,
	}
	for k, v := range fact {
		switch {
		case strings.HasPrefix(k, "score:"):
			if score, err := strconv.Atoi(v); err == nil {
				scoreBot := strings.SplitN(k, ":", 2)[1]
				if scoreBot == bot {
					o.Score = score
				} else {
					o.OpponentScore = score
					o.Opponent = scoreBot
				}
			}
		case k == "boards":
			fmt.Printf("[driver] boards: %s\n", v)
			for _, b := range strings.Split(v, "\n") {
				o.Boards = append(o.Boards, []byte(b))
			}
		case k == "round":
			if r, err := strconv.Atoi(v); err == nil {
				o.Round = r
			}
		case k == "turn":
			o.MyTurn = v == bot
		case k == "turnExpires":
			if ms, err := strconv.Atoi(v); err == nil {
				o.MoveTimeout = time.Duration(ms) * time.Millisecond
			}
		case k == "matchExpires":
			if ms, err := strconv.Atoi(v); err == nil {
				o.MatchTimeout = time.Duration(ms) * time.Millisecond
			}
		}
	}
	return o
}

func (d Driver) Act(bot, game, actions string) error {
	// act:tictactoe:id (stream)
	if token, ok := d.token[fmt.Sprintf("%s:%s", bot, game)]; ok {
		err := d.pool.Do(radix.Cmd(nil, "XADD", fmt.Sprintf("act:tictactoe:%s", game), "*",
			"token", token,
			"actions", actions))
		if err != nil {
			return err
		}
	}
	return nil
}
