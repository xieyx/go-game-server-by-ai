package combat

import (
	"github.com/xieyx/go-game-server-by-ai/internal/character"
)

// SkillType represents the type of skill
type SkillType int

const (
	// Basic attack skill
	BasicAttack SkillType = iota
	// Warrior skills
	PowerStrike
	Whirlwind
	ShieldBash
	// Mage skills
	Fireball
	Frostbolt
	Lightning
	Heal
	// Archer skills
	PiercingShot
	Multishot
	Trap
)

// Skill represents a combat skill
type Skill struct {
	ID          string
	Name        string
	Type        SkillType
	Description string
	MPCost      int
	Damage      int
	HealAmount  int
	TargetType  TargetType
}

// TargetType represents the type of target for a skill
type TargetType int

const (
	SingleTarget TargetType = iota
	AllEnemies
	AllAllies
	Self
)

// NewSkill creates a new skill
func NewSkill(id string, name string, skillType SkillType, description string, mpCost int, damage int, healAmount int, targetType TargetType) *Skill {
	return &Skill{
		ID:          id,
		Name:        name,
		Type:        skillType,
		Description: description,
		MPCost:      mpCost,
		Damage:      damage,
		HealAmount:  healAmount,
		TargetType:  targetType,
	}
}

// GetBasicAttack returns the basic attack skill
func GetBasicAttack() *Skill {
	return &Skill{
		ID:          "basic_attack",
		Name:        "普通攻击",
		Type:        BasicAttack,
		Description: "对单个敌人进行普通攻击",
		MPCost:      0,
		Damage:      10,
		HealAmount:  0,
		TargetType:  SingleTarget,
	}
}

// GetCharacterSkills returns the skills available to a character based on their type
func GetCharacterSkills(charType character.CharacterType) []*Skill {
	var skills []*Skill

	switch charType {
	case character.Warrior:
		skills = append(skills, GetBasicAttack())
		skills = append(skills, NewSkill("power_strike", "强力打击", PowerStrike, "对单个敌人造成大量伤害", 10, 25, 0, SingleTarget))
		skills = append(skills, NewSkill("whirlwind", "旋风斩", Whirlwind, "对所有敌人造成伤害", 20, 15, 0, AllEnemies))
		skills = append(skills, NewSkill("shield_bash", "盾击", ShieldBash, "对单个敌人造成伤害并使其眩晕1回合", 15, 20, 0, SingleTarget))
	case character.Mage:
		skills = append(skills, GetBasicAttack())
		skills = append(skills, NewSkill("fireball", "火球术", Fireball, "对单个敌人造成火焰伤害", 15, 30, 0, SingleTarget))
		skills = append(skills, NewSkill("frostbolt", "寒冰箭", Frostbolt, "对单个敌人造成冰霜伤害并减速", 12, 20, 0, SingleTarget))
		skills = append(skills, NewSkill("lightning", "闪电链", Lightning, "对所有敌人造成闪电伤害", 25, 18, 0, AllEnemies))
		skills = append(skills, NewSkill("heal", "治疗术", Heal, "恢复单个友方单位的生命值", 20, 0, 30, SingleTarget))
	case character.Archer:
		skills = append(skills, GetBasicAttack())
		skills = append(skills, NewSkill("piercing_shot", "穿刺射击", PiercingShot, "对单个敌人造成大量伤害，无视部分防御", 15, 28, 0, SingleTarget))
		skills = append(skills, NewSkill("multishot", "多重射击", Multishot, "对所有敌人造成伤害", 20, 12, 0, AllEnemies))
		skills = append(skills, NewSkill("trap", "陷阱", Trap, "在战场上设置陷阱，敌人触发时受到伤害", 10, 15, 0, SingleTarget))
	default:
		skills = append(skills, GetBasicAttack())
	}

	return skills
}

// CanUse checks if a character can use this skill
func (s *Skill) CanUse(c *character.Character) bool {
	return c.MP >= s.MPCost
}

// Use consumes the skill's MP cost
func (s *Skill) Use(c *character.Character) {
	if c.MP >= s.MPCost {
		c.MP -= s.MPCost
	}
}

// IsHealingSkill checks if the skill is a healing skill
func (s *Skill) IsHealingSkill() bool {
	return s.HealAmount > 0
}

// IsDamageSkill checks if the skill is a damage skill
func (s *Skill) IsDamageSkill() bool {
	return s.Damage > 0
}

// GetSkillName returns the name of the skill based on its type
func GetSkillName(skillType SkillType) string {
	switch skillType {
	case BasicAttack:
		return "普通攻击"
	case PowerStrike:
		return "强力打击"
	case Whirlwind:
		return "旋风斩"
	case ShieldBash:
		return "盾击"
	case Fireball:
		return "火球术"
	case Frostbolt:
		return "寒冰箭"
	case Lightning:
		return "闪电链"
	case Heal:
		return "治疗术"
	case PiercingShot:
		return "穿刺射击"
	case Multishot:
		return "多重射击"
	case Trap:
		return "陷阱"
	default:
		return "未知技能"
	}
}
