package player

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer("1", "Alice")

	if player.GetID() != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", player.GetID())
	}

	if player.GetName() != "Alice" {
		t.Errorf("Expected Name to be 'Alice', got '%s'", player.GetName())
	}

	if player.GetLevel() != 1 {
		t.Errorf("Expected Level to be 1, got %d", player.GetLevel())
	}

	if player.GetScore() != 0 {
		t.Errorf("Expected Score to be 0, got %d", player.GetScore())
	}

	if player.GetCreatedAt().IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if player.GetUpdatedAt().IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestPlayerSetName(t *testing.T) {
	player := NewPlayer("1", "Alice")
	player.SetName("Bob")

	if player.GetName() != "Bob" {
		t.Errorf("Expected Name to be 'Bob', got '%s'", player.GetName())
	}

	// Check that UpdatedAt was updated
	if player.GetUpdatedAt().Before(player.GetCreatedAt()) {
		t.Error("Expected UpdatedAt to be after CreatedAt")
	}
}

func TestPlayerLevelUp(t *testing.T) {
	player := NewPlayer("1", "Alice")
	initialLevel := player.GetLevel()
	player.LevelUp()

	if player.GetLevel() != initialLevel+1 {
		t.Errorf("Expected Level to be %d, got %d", initialLevel+1, player.GetLevel())
	}

	// Check that UpdatedAt was updated
	if player.GetUpdatedAt().Before(player.GetCreatedAt()) {
		t.Error("Expected UpdatedAt to be after CreatedAt")
	}
}

func TestPlayerAddScore(t *testing.T) {
	player := NewPlayer("1", "Alice")
	initialScore := player.GetScore()
	player.AddScore(100)

	if player.GetScore() != initialScore+100 {
		t.Errorf("Expected Score to be %d, got %d", initialScore+100, player.GetScore())
	}

	// Check that UpdatedAt was updated
	if player.GetUpdatedAt().Before(player.GetCreatedAt()) {
		t.Error("Expected UpdatedAt to be after CreatedAt")
	}

	// Test adding negative score
	player.AddScore(-50)
	if player.GetScore() != initialScore+50 {
		t.Errorf("Expected Score to be %d, got %d", initialScore+50, player.GetScore())
	}
}
