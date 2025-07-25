package battle

import (
	"testing"

	"github.com/xieyx/go-game-server-by-ai/internal/character"
	"github.com/xieyx/go-game-server-by-ai/internal/combat"
)

func TestNewBattle(t *testing.T) {
	battle := NewBattle("test_battle")

	if battle.ID != "test_battle" {
		t.Errorf("Expected ID to be 'test_battle', got '%s'", battle.ID)
	}

	if battle.State != BattleStateNotStarted {
		t.Errorf("Expected State to be BattleStateNotStarted, got %d", battle.State)
	}

	if len(battle.Participants) != 0 {
		t.Errorf("Expected Participants to be empty, got %d", len(battle.Participants))
	}

	if len(battle.Log) != 0 {
		t.Errorf("Expected Log to be empty, got %d", len(battle.Log))
	}
}

func TestAddParticipant(t *testing.T) {
	battle := NewBattle("test_battle")
	char := character.NewCharacter("1", "战士一号", character.Warrior)

	battle.AddParticipant(char, true)

	if len(battle.Participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(battle.Participants))
	}

	participant := battle.Participants[0]
	if participant.Character != char {
		t.Error("Expected participant character to match")
	}

	if participant.IsPlayer != true {
		t.Error("Expected participant to be player")
	}

	if participant.IsAlive != true {
		t.Error("Expected participant to be alive")
	}
}

func TestBattleStart(t *testing.T) {
	battle := NewBattle("test_battle")

	// Try to start battle with no participants
	battle.Start()

	if len(battle.Log) == 0 {
		t.Error("Expected log entry for battle with no participants")
	}

	// Add participants and start battle
	player := character.NewCharacter("1", "战士一号", character.Warrior)
	enemy := character.NewCharacter("2", "哥布林", character.Warrior)

	battle.AddParticipant(player, true)
	battle.AddParticipant(enemy, false)

	battle.Start()

	if battle.State != BattleStateInProgress {
		t.Errorf("Expected State to be BattleStateInProgress, got %d", battle.State)
	}

	if battle.CurrentRound != 1 {
		t.Errorf("Expected CurrentRound to be 1, got %d", battle.CurrentRound)
	}

	if len(battle.Log) < 3 {
		t.Errorf("Expected at least 3 log entries, got %d", len(battle.Log))
	}
}

func TestGetCurrentParticipant(t *testing.T) {
	battle := NewBattle("test_battle")
	player := character.NewCharacter("1", "战士一号", character.Warrior)
	enemy := character.NewCharacter("2", "哥布林", character.Warrior)

	battle.AddParticipant(player, true)
	battle.AddParticipant(enemy, false)
	battle.Start()

	current := battle.GetCurrentParticipant()

	if current == nil {
		t.Error("Expected current participant to not be nil")
	}

	// The participant with higher speed should go first
	// Since we're using random values, we can't guarantee which will be faster
	// But we can check that it's one of our participants
	if current != battle.Participants[0] && current != battle.Participants[1] {
		t.Error("Expected current participant to be one of the battle participants")
	}
}

func TestSelectSkill(t *testing.T) {
	battle := NewBattle("test_battle")
	player := character.NewCharacter("1", "法师一号", character.Mage)
	player.MP = 30 // Ensure enough MP

	battle.AddParticipant(player, true)
	battle.AddParticipant(character.NewCharacter("2", "哥布林", character.Warrior), false)
	battle.Start()

	participant := battle.Participants[0]
	skill := combat.GetBasicAttack()

	err := battle.SelectSkill(participant, skill)

	if err != nil {
		t.Errorf("Expected no error when selecting skill, got %v", err)
	}

	if participant.SelectedSkill != skill {
		t.Error("Expected participant's selected skill to match")
	}

	// Test selecting skill when not enough MP
	player.MP = 0
	err = battle.SelectSkill(participant, skill)

	// Basic attack costs 0 MP, so it should still work
	// Let's test with a skill that costs MP
	fireball := &combat.Skill{
		ID:          "fireball",
		Name:        "火球术",
		Type:        combat.Fireball,
		Description: "对单个敌人造成火焰伤害",
		MPCost:      15,
		Damage:      30,
		HealAmount:  0,
		TargetType:  combat.SingleTarget,
	}

	err = battle.SelectSkill(participant, fireball)

	if err == nil {
		t.Error("Expected error when not enough MP")
	}
}

func TestSelectTarget(t *testing.T) {
	battle := NewBattle("test_battle")
	player := character.NewCharacter("1", "战士一号", character.Warrior)
	enemy := character.NewCharacter("2", "哥布林", character.Warrior)

	battle.AddParticipant(player, true)
	battle.AddParticipant(enemy, false)
	battle.Start()

	participant := battle.Participants[0]
	target := battle.Participants[1]

	err := battle.SelectTarget(participant, target)

	if err != nil {
		t.Errorf("Expected no error when selecting target, got %v", err)
	}

	if participant.SelectedTarget != target {
		t.Error("Expected participant's selected target to match")
	}

	// Test selecting dead target
	target.IsAlive = false
	target.Character.HP = 0

	err = battle.SelectTarget(participant, target)

	if err == nil {
		t.Error("Expected error when selecting dead target")
	}
}

func TestExecuteTurn(t *testing.T) {
	battle := NewBattle("test_battle")
	player := character.NewCharacter("1", "战士一号", character.Warrior)
	enemy := character.NewCharacter("2", "哥布林", character.Warrior)

	battle.AddParticipant(player, true)
	battle.AddParticipant(enemy, false)
	battle.Start()

	// Execute a few turns
	for i := 0; i < 3; i++ {
		battle.ExecuteTurn()
	}

	// Check that battle is still in progress or has ended
	if battle.State != BattleStateInProgress && battle.State != BattleStatePlayerWon && battle.State != BattleStateEnemiesWon {
		t.Errorf("Expected battle state to be in progress or ended, got %d", battle.State)
	}
}

func TestBattleEnd(t *testing.T) {
	battle := NewBattle("test_battle")
	player := character.NewCharacter("1", "战士一号", character.Warrior)
	enemy := character.NewCharacter("2", "哥布林", character.Warrior)

	battle.AddParticipant(player, true)
	battle.AddParticipant(enemy, false)
	battle.Start()

	// Kill the enemy
	enemy.HP = 0
	enemy.Alive = false

	// Execute turn to check battle end
	battle.ExecuteTurn()

	// Battle should be over with player winning
	if battle.State != BattleStatePlayerWon {
		t.Errorf("Expected battle state to be BattleStatePlayerWon, got %d", battle.State)
	}

	// Check that rewards were calculated
	if battle.Reward.Exp == 0 {
		t.Error("Expected experience reward")
	}
}

func TestApplyEffect(t *testing.T) {
	battle := NewBattle("test_battle")
	participant := &BattleParticipant{
		Character: character.NewCharacter("1", "战士一号", character.Warrior),
		IsPlayer:  true,
		IsAlive:   true,
		Effects:   make([]BattleEffect, 0),
	}

	battle.applyEffect(participant, EffectStun, 2, 0)

	if len(participant.Effects) != 1 {
		t.Errorf("Expected 1 effect, got %d", len(participant.Effects))
	}

	effect := participant.Effects[0]
	if effect.Type != EffectStun {
		t.Errorf("Expected effect type to be EffectStun, got %d", effect.Type)
	}

	if effect.Duration != 2 {
		t.Errorf("Expected effect duration to be 2, got %d", effect.Duration)
	}
}

func TestIsStunned(t *testing.T) {
	battle := NewBattle("test_battle")
	participant := &BattleParticipant{
		Character: character.NewCharacter("1", "战士一号", character.Warrior),
		IsPlayer:  true,
		IsAlive:   true,
		Effects:   make([]BattleEffect, 0),
	}

	// Test with no stun effect
	if battle.isStunned(participant) {
		t.Error("Expected participant to not be stunned")
	}

	// Add stun effect
	battle.applyEffect(participant, EffectStun, 1, 0)

	// Test with stun effect
	if !battle.isStunned(participant) {
		t.Error("Expected participant to be stunned")
	}
}

func TestGetLog(t *testing.T) {
	battle := NewBattle("test_battle")

	// Add some log entries
	battle.Log = append(battle.Log, BattleLogEntry{Round: 1, Turn: 1, Text: "Test entry"})
	battle.Log = append(battle.Log, BattleLogEntry{Round: 1, Turn: 2, Text: "Another entry"})

	log := battle.GetLog()

	if len(log) != 2 {
		t.Errorf("Expected log length to be 2, got %d", len(log))
	}

	if log[0].Text != "Test entry" {
		t.Errorf("Expected first log entry to be 'Test entry', got '%s'", log[0].Text)
	}
}

func TestGetState(t *testing.T) {
	battle := NewBattle("test_battle")

	state := battle.GetState()

	if state != BattleStateNotStarted {
		t.Errorf("Expected state to be BattleStateNotStarted, got %d", state)
	}
}

func TestGetReward(t *testing.T) {
	battle := NewBattle("test_battle")
	battle.Reward = BattleReward{Exp: 100, Gold: 50, Items: []string{"sword"}}

	reward := battle.GetReward()

	if reward.Exp != 100 {
		t.Errorf("Expected reward exp to be 100, got %d", reward.Exp)
	}

	if reward.Gold != 50 {
		t.Errorf("Expected reward gold to be 50, got %d", reward.Gold)
	}

	if len(reward.Items) != 1 {
		t.Errorf("Expected 1 reward item, got %d", len(reward.Items))
	}
}

func TestGetAliveParticipants(t *testing.T) {
	battle := NewBattle("test_battle")
	player := character.NewCharacter("1", "战士一号", character.Warrior)
	enemy := character.NewCharacter("2", "哥布林", character.Warrior)

	battle.AddParticipant(player, true)
	battle.AddParticipant(enemy, false)

	// Kill the enemy
	enemy.HP = 0
	enemy.Alive = false
	battle.Participants[1].IsAlive = false

	alive := battle.GetAliveParticipants()

	if len(alive) != 1 {
		t.Errorf("Expected 1 alive participant, got %d", len(alive))
	}

	if alive[0].IsPlayer != true {
		t.Error("Expected alive participant to be player")
	}
}
