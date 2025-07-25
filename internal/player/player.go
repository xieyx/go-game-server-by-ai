package player

import (
	"sync"
	"time"
)

// Player represents a game player
type Player struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	mutex     sync.RWMutex
}

// NewPlayer creates a new player
func NewPlayer(id, name string) *Player {
	now := time.Now()
	return &Player{
		ID:        id,
		Name:      name,
		Level:     1,
		Score:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetID returns the player's ID
func (p *Player) GetID() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.ID
}

// GetName returns the player's name
func (p *Player) GetName() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Name
}

// SetName updates the player's name
func (p *Player) SetName(name string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Name = name
	p.UpdatedAt = time.Now()
}

// GetLevel returns the player's level
func (p *Player) GetLevel() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Level
}

// LevelUp increases the player's level
func (p *Player) LevelUp() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Level++
	p.UpdatedAt = time.Now()
}

// GetScore returns the player's score
func (p *Player) GetScore() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Score
}

// AddScore adds points to the player's score
func (p *Player) AddScore(points int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Score += points
	p.UpdatedAt = time.Now()
}

// GetCreatedAt returns when the player was created
func (p *Player) GetCreatedAt() time.Time {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.CreatedAt
}

// GetUpdatedAt returns when the player was last updated
func (p *Player) GetUpdatedAt() time.Time {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.UpdatedAt
}
