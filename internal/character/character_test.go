package character

import (
	"testing"
)

func TestNewCharacter(t *testing.T) {
	char := NewCharacter("1", "战士一号", Warrior)

	if char.ID != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", char.ID)
	}

	if char.Name != "战士一号" {
		t.Errorf("Expected Name to be '战士一号', got '%s'", char.Name)
	}

	if char.Type != Warrior {
		t.Errorf("Expected Type to be Warrior, got %d", char.Type)
	}

	if char.Level != 1 {
		t.Errorf("Expected Level to be 1, got %d", char.Level)
	}

	if char.HP <= 0 {
		t.Error("Expected HP to be greater than 0")
	}

	if char.MaxHP <= 0 {
		t.Error("Expected MaxHP to be greater than 0")
	}

	if char.HP != char.MaxHP {
		t.Error("Expected HP to equal MaxHP at creation")
	}

	if char.Alive != true {
		t.Error("Expected Alive to be true at creation")
	}
}

func TestCharacterTakeDamage(t *testing.T) {
	char := NewCharacter("1", "法师一号", Mage)
	initialHP := char.HP
	damage := 20

	actualDamage := char.TakeDamage(damage)

	// Check that actual damage is calculated correctly (damage - defense)
	if actualDamage < 0 {
		t.Errorf("Expected actual damage to be non-negative, got %d", actualDamage)
	}

	// Check that HP was reduced
	if char.HP >= initialHP {
		t.Error("Expected HP to be reduced after taking damage")
	}

	// Test taking damage that would kill the character
	char.HP = 5
	char.TakeDamage(100)

	if char.Alive != false {
		t.Error("Expected character to be dead after taking lethal damage")
	}

	if char.HP != 0 {
		t.Error("Expected HP to be 0 when character is dead")
	}
}

func TestCharacterHeal(t *testing.T) {
	char := NewCharacter("1", "弓箭手一号", Archer)

	// Take some damage first
	char.TakeDamage(30)
	damageHP := char.HP

	// Heal some HP
	healAmount := 20
	actuallyHealed := char.Heal(healAmount)

	if actuallyHealed != healAmount {
		t.Errorf("Expected to heal %d HP, actually healed %d", healAmount, actuallyHealed)
	}

	if char.HP != damageHP+healAmount {
		t.Errorf("Expected HP to be %d, got %d", damageHP+healAmount, char.HP)
	}

	// Test overhealing
	char.HP = char.MaxHP - 10
	overheal := char.Heal(20)

	if overheal != 10 {
		t.Errorf("Expected to heal 10 HP from overheal, actually healed %d", overheal)
	}

	if char.HP != char.MaxHP {
		t.Error("Expected HP to be at max after overheal")
	}
}

func TestCharacterRestoreMP(t *testing.T) {
	char := NewCharacter("1", "法师一号", Mage)

	// Use some MP
	char.MP -= 30
	usedMP := char.MP

	// Restore some MP
	restoreAmount := 20
	actuallyRestored := char.RestoreMP(restoreAmount)

	if actuallyRestored != restoreAmount {
		t.Errorf("Expected to restore %d MP, actually restored %d", restoreAmount, actuallyRestored)
	}

	if char.MP != usedMP+restoreAmount {
		t.Errorf("Expected MP to be %d, got %d", usedMP+restoreAmount, char.MP)
	}

	// Test overrestoring
	char.MP = char.MaxMP - 10
	overrestore := char.RestoreMP(20)

	if overrestore != 10 {
		t.Errorf("Expected to restore 10 MP from overrestore, actually restored %d", overrestore)
	}

	if char.MP != char.MaxMP {
		t.Error("Expected MP to be at max after overrestore")
	}
}

func TestCharacterGainExp(t *testing.T) {
	char := NewCharacter("1", "战士一号", Warrior)
	initialLevel := char.Level
	initialExp := char.Exp

	// Gain some exp
	expGained := 50
	leveledUp := char.GainExp(expGained)

	if char.Exp != initialExp+expGained {
		t.Errorf("Expected Exp to be %d, got %d", initialExp+expGained, char.Exp)
	}

	if leveledUp != false {
		t.Error("Expected character not to level up with 50 exp")
	}

	// Gain enough exp to level up
	expToLevel := char.ExpToNext - char.Exp
	leveledUp = char.GainExp(expToLevel)

	if leveledUp != true {
		t.Error("Expected character to level up")
	}

	if char.Level != initialLevel+1 {
		t.Errorf("Expected Level to be %d, got %d", initialLevel+1, char.Level)
	}

	// Check that HP and MP are fully restored on level up
	if char.HP != char.MaxHP {
		t.Error("Expected HP to be fully restored on level up")
	}

	if char.MP != char.MaxMP {
		t.Error("Expected MP to be fully restored on level up")
	}
}

func TestCharacterGetAttackDamage(t *testing.T) {
	char := NewCharacter("1", "战士一号", Warrior)

	// Test that we get a reasonable damage value
	damage := char.GetAttackDamage()

	if damage <= 0 {
		t.Error("Expected attack damage to be greater than 0")
	}

	// Test that damage is based on character's attack stat
	if damage < char.Attack-2 || damage > char.Attack+2 {
		t.Error("Expected attack damage to be within variance of attack stat")
	}
}

func TestCharacterGetNameWithType(t *testing.T) {
	warrior := NewCharacter("1", "战士一号", Warrior)
	mage := NewCharacter("2", "法师一号", Mage)
	archer := NewCharacter("3", "弓箭手一号", Archer)

	if warrior.GetNameWithType() != "战士一号(战士)" {
		t.Errorf("Expected warrior name to be '战士一号(战士)', got '%s'", warrior.GetNameWithType())
	}

	if mage.GetNameWithType() != "法师一号(法师)" {
		t.Errorf("Expected mage name to be '法师一号(法师)', got '%s'", mage.GetNameWithType())
	}

	if archer.GetNameWithType() != "弓箭手一号(弓箭手)" {
		t.Errorf("Expected archer name to be '弓箭手一号(弓箭手)', got '%s'", archer.GetNameWithType())
	}
}
