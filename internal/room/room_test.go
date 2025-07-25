package room

import (
	"testing"

	"github.com/xieyx/go-game-server-by-ai/internal/player"
)

func TestNewRoom(t *testing.T) {
	room := NewRoom("1", "Test Room", 4)

	if room.GetID() != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", room.GetID())
	}

	if room.GetName() != "Test Room" {
		t.Errorf("Expected Name to be 'Test Room', got '%s'", room.GetName())
	}

	if room.GetMaxPlayers() != 4 {
		t.Errorf("Expected MaxPlayers to be 4, got %d", room.GetMaxPlayers())
	}

	if room.GetStatus() != RoomStatusWaiting {
		t.Errorf("Expected Status to be RoomStatusWaiting, got '%s'", room.GetStatus())
	}

	if room.GetPlayerCount() != 0 {
		t.Errorf("Expected PlayerCount to be 0, got %d", room.GetPlayerCount())
	}
}

func TestRoomSetName(t *testing.T) {
	room := NewRoom("1", "Test Room", 4)
	room.SetName("New Room Name")

	if room.GetName() != "New Room Name" {
		t.Errorf("Expected Name to be 'New Room Name', got '%s'", room.GetName())
	}
}

func TestRoomSetStatus(t *testing.T) {
	room := NewRoom("1", "Test Room", 4)
	room.SetStatus(RoomStatusPlaying)

	if room.GetStatus() != RoomStatusPlaying {
		t.Errorf("Expected Status to be RoomStatusPlaying, got '%s'", room.GetStatus())
	}

	room.SetStatus(RoomStatusClosed)
	if room.GetStatus() != RoomStatusClosed {
		t.Errorf("Expected Status to be RoomStatusClosed, got '%s'", room.GetStatus())
	}
}

func TestRoomAddPlayer(t *testing.T) {
	room := NewRoom("1", "Test Room", 2)
	player1 := player.NewPlayer("1", "Alice")
	player2 := player.NewPlayer("2", "Bob")

	// Add first player
	err := room.AddPlayer(player1)
	if err != nil {
		t.Errorf("Expected no error when adding first player, got %v", err)
	}

	if room.GetPlayerCount() != 1 {
		t.Errorf("Expected PlayerCount to be 1, got %d", room.GetPlayerCount())
	}

	// Add second player
	err = room.AddPlayer(player2)
	if err != nil {
		t.Errorf("Expected no error when adding second player, got %v", err)
	}

	if room.GetPlayerCount() != 2 {
		t.Errorf("Expected PlayerCount to be 2, got %d", room.GetPlayerCount())
	}

	// Try to add third player (should fail because room is full)
	player3 := player.NewPlayer("3", "Charlie")
	err = room.AddPlayer(player3)
	if err == nil {
		t.Error("Expected error when adding player to full room, got nil")
	}
}

func TestRoomRemovePlayer(t *testing.T) {
	room := NewRoom("1", "Test Room", 2)
	player1 := player.NewPlayer("1", "Alice")
	player2 := player.NewPlayer("2", "Bob")

	// Add players
	room.AddPlayer(player1)
	room.AddPlayer(player2)

	// Remove first player
	err := room.RemovePlayer("1")
	if err != nil {
		t.Errorf("Expected no error when removing player, got %v", err)
	}

	if room.GetPlayerCount() != 1 {
		t.Errorf("Expected PlayerCount to be 1, got %d", room.GetPlayerCount())
	}

	// Try to remove non-existent player
	err = room.RemovePlayer("3")
	if err == nil {
		t.Error("Expected error when removing non-existent player, got nil")
	}
}

func TestRoomGetPlayer(t *testing.T) {
	room := NewRoom("1", "Test Room", 2)
	player1 := player.NewPlayer("1", "Alice")

	room.AddPlayer(player1)

	// Get existing player
	p, exists := room.GetPlayer("1")
	if !exists {
		t.Error("Expected player to exist")
	}

	if p.GetID() != "1" {
		t.Errorf("Expected player ID to be '1', got '%s'", p.GetID())
	}

	// Get non-existent player
	_, exists = room.GetPlayer("2")
	if exists {
		t.Error("Expected player to not exist")
	}
}

func TestRoomIsFullAndIsEmpty(t *testing.T) {
	room := NewRoom("1", "Test Room", 2)

	// Empty room
	if !room.IsEmpty() {
		t.Error("Expected room to be empty")
	}

	if room.IsFull() {
		t.Error("Expected room to not be full")
	}

	// Add one player
	player1 := player.NewPlayer("1", "Alice")
	room.AddPlayer(player1)

	if room.IsEmpty() {
		t.Error("Expected room to not be empty")
	}

	if room.IsFull() {
		t.Error("Expected room to not be full")
	}

	// Add second player (room should now be full)
	player2 := player.NewPlayer("2", "Bob")
	room.AddPlayer(player2)

	if room.IsEmpty() {
		t.Error("Expected room to not be empty")
	}

	if !room.IsFull() {
		t.Error("Expected room to be full")
	}
}
