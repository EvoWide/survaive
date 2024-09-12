package server

import (
	"context"
	"encoding/json"
	"fmt"
	"survaive/bus"
	"sync"
	"time"
)

type Game struct {
	gameName   string
	channel    chan string
	eventQueue chan string
	State      *GameState
	mu         sync.Mutex
}

type GameState struct {
	Running bool
}

func NewGame(name string, dataChan chan string) *Game {
	gameCreated := &Game{
		gameName:   name,
		channel:    dataChan,
		eventQueue: make(chan string, 100), //100 events storaze size
		State:      &GameState{Running: false},
	}

	data, _ := json.Marshal(gameCreated.State)
	err := bus.GetRedisClient().HSet(context.Background(), name, data, 0).Err()

	if err != nil {
		panic(err)
	}

	return gameCreated
}

func (g *Game) enqueueEvents() {
	for event := range g.channel {
		g.mu.Lock()
		g.eventQueue <- event
		g.mu.Unlock()
	}
}

func (g *Game) Run() {
	ticker := time.NewTicker(200 * time.Millisecond) // 100 ms between each tick (10 TPS)
	defer ticker.Stop()

	// Start a goroutine to read from the channel and enqueue events
	go g.enqueueEvents()

	for g.State.Running {
		<-ticker.C
		g.gameTick()
	}
}

func (g *Game) Stop() {
	g.State.Running = false
	close(g.eventQueue)
}

func (g *Game) gameTick() {
	fmt.Printf("Tick - Game running : %v\n", g.State.Running)
}
