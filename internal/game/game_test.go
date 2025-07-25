package game

import (
	"testing"

	"github.com/xieyx/go-game-server-by-ai/internal/room"
)

func TestNewGame(t *testing.T) {
	game := NewGame()

	if game.GetRoomCount() != 0 {
		t.Errorf("Expected RoomCount to be 0, got %d", game.GetRoomCount())
	}

	if game.GetPlayerCount() != 0 {
		t.Errorf("Expected PlayerCount to be 0, got %d", game.GetPlayerCount())
	}
}

func TestGameCreateRoom(t *testing.T) {
	game := NewGame()
	room := game.CreateRoom("1", "Test Room", 4)

	if room.GetID() != "1" {
		t.Errorf("Expected room ID to be '1', got '%s'", room.GetID())
	}

	if game.GetRoomCount() != 1 {
		t.Errorf("Expected RoomCount to be 1, got %d", game.GetRoomCount())
	}

	// Check that room exists in game
	r, exists := game.GetRoom("1")
	if !exists {
		t.Error("Expected room to exist in game")
	}

	if r.GetID() != "1" {
		t.Errorf("Expected room ID to be '1', got '%s'", r.GetID())
	}
}

func TestGameRegisterPlayer(t *testing.T) {
	game := NewGame()
	player := game.RegisterPlayer("1", "Alice")

	if player.GetID() != "1" {
		t.Errorf("Expected player ID to be '1', got '%s'", player.GetID())
	}

	if game.GetPlayerCount() != 1 {
		t.Errorf("Expected PlayerCount to be 1, got %d", game.GetPlayerCount())
	}

	// Check that player exists in game
	p, exists := game.GetPlayer("1")
	if !exists {
		t.Error("Expected player to exist in game")
	}

	if p.GetID() != "1" {
		t.Errorf("Expected player ID to be '1', got '%s'", p.GetID())
	}
}

func TestGameJoinAndLeaveRoom(t *testing.T) {
	game := NewGame()

	// Create room and player
	game.CreateRoom("1", "Test Room", 4)
	game.RegisterPlayer("1", "Alice")

	// Join room
	err := game.JoinRoom("1", "1")
	if err != nil {
		t.Errorf("Expected no error when joining room, got %v", err)
	}

	// Check that player is in room
	r, _ := game.GetRoom("1")
	p, exists := r.GetPlayer("1")
	if !exists {
		t.Error("Expected player to be in room")
	}

	if p.GetID() != "1" {
		t.Errorf("Expected player ID to be '1', got '%s'", p.GetID())
	}

	// Leave room
	err = game.LeaveRoom("1", "1")
	if err != nil {
		t.Errorf("Expected no error when leaving room, got %v", err)
	}

	// Check that player is no longer in room
	_, exists = r.GetPlayer("1")
	if exists {
		t.Error("Expected player to not be in room")
	}
}

func TestGameStartAndEndGame(t *testing.T) {
	game := NewGame()

	// Create room and players
	game.CreateRoom("1", "Test Room", 4)
	game.RegisterPlayer("1", "Alice")
	game.RegisterPlayer("2", "Bob")

	// Add players to room
	game.JoinRoom("1", "1")
	game.JoinRoom("2", "1")

	// Start game
	err := game.StartGame("1")
	if err != nil {
		t.Errorf("Expected no error when starting game, got %v", err)
	}

	// Check that room status is playing
	r, _ := game.GetRoom("1")
	if r.GetStatus() != room.RoomStatusPlaying {
		t.Errorf("Expected room status to be RoomStatusPlaying, got '%s'", r.GetStatus())
	}

	// End game
	err = game.EndGame("1")
	if err != nil {
		t.Errorf("Expected no error when ending game, got %v", err)
	}

	// Check that room status is waiting
	if r.GetStatus() != room.RoomStatusWaiting {
		t.Errorf("Expected room status to be RoomStatusWaiting, got '%s'", r.GetStatus())
	}
}

func TestGameDeleteRoom(t *testing.T) {
	game := NewGame()

	// Create room and player
	game.CreateRoom("1", "Test Room", 4)
	game.RegisterPlayer("1", "Alice")

	// Add player to room
	game.JoinRoom("1", "1")

	// Delete room
	err := game.DeleteRoom("1")
	if err != nil {
		t.Errorf("Expected no error when deleting room, got %v", err)
	}

	// Check that room no longer exists
	_, exists := game.GetRoom("1")
	if exists {
		t.Error("Expected room to not exist")
	}
}

func TestGameUnregisterPlayer(t *testing.T) {
	game := NewGame()

	// Create room and player
	game.CreateRoom("1", "Test Room", 4)
	game.RegisterPlayer("1", "Alice")

	// Add player to room
	game.JoinRoom("1", "1")

	// Unregister player
	err := game.UnregisterPlayer("1")
	if err != nil {
		t.Errorf("Expected no error when unregistering player, got %v", err)
	}

	// Check that player no longer exists in game
	_, exists := game.GetPlayer("1")
	if exists {
		t.Error("Expected player to not exist in game")
	}

	// Check that player is no longer in room
	r, _ := game.GetRoom("1")
	_, exists = r.GetPlayer("1")
	if exists {
		t.Error("Expected player to not be in room")
	}
}
