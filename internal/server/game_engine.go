package server

import (
	"fmt"
	"survaive/sse"
	"sync"
)

type GameEngine struct {
	broker *sse.Broker
	rwm    *sync.RWMutex
	games  map[string]*Game
}

var (
	instance *GameEngine
	once     sync.Once
)

const (
	GAME_CHANNEL_NAME = "transport:game"
)

func InitGameEngine(broker *sse.Broker) *GameEngine {
	once.Do(func() {
		instance = &GameEngine{
			broker: broker,
			rwm:    &sync.RWMutex{},
			games:  make(map[string]*Game),
		}
	})
	return instance
}

func GetGameEngine() *GameEngine {
	return instance
}

func (ge *GameEngine) GetOrCreateGame(name string) *Game {
	game := ge.GetGame(name)
	if game != nil {
		return game
	}

	return ge.CreateGame(name)
}

func (ge *GameEngine) GetGame(name string) *Game {
	fmt.Println(ge.games)
	fmt.Println(name)
	ge.rwm.RLock()
	defer ge.rwm.RUnlock()

	return ge.games[name]
}

func (ge *GameEngine) CreateGame(name string) *Game {
	ge.rwm.Lock()
	defer ge.rwm.Unlock()

	// Subscribe to the room
	dataChan := ge.broker.AddSubscriber(name, "server")
	game := NewGame(name, dataChan)

	ge.games[name] = game
	return game
}
