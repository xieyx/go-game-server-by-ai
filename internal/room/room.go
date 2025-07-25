package room

import (
	"errors"
	"sync"
	"time"

	"github.com/xieyx/go-game-server-by-ai/internal/player"
)

// Room represents a game room
type Room struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Players   map[string]*player.Player `json:"players"`
	MaxPlayers int              `json:"max_players"`
	Status    RoomStatus        `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	mutex     sync.RWMutex
}

// RoomStatus represents the status of a room
type RoomStatus string

const (
	RoomStatusWaiting RoomStatus = "waiting"
	RoomStatusPlaying RoomStatus = "playing"
	RoomStatusClosed  RoomStatus = "closed"
)

// NewRoom creates a new room
func NewRoom(id, name string, maxPlayers int) *Room {
	now := time.Now()
	return &Room{
		ID:         id,
		Name:       name,
		Players:    make(map[string]*player.Player),
		MaxPlayers: maxPlayers,
		Status:     RoomStatusWaiting,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// GetID returns the room's ID
func (r *Room) GetID() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.ID
}

// GetName returns the room's name
func (r *Room) GetName() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Name
}

// SetName updates the room's name
func (r *Room) SetName(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Name = name
	r.UpdatedAt = time.Now()
}

// GetStatus returns the room's status
func (r *Room) GetStatus() RoomStatus {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Status
}

// SetStatus updates the room's status
func (r *Room) SetStatus(status RoomStatus) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Status = status
	r.UpdatedAt = time.Now()
}

// GetPlayerCount returns the number of players in the room
func (r *Room) GetPlayerCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.Players)
}

// GetMaxPlayers returns the maximum number of players allowed in the room
func (r *Room) GetMaxPlayers() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.MaxPlayers
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(p *player.Player) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Status == RoomStatusClosed {
		return errors.New("room is closed")
	}

	if len(r.Players) >= r.MaxPlayers {
		return errors.New("room is full")
	}

	if _, exists := r.Players[p.GetID()]; exists {
		return errors.New("player already in room")
	}

	r.Players[p.GetID()] = p
	r.UpdatedAt = time.Now()
	return nil
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(playerID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Players[playerID]; !exists {
		return errors.New("player not in room")
	}

	delete(r.Players, playerID)
	r.UpdatedAt = time.Now()
	return nil
}

// GetPlayer returns a player by ID
func (r *Room) GetPlayer(playerID string) (*player.Player, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	p, exists := r.Players[playerID]
	return p, exists
}

// GetAllPlayers returns all players in the room
func (r *Room) GetAllPlayers() map[string]*player.Player {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Create a copy to avoid concurrent access issues
	players := make(map[string]*player.Player)
	for id, player := range r.Players {
		players[id] = player
	}
	return players
}

// IsFull checks if the room is full
func (r *Room) IsFull() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.Players) >= r.MaxPlayers
}

// IsEmpty checks if the room is empty
func (r *Room) IsEmpty() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.Players) == 0
}

// GetCreatedAt returns when the room was created
func (r *Room) GetCreatedAt() time.Time {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.CreatedAt
}

// GetUpdatedAt returns when the room was last updated
func (r *Room) GetUpdatedAt() time.Time {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.UpdatedAt
}
