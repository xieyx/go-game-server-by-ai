package game

import (
	"errors"
	"sync"

	"github.com/xieyx/go-game-server-by-ai/internal/player"
	"github.com/xieyx/go-game-server-by-ai/internal/room"
)

// Game represents the game server
type Game struct {
	Rooms     map[string]*room.Room `json:"rooms"`
	Players   map[string]*player.Player `json:"players"`
	mutex     sync.RWMutex
}

// NewGame creates a new game server instance
func NewGame() *Game {
	return &Game{
		Rooms:   make(map[string]*room.Room),
		Players: make(map[string]*player.Player),
	}
}

// CreateRoom creates a new room
func (g *Game) CreateRoom(id, name string, maxPlayers int) *room.Room {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	r := room.NewRoom(id, name, maxPlayers)
	g.Rooms[id] = r
	return r
}

// GetRoom returns a room by ID
func (g *Game) GetRoom(roomID string) (*room.Room, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	r, exists := g.Rooms[roomID]
	return r, exists
}

// DeleteRoom deletes a room by ID
func (g *Game) DeleteRoom(roomID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	r, exists := g.Rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	// Remove all players from the room
	players := r.GetAllPlayers()
	for playerID := range players {
		r.RemovePlayer(playerID)
	}

	// Remove the room
	delete(g.Rooms, roomID)
	return nil
}

// GetAllRooms returns all rooms
func (g *Game) GetAllRooms() map[string]*room.Room {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	// Create a copy to avoid concurrent access issues
	rooms := make(map[string]*room.Room)
	for id, room := range g.Rooms {
		rooms[id] = room
	}
	return rooms
}

// RegisterPlayer registers a new player
func (g *Game) RegisterPlayer(id, name string) *player.Player {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	p := player.NewPlayer(id, name)
	g.Players[id] = p
	return p
}

// GetPlayer returns a player by ID
func (g *Game) GetPlayer(playerID string) (*player.Player, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	p, exists := g.Players[playerID]
	return p, exists
}

// UnregisterPlayer removes a player from the game
func (g *Game) UnregisterPlayer(playerID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, exists := g.Players[playerID]
	if !exists {
		return errors.New("player not found")
	}

	// Remove player from any rooms they're in
	for _, r := range g.Rooms {
		if _, inRoom := r.GetPlayer(playerID); inRoom {
			r.RemovePlayer(playerID)
		}
	}

	// Remove the player
	delete(g.Players, playerID)
	return nil
}

// GetAllPlayers returns all players
func (g *Game) GetAllPlayers() map[string]*player.Player {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	// Create a copy to avoid concurrent access issues
	players := make(map[string]*player.Player)
	for id, player := range g.Players {
		players[id] = player
	}
	return players
}

// JoinRoom adds a player to a room
func (g *Game) JoinRoom(playerID, roomID string) error {
	g.mutex.RLock()
	p, playerExists := g.Players[playerID]
	r, roomExists := g.Rooms[roomID]
	g.mutex.RUnlock()

	if !playerExists {
		return errors.New("player not found")
	}

	if !roomExists {
		return errors.New("room not found")
	}

	return r.AddPlayer(p)
}

// LeaveRoom removes a player from a room
func (g *Game) LeaveRoom(playerID, roomID string) error {
	g.mutex.RLock()
	r, roomExists := g.Rooms[roomID]
	g.mutex.RUnlock()

	if !roomExists {
		return errors.New("room not found")
	}

	return r.RemovePlayer(playerID)
}

// GetRoomCount returns the number of rooms
func (g *Game) GetRoomCount() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return len(g.Rooms)
}

// GetPlayerCount returns the number of players
func (g *Game) GetPlayerCount() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return len(g.Players)
}

// StartGame starts a game in a room
func (g *Game) StartGame(roomID string) error {
	g.mutex.RLock()
	r, exists := g.Rooms[roomID]
	g.mutex.RUnlock()

	if !exists {
		return errors.New("room not found")
	}

	if r.GetPlayerCount() < 2 {
		return errors.New("need at least 2 players to start game")
	}

	r.SetStatus(room.RoomStatusPlaying)
	return nil
}

// EndGame ends a game in a room
func (g *Game) EndGame(roomID string) error {
	g.mutex.RLock()
	r, exists := g.Rooms[roomID]
	g.mutex.RUnlock()

	if !exists {
		return errors.New("room not found")
	}

	r.SetStatus(room.RoomStatusWaiting)
	return nil
}
