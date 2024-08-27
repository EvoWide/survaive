package managers

import (
	"log"
	"survaive/internal/models"
)

var gameSessions = make(map[string]*models.GameSession)

func RetrieveOrCreate(channelName string) *models.GameSession {
	if _, ok := gameSessions[channelName]; !ok {
		gameSessions[channelName] = models.NewGameSession(channelName)
	}
	return gameSessions[channelName]
}

func SendMessageToRoom(gs *models.GameSession, message string) {
	log.Printf("message sended in channel")

	gs.Mu.Lock()
	for _, client := range gs.Players {
		client.Events <- message // Send the update to each client
	}
	gs.Mu.Unlock()
}
