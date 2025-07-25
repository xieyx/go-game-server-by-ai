package character

import (
	"math/rand"
	"time"
)

// CharacterType represents the type of character
type CharacterType int

const (
	Warrior CharacterType = iota
	Mage
	Archer
)

// Character represents a game character
type Character struct {
	ID       string
	Name     string
	Type     CharacterType
	Level    int
	HP       int
	MaxHP    int
	MP       int
	MaxMP    int
	Attack   int
	Defense  int
	Speed    int
	Exp      int
	ExpToNext int
	Alive    bool
}

// NewCharacter creates a new character
func NewCharacter(id, name string, charType CharacterType) *Character {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	char := &Character{
		ID:      id,
		Name:    name,
		Type:    charType,
		Level:   1,
		Alive:   true,
	}

	// Set stats based on character type
	switch charType {
	case Warrior:
		char.MaxHP = 100 + rand.Intn(20)
		char.HP = char.MaxHP
		char.MaxMP = 30 + rand.Intn(10)
		char.MP = char.MaxMP
		char.Attack = 15 + rand.Intn(5)
		char.Defense = 10 + rand.Intn(5)
		char.Speed = 5 + rand.Intn(3)
	case Mage:
		char.MaxHP = 60 + rand.Intn(15)
		char.HP = char.MaxHP
		char.MaxMP = 100 + rand.Intn(20)
		char.MP = char.MaxMP
		char.Attack = 20 + rand.Intn(7)
		char.Defense = 5 + rand.Intn(3)
		char.Speed = 8 + rand.Intn(4)
	case Archer:
		char.MaxHP = 70 + rand.Intn(15)
		char.HP = char.MaxHP
		char.MaxMP = 50 + rand.Intn(15)
		char.MP = char.MaxMP
		char.Attack = 18 + rand.Intn(6)
		char.Defense = 7 + rand.Intn(4)
		char.Speed = 12 + rand.Intn(5)
	}

	char.ExpToNext = char.Level * 100

	return char
}

// IsAlive returns whether the character is alive
func (c *Character) IsAlive() bool {
	return c.Alive
}

// TakeDamage applies damage to the character
func (c *Character) TakeDamage(damage int) int {
	actualDamage := damage - c.Defense
	if actualDamage < 0 {
		actualDamage = 0
	}

	c.HP -= actualDamage
	if c.HP <= 0 {
		c.HP = 0
		c.Alive = false
	}

	return actualDamage
}

// Heal restores HP to the character
func (c *Character) Heal(amount int) int {
	c.HP += amount
	if c.HP > c.MaxHP {
		amount -= (c.HP - c.MaxHP)
		c.HP = c.MaxHP
	}

	return amount
}

// RestoreMP restores MP to the character
func (c *Character) RestoreMP(amount int) int {
	c.MP += amount
	if c.MP > c.MaxMP {
		amount -= (c.MP - c.MaxMP)
		c.MP = c.MaxMP
	}

	return amount
}

// GainExp adds experience points to the character
func (c *Character) GainExp(exp int) bool {
	c.Exp += exp

	// Check if character levels up
	if c.Exp >= c.ExpToNext {
		c.levelUp()
		return true
	}

	return false
}

// levelUp increases the character's level and stats
func (c *Character) levelUp() {
	c.Level++
	c.Exp -= c.ExpToNext
	c.ExpToNext = c.Level * 100

	// Increase stats based on character type
	switch c.Type {
	case Warrior:
		c.MaxHP += 10 + rand.Intn(5)
		c.MaxMP += 2 + rand.Intn(3)
		c.Attack += 3 + rand.Intn(2)
		c.Defense += 2 + rand.Intn(2)
		c.Speed += 1
	case Mage:
		c.MaxHP += 5 + rand.Intn(3)
		c.MaxMP += 8 + rand.Intn(4)
		c.Attack += 4 + rand.Intn(2)
		c.Defense += 1
		c.Speed += 2 + rand.Intn(1)
	case Archer:
		c.MaxHP += 7 + rand.Intn(4)
		c.MaxMP += 4 + rand.Intn(3)
		c.Attack += 3 + rand.Intn(2)
		c.Defense += 1 + rand.Intn(1)
		c.Speed += 3 + rand.Intn(2)
	}

	// Fully restore HP and MP on level up
	c.HP = c.MaxHP
	c.MP = c.MaxMP
}

// GetAttackDamage calculates the damage for a basic attack
func (c *Character) GetAttackDamage() int {
	baseDamage := c.Attack
	variance := rand.Intn(5) - 2 // -2 to +2 variance
	return baseDamage + variance
}

// GetNameWithType returns the character's name with type
func (c *Character) GetNameWithType() string {
	typeName := ""
	switch c.Type {
	case Warrior:
		typeName = "战士"
	case Mage:
		typeName = "法师"
	case Archer:
		typeName = "弓箭手"
	}

	return c.Name + "(" + typeName + ")"
}
