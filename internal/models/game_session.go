package models

import "sync"

type GameSession struct {
	SessionId string
	Broadcast chan string
	Players   map[string]*Player
	Mu        sync.Mutex
}

type Player struct {
	Username string
	Id       string
	Events   chan string
}

func NewGameSession(sessionId string) *GameSession {
	return &GameSession{
		SessionId: sessionId,
		Broadcast: make(chan string),
		Players:   make(map[string]*Player),
	}
}

func (gs *GameSession) GetChannel() chan string {
	return gs.Broadcast
}

func (gs *GameSession) AddPlayerToSession(p *Player) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	gs.Players[p.Id] = p
}

func (gs *GameSession) RemovePlayerFromSession(playerId string) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	delete(gs.Players, playerId)
}
