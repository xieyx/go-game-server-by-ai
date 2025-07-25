package battle

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/xieyx/go-game-server-by-ai/internal/character"
	"github.com/xieyx/go-game-server-by-ai/internal/combat"
)

// Battle represents a combat battle
type Battle struct {
	ID            string
	Participants  []*BattleParticipant
	TurnOrder     []*BattleParticipant
	CurrentTurn   int
	CurrentRound  int
	State         BattleState
	Log           []BattleLogEntry
	Reward        BattleReward
}

// BattleParticipant represents a participant in a battle
type BattleParticipant struct {
	Character    *character.Character
	IsPlayer     bool
	IsAlive      bool
	Speed        int
	Effects      []BattleEffect
	SelectedSkill *combat.Skill
	SelectedTarget *BattleParticipant
}

// BattleState represents the state of a battle
type BattleState int

const (
	BattleStateNotStarted BattleState = iota
	BattleStateInProgress
	BattleStatePlayerWon
	BattleStateEnemiesWon
	BattleStateDraw
)

// BattleEffect represents a status effect in battle
type BattleEffect struct {
	Type      EffectType
	Duration  int
	Remaining int
	Value     int
}

// EffectType represents the type of effect
type EffectType int

const (
	EffectNone EffectType = iota
	EffectStun
	EffectPoison
	EffectSlow
	EffectBuff
)

// BattleLogEntry represents an entry in the battle log
type BattleLogEntry struct {
	Round int
	Turn  int
	Text  string
}

// BattleReward represents the rewards for winning a battle
type BattleReward struct {
	Exp  int
	Gold int
	Items []string
}

// NewBattle creates a new battle
func NewBattle(id string) *Battle {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	return &Battle{
		ID:          id,
		Participants: make([]*BattleParticipant, 0),
		TurnOrder:   make([]*BattleParticipant, 0),
		CurrentTurn: 0,
		CurrentRound: 0,
		State:       BattleStateNotStarted,
		Log:         make([]BattleLogEntry, 0),
		Reward:      BattleReward{Exp: 0, Gold: 0, Items: make([]string, 0)},
	}
}

// AddParticipant adds a participant to the battle
func (b *Battle) AddParticipant(char *character.Character, isPlayer bool) {
	participant := &BattleParticipant{
		Character: char,
		IsPlayer:  isPlayer,
		IsAlive:   char.IsAlive(),
		Speed:     char.Speed,
		Effects:   make([]BattleEffect, 0),
	}

	b.Participants = append(b.Participants, participant)
}

// Start begins the battle
func (b *Battle) Start() {
	if len(b.Participants) < 2 {
		b.Log = append(b.Log, BattleLogEntry{
			Round: 0,
			Turn:  0,
			Text:  "战斗至少需要两个参与者",
		})
		return
	}

	// Initialize turn order based on speed
	b.TurnOrder = make([]*BattleParticipant, len(b.Participants))
	copy(b.TurnOrder, b.Participants)

	// Sort by speed (higher speed goes first)
	sort.Slice(b.TurnOrder, func(i, j int) bool {
		return b.TurnOrder[i].Speed > b.TurnOrder[j].Speed
	})

	b.State = BattleStateInProgress
	b.CurrentRound = 1
	b.CurrentTurn = 0

	b.Log = append(b.Log, BattleLogEntry{
		Round: 0,
		Turn:  0,
		Text:  fmt.Sprintf("战斗开始！第%d回合", b.CurrentRound),
	})

	// Log participants
	for _, p := range b.Participants {
		b.Log = append(b.Log, BattleLogEntry{
			Round: 0,
			Turn:  0,
			Text:  fmt.Sprintf("%s 加入了战斗", p.Character.GetNameWithType()),
		})
	}
}

// GetCurrentParticipant returns the participant whose turn it is
func (b *Battle) GetCurrentParticipant() *BattleParticipant {
	if b.State != BattleStateInProgress || len(b.TurnOrder) == 0 {
		return nil
	}

	return b.TurnOrder[b.CurrentTurn]
}

// SelectSkill allows a participant to select a skill for their turn
func (b *Battle) SelectSkill(participant *BattleParticipant, skill *combat.Skill) error {
	if b.State != BattleStateInProgress {
		return fmt.Errorf("战斗未进行中")
	}

	if !participant.IsAlive {
		return fmt.Errorf("角色已死亡")
	}

	if !skill.CanUse(participant.Character) {
		return fmt.Errorf("角色没有足够的MP使用此技能")
	}

	participant.SelectedSkill = skill
	return nil
}

// SelectTarget allows a participant to select a target for their turn
func (b *Battle) SelectTarget(participant *BattleParticipant, target *BattleParticipant) error {
	if b.State != BattleStateInProgress {
		return fmt.Errorf("战斗未进行中")
	}

	if !participant.IsAlive {
		return fmt.Errorf("角色已死亡")
	}

	if !target.IsAlive {
		return fmt.Errorf("目标已死亡")
	}

	// Check if target is valid based on skill type
	if participant.SelectedSkill != nil {
		switch participant.SelectedSkill.TargetType {
		case combat.SingleTarget:
			// Any alive participant is a valid target
		case combat.AllEnemies:
			// Any enemy is a valid target (we'll apply to all enemies later)
		case combat.AllAllies:
			// Any ally is a valid target (we'll apply to all allies later)
		case combat.Self:
			if target != participant {
				return fmt.Errorf("该技能只能对自己使用")
			}
		}
	}

	participant.SelectedTarget = target
	return nil
}

// ExecuteTurn executes the current participant's turn
func (b *Battle) ExecuteTurn() {
	if b.State != BattleStateInProgress {
		return
	}

	currentParticipant := b.GetCurrentParticipant()
	if currentParticipant == nil {
		return
	}

	// Check if participant is stunned
	if b.isStunned(currentParticipant) {
		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  fmt.Sprintf("%s 被眩晕，跳过回合", currentParticipant.Character.GetNameWithType()),
		})
		b.nextTurn()
		return
	}

	// If no skill selected, use basic attack
	if currentParticipant.SelectedSkill == nil {
		currentParticipant.SelectedSkill = combat.GetBasicAttack()
	}

	// If no target selected, select a random valid target
	if currentParticipant.SelectedTarget == nil {
		target, err := b.getRandomValidTarget(currentParticipant)
		if err != nil {
			b.Log = append(b.Log, BattleLogEntry{
				Round: b.CurrentRound,
				Turn:  b.CurrentTurn,
				Text:  fmt.Sprintf("%s 无法选择目标: %s", currentParticipant.Character.GetNameWithType(), err.Error()),
			})
			b.nextTurn()
			return
		}
		currentParticipant.SelectedTarget = target
	}

	// Execute the skill
	b.executeSkill(currentParticipant, currentParticipant.SelectedSkill, currentParticipant.SelectedTarget)

	// Reset selections for next turn
	currentParticipant.SelectedSkill = nil
	currentParticipant.SelectedTarget = nil

	// Check if battle is over
	if b.checkBattleEnd() {
		return
	}

	// Move to next turn
	b.nextTurn()
}

// executeSkill executes a skill
func (b *Battle) executeSkill(user *BattleParticipant, skill *combat.Skill, target *BattleParticipant) {
	// Use the skill (consume MP)
	skill.Use(user.Character)

	// Log skill usage
	b.Log = append(b.Log, BattleLogEntry{
		Round: b.CurrentRound,
		Turn:  b.CurrentTurn,
		Text:  fmt.Sprintf("%s 使用了 %s", user.Character.GetNameWithType(), skill.Name),
	})

	// Apply skill effects based on target type
	switch skill.TargetType {
	case combat.SingleTarget:
		b.applySkillToTarget(user, skill, target)
	case combat.AllEnemies:
		// Apply to all enemies
		for _, p := range b.Participants {
			if p.IsAlive && p.IsPlayer != user.IsPlayer {
				b.applySkillToTarget(user, skill, p)
			}
		}
	case combat.AllAllies:
		// Apply to all allies
		for _, p := range b.Participants {
			if p.IsAlive && p.IsPlayer == user.IsPlayer {
				b.applySkillToTarget(user, skill, p)
			}
		}
	case combat.Self:
		b.applySkillToTarget(user, skill, user)
	}
}

// applySkillToTarget applies a skill to a target
func (b *Battle) applySkillToTarget(user *BattleParticipant, skill *combat.Skill, target *BattleParticipant) {
	if skill.IsDamageSkill() {
		// Calculate damage
		damage := skill.Damage
		if skill.Type == combat.BasicAttack {
			damage = user.Character.GetAttackDamage()
		}

		// Apply damage
		actualDamage := target.Character.TakeDamage(damage)

		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  fmt.Sprintf("%s 对 %s 造成了 %d 点伤害", user.Character.GetNameWithType(), target.Character.GetNameWithType(), actualDamage),
		})

		// Check if target is defeated
		if !target.Character.IsAlive() {
			target.IsAlive = false
			b.Log = append(b.Log, BattleLogEntry{
				Round: b.CurrentRound,
				Turn:  b.CurrentTurn,
				Text:  fmt.Sprintf("%s 被击败了", target.Character.GetNameWithType()),
			})
		}
	}

	if skill.IsHealingSkill() {
		// Apply healing
		healed := target.Character.Heal(skill.HealAmount)

		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  fmt.Sprintf("%s 恢复了 %d 点生命值", target.Character.GetNameWithType(), healed),
		})
	}

	// Apply status effects for certain skills
	switch skill.Type {
	case combat.ShieldBash:
		// 25% chance to stun
		if rand.Intn(100) < 25 {
			b.applyEffect(target, EffectStun, 1, 0)
			b.Log = append(b.Log, BattleLogEntry{
				Round: b.CurrentRound,
				Turn:  b.CurrentTurn,
				Text:  fmt.Sprintf("%s 被眩晕了", target.Character.GetNameWithType()),
			})
		}
	case combat.Frostbolt:
		// 30% chance to slow
		if rand.Intn(100) < 30 {
			b.applyEffect(target, EffectSlow, 2, -2) // Reduce speed by 2 for 2 rounds
			b.Log = append(b.Log, BattleLogEntry{
				Round: b.CurrentRound,
				Turn:  b.CurrentTurn,
				Text:  fmt.Sprintf("%s 被减速了", target.Character.GetNameWithType()),
			})
		}
	}
}

// applyEffect applies a status effect to a participant
func (b *Battle) applyEffect(target *BattleParticipant, effectType EffectType, duration int, value int) {
	effect := BattleEffect{
		Type:      effectType,
		Duration:  duration,
		Remaining: duration,
		Value:     value,
	}

	target.Effects = append(target.Effects, effect)
}

// isStunned checks if a participant is stunned
func (b *Battle) isStunned(participant *BattleParticipant) bool {
	for _, effect := range participant.Effects {
		if effect.Type == EffectStun && effect.Remaining > 0 {
			return true
		}
	}
	return false
}

// getRandomValidTarget gets a random valid target for a participant
func (b *Battle) getRandomValidTarget(participant *BattleParticipant) (*BattleParticipant, error) {
	// Get all alive participants
	var validTargets []*BattleParticipant
	for _, p := range b.Participants {
		if p.IsAlive {
			validTargets = append(validTargets, p)
		}
	}

	if len(validTargets) == 0 {
		return nil, fmt.Errorf("没有有效的目标")
	}

	// If only one target, return it
	if len(validTargets) == 1 {
		return validTargets[0], nil
	}

	// Get a random target
	return validTargets[rand.Intn(len(validTargets))], nil
}

// nextTurn moves to the next turn
func (b *Battle) nextTurn() {
	b.CurrentTurn++

	// If we've gone through all participants, start a new round
	if b.CurrentTurn >= len(b.TurnOrder) {
		b.CurrentRound++
		b.CurrentTurn = 0

		// Update effects
		b.updateEffects()

		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  0,
			Text:  fmt.Sprintf("第%d回合开始", b.CurrentRound),
		})
	}
}

// updateEffects updates status effects at the end of each round
func (b *Battle) updateEffects() {
	for _, participant := range b.Participants {
		// Update effects for this participant
		var updatedEffects []BattleEffect
		for _, effect := range participant.Effects {
			effect.Remaining--

			if effect.Remaining > 0 {
				// Effect still active
				updatedEffects = append(updatedEffects, effect)

				// Apply effect values (e.g., damage over time)
				if effect.Type == EffectPoison {
					damage := effect.Value
					actualDamage := participant.Character.TakeDamage(damage)
					b.Log = append(b.Log, BattleLogEntry{
						Round: b.CurrentRound,
						Turn:  0,
						Text:  fmt.Sprintf("%s 受到 %d 点毒素伤害", participant.Character.GetNameWithType(), actualDamage),
					})

					// Check if participant is defeated
					if !participant.Character.IsAlive() {
						participant.IsAlive = false
						b.Log = append(b.Log, BattleLogEntry{
							Round: b.CurrentRound,
							Turn:  0,
							Text:  fmt.Sprintf("%s 被毒素击败了", participant.Character.GetNameWithType()),
						})
					}
				}
			} else {
				// Effect expired
				effectName := ""
				switch effect.Type {
				case EffectStun:
					effectName = "眩晕"
				case EffectPoison:
					effectName = "中毒"
				case EffectSlow:
					effectName = "减速"
				case EffectBuff:
					effectName = "增益"
				}

				if effectName != "" {
					b.Log = append(b.Log, BattleLogEntry{
						Round: b.CurrentRound,
						Turn:  0,
						Text:  fmt.Sprintf("%s 的 %s 效果已消失", participant.Character.GetNameWithType(), effectName),
					})
				}
			}
		}

		participant.Effects = updatedEffects
	}
}

// checkBattleEnd checks if the battle has ended
func (b *Battle) checkBattleEnd() bool {
	// Count alive players and enemies
	alivePlayers := 0
	aliveEnemies := 0

	for _, p := range b.Participants {
		if p.IsAlive {
			if p.IsPlayer {
				alivePlayers++
			} else {
				aliveEnemies++
			}
		}
	}

	// Determine battle outcome
	if alivePlayers == 0 && aliveEnemies == 0 {
		b.State = BattleStateDraw
		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  "战斗结束！平局",
		})
		return true
	} else if alivePlayers == 0 {
		b.State = BattleStateEnemiesWon
		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  "战斗结束！敌人获胜",
		})
		return true
	} else if aliveEnemies == 0 {
		b.State = BattleStatePlayerWon
		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  "战斗结束！玩家获胜",
		})

		// Calculate rewards
		b.calculateRewards()
		return true
	}

	return false
}

// calculateRewards calculates rewards for winning the battle
func (b *Battle) calculateRewards() {
	totalExp := 0
	totalGold := 0

	// Calculate rewards based on defeated enemies
	for _, p := range b.Participants {
		if !p.IsPlayer && !p.IsAlive {
			// Reward based on enemy level
			totalExp += p.Character.Level * 10
			totalGold += p.Character.Level * 5
		}
	}

	b.Reward.Exp = totalExp
	b.Reward.Gold = totalGold

	// Log rewards
	if totalExp > 0 || totalGold > 0 {
		b.Log = append(b.Log, BattleLogEntry{
			Round: b.CurrentRound,
			Turn:  b.CurrentTurn,
			Text:  fmt.Sprintf("获得经验: %d, 金币: %d", totalExp, totalGold),
		})
	}
}

// GetLog returns the battle log
func (b *Battle) GetLog() []BattleLogEntry {
	return b.Log
}

// GetState returns the battle state
func (b *Battle) GetState() BattleState {
	return b.State
}

// GetReward returns the battle reward
func (b *Battle) GetReward() BattleReward {
	return b.Reward
}

// GetAliveParticipants returns all alive participants
func (b *Battle) GetAliveParticipants() []*BattleParticipant {
	var alive []*BattleParticipant
	for _, p := range b.Participants {
		if p.IsAlive {
			alive = append(alive, p)
		}
	}
	return alive
}
