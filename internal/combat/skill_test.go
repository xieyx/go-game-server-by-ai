package combat

import (
	"testing"

	"github.com/xieyx/go-game-server-by-ai/internal/character"
)

func TestNewSkill(t *testing.T) {
	skill := NewSkill("test_skill", "测试技能", Fireball, "这是一个测试技能", 10, 20, 0, SingleTarget)

	if skill.ID != "test_skill" {
		t.Errorf("Expected ID to be 'test_skill', got '%s'", skill.ID)
	}

	if skill.Name != "测试技能" {
		t.Errorf("Expected Name to be '测试技能', got '%s'", skill.Name)
	}

	if skill.Type != Fireball {
		t.Errorf("Expected Type to be Fireball, got %d", skill.Type)
	}

	if skill.Description != "这是一个测试技能" {
		t.Errorf("Expected Description to be '这是一个测试技能', got '%s'", skill.Description)
	}

	if skill.MPCost != 10 {
		t.Errorf("Expected MPCost to be 10, got %d", skill.MPCost)
	}

	if skill.Damage != 20 {
		t.Errorf("Expected Damage to be 20, got %d", skill.Damage)
	}

	if skill.HealAmount != 0 {
		t.Errorf("Expected HealAmount to be 0, got %d", skill.HealAmount)
	}

	if skill.TargetType != SingleTarget {
		t.Errorf("Expected TargetType to be SingleTarget, got %d", skill.TargetType)
	}
}

func TestGetBasicAttack(t *testing.T) {
	skill := GetBasicAttack()

	if skill.ID != "basic_attack" {
		t.Errorf("Expected ID to be 'basic_attack', got '%s'", skill.ID)
	}

	if skill.Name != "普通攻击" {
		t.Errorf("Expected Name to be '普通攻击', got '%s'", skill.Name)
	}

	if skill.Type != BasicAttack {
		t.Errorf("Expected Type to be BasicAttack, got %d", skill.Type)
	}

	if skill.MPCost != 0 {
		t.Errorf("Expected MPCost to be 0, got %d", skill.MPCost)
	}

	if skill.Damage != 10 {
		t.Errorf("Expected Damage to be 10, got %d", skill.Damage)
	}

	if skill.HealAmount != 0 {
		t.Errorf("Expected HealAmount to be 0, got %d", skill.HealAmount)
	}

	if skill.TargetType != SingleTarget {
		t.Errorf("Expected TargetType to be SingleTarget, got %d", skill.TargetType)
	}
}

func TestGetCharacterSkills(t *testing.T) {
	// Test warrior skills
	warriorSkills := GetCharacterSkills(character.Warrior)
	if len(warriorSkills) == 0 {
		t.Error("Expected warrior to have skills")
	}

	// Check that warrior has basic attack
	hasBasicAttack := false
	for _, skill := range warriorSkills {
		if skill.Type == BasicAttack {
			hasBasicAttack = true
			break
		}
	}

	if !hasBasicAttack {
		t.Error("Expected warrior to have basic attack skill")
	}

	// Test mage skills
	mageSkills := GetCharacterSkills(character.Mage)
	if len(mageSkills) == 0 {
		t.Error("Expected mage to have skills")
	}

	// Check that mage has heal skill
	hasHeal := false
	for _, skill := range mageSkills {
		if skill.Type == Heal {
			hasHeal = true
			break
		}
	}

	if !hasHeal {
		t.Error("Expected mage to have heal skill")
	}

	// Test archer skills
	archerSkills := GetCharacterSkills(character.Archer)
	if len(archerSkills) == 0 {
		t.Error("Expected archer to have skills")
	}
}

func TestSkillCanUse(t *testing.T) {
	skill := NewSkill("test_skill", "测试技能", Fireball, "这是一个测试技能", 10, 20, 0, SingleTarget)

	// Create a character with enough MP
	char := character.NewCharacter("1", "法师一号", character.Mage)
	char.MP = 20

	if !skill.CanUse(char) {
		t.Error("Expected character with enough MP to be able to use skill")
	}

	// Create a character with not enough MP
	char.MP = 5

	if skill.CanUse(char) {
		t.Error("Expected character with not enough MP to not be able to use skill")
	}
}

func TestSkillUse(t *testing.T) {
	skill := NewSkill("test_skill", "测试技能", Fireball, "这是一个测试技能", 10, 20, 0, SingleTarget)

	// Create a character with enough MP
	char := character.NewCharacter("1", "法师一号", character.Mage)
	char.MP = 20
	initialMP := char.MP

	skill.Use(char)

	if char.MP != initialMP-skill.MPCost {
		t.Errorf("Expected MP to be %d, got %d", initialMP-skill.MPCost, char.MP)
	}

	// Test using skill with not enough MP
	char.MP = 5
	skill.Use(char)

	// MP should not go below 0
	if char.MP < 0 {
		t.Error("Expected MP to not go below 0")
	}
}

func TestSkillIsHealingSkill(t *testing.T) {
	// Test healing skill
	healSkill := NewSkill("heal", "治疗术", Heal, "恢复生命值", 20, 0, 30, SingleTarget)

	if !healSkill.IsHealingSkill() {
		t.Error("Expected heal skill to be a healing skill")
	}

	// Test damage skill
	damageSkill := NewSkill("fireball", "火球术", Fireball, "造成伤害", 15, 30, 0, SingleTarget)

	if damageSkill.IsHealingSkill() {
		t.Error("Expected damage skill to not be a healing skill")
	}
}

func TestSkillIsDamageSkill(t *testing.T) {
	// Test damage skill
	damageSkill := NewSkill("fireball", "火球术", Fireball, "造成伤害", 15, 30, 0, SingleTarget)

	if !damageSkill.IsDamageSkill() {
		t.Error("Expected damage skill to be a damage skill")
	}

	// Test healing skill
	healSkill := NewSkill("heal", "治疗术", Heal, "恢复生命值", 20, 0, 30, SingleTarget)

	if healSkill.IsDamageSkill() {
		t.Error("Expected heal skill to not be a damage skill")
	}
}

func TestGetSkillName(t *testing.T) {
	if GetSkillName(BasicAttack) != "普通攻击" {
		t.Errorf("Expected BasicAttack name to be '普通攻击', got '%s'", GetSkillName(BasicAttack))
	}

	if GetSkillName(Fireball) != "火球术" {
		t.Errorf("Expected Fireball name to be '火球术', got '%s'", GetSkillName(Fireball))
	}

	if GetSkillName(Heal) != "治疗术" {
		t.Errorf("Expected Heal name to be '治疗术', got '%s'", GetSkillName(Heal))
	}

	// Test unknown skill
	if GetSkillName(999) != "未知技能" {
		t.Errorf("Expected unknown skill name to be '未知技能', got '%s'", GetSkillName(999))
	}
}
