package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xieyx/go-game-server-by-ai/internal/battle"
	"github.com/xieyx/go-game-server-by-ai/internal/character"
	"github.com/xieyx/go-game-server-by-ai/internal/combat"
)

func main() {
	fmt.Println("欢迎来到回合制战斗游戏演示！")
	fmt.Println("==============================")

	// Create a battle
	b := battle.NewBattle("demo_battle")

	// Create player character
	fmt.Println("创建玩家角色...")
	player := character.NewCharacter("player_1", "勇士", character.Warrior)
	fmt.Printf("玩家角色: %s (等级: %d, HP: %d/%d, MP: %d/%d)\n",
		player.GetNameWithType(), player.Level, player.HP, player.MaxHP, player.MP, player.MaxMP)

	// Create enemy characters
	fmt.Println("\n创建敌人角色...")
	enemy1 := character.NewCharacter("enemy_1", "哥布林", character.Warrior)
	enemy2 := character.NewCharacter("enemy_2", "兽人", character.Warrior)

	fmt.Printf("敌人1: %s (等级: %d, HP: %d/%d, MP: %d/%d)\n",
		enemy1.GetNameWithType(), enemy1.Level, enemy1.HP, enemy1.MaxHP, enemy1.MP, enemy1.MaxMP)
	fmt.Printf("敌人2: %s (等级: %d, HP: %d/%d, MP: %d/%d)\n",
		enemy2.GetNameWithType(), enemy2.Level, enemy2.HP, enemy2.MaxHP, enemy2.MP, enemy2.MaxMP)

	// Add participants to battle
	b.AddParticipant(player, true)
	b.AddParticipant(enemy1, false)
	b.AddParticipant(enemy2, false)

	// Start battle
	fmt.Println("\n开始战斗！")
	b.Start()

	// Print initial battle log
	printBattleLog(b.GetLog())

	// Get player's battle participant
	playerParticipant := b.Participants[0]

	// Battle loop
	reader := bufio.NewReader(os.Stdin)
	for b.GetState() == battle.BattleStateInProgress {
		fmt.Println("\n" + strings.Repeat("-", 30))

		// Show current participant
		current := b.GetCurrentParticipant()
		fmt.Printf("当前回合: %s\n", current.Character.GetNameWithType())

		// If it's player's turn, let them choose action
		if current == playerParticipant {
			// Show player's skills
			skills := combat.GetCharacterSkills(player.Type)
			fmt.Println("\n可用技能:")
			for i, skill := range skills {
				fmt.Printf("%d. %s (消耗MP: %d)\n", i+1, skill.Name, skill.MPCost)
			}

			// Get player's skill choice
			fmt.Print("\n请选择技能 (输入数字): ")
			skillChoiceStr, _ := reader.ReadString('\n')
			skillChoiceStr = strings.TrimSpace(skillChoiceStr)
			skillChoice, err := strconv.Atoi(skillChoiceStr)

			if err != nil || skillChoice < 1 || skillChoice > len(skills) {
				fmt.Println("无效选择，使用普通攻击")
				skillChoice = 1
			}

			selectedSkill := skills[skillChoice-1]

			// Check if player can use the skill
			if !selectedSkill.CanUse(player) {
				fmt.Printf("MP不足，无法使用 %s\n", selectedSkill.Name)
				selectedSkill = combat.GetBasicAttack()
			}

			// Select skill
			b.SelectSkill(current, selectedSkill)

			// Show targets
			fmt.Println("\n选择目标:")
			aliveParticipants := b.GetAliveParticipants()
			targetOptions := make([]*battle.BattleParticipant, 0)

			for i, p := range aliveParticipants {
				if p.IsAlive {
					targetOptions = append(targetOptions, p)
					fmt.Printf("%d. %s (HP: %d/%d)\n", i+1, p.Character.GetNameWithType(), p.Character.HP, p.Character.MaxHP)
				}
			}

			// Get player's target choice
			fmt.Print("\n请选择目标 (输入数字): ")
			targetChoiceStr, _ := reader.ReadString('\n')
			targetChoiceStr = strings.TrimSpace(targetChoiceStr)
			targetChoice, err := strconv.Atoi(targetChoiceStr)

			if err != nil || targetChoice < 1 || targetChoice > len(targetOptions) {
				fmt.Println("无效选择，选择第一个目标")
				targetChoice = 1
			}

			selectedTarget := targetOptions[targetChoice-1]

			// Select target
			b.SelectTarget(current, selectedTarget)
		}

		// Execute turn
		b.ExecuteTurn()

		// Print battle log
		log := b.GetLog()
		if len(log) > 0 {
			lastEntry := log[len(log)-1]
			fmt.Printf("[%d-%d] %s\n", lastEntry.Round, lastEntry.Turn, lastEntry.Text)
		}

		// Show current status
		fmt.Println("\n当前状态:")
		for _, p := range b.GetAliveParticipants() {
			fmt.Printf("%s: HP %d/%d, MP %d/%d\n",
				p.Character.GetNameWithType(),
				p.Character.HP, p.Character.MaxHP,
				p.Character.MP, p.Character.MaxMP)
		}

		// Pause for user to read
		if current == playerParticipant {
			fmt.Print("\n按回车键继续...")
			reader.ReadString('\n')
		}
	}

	// Print final battle log
	fmt.Println("\n" + strings.Repeat("=", 30))
	fmt.Println("战斗结束！")
	printBattleLog(b.GetLog())

	// Show rewards if player won
	if b.GetState() == battle.BattleStatePlayerWon {
		reward := b.GetReward()
		fmt.Printf("\n恭喜获胜！获得奖励:\n")
		fmt.Printf("经验: %d\n", reward.Exp)
		fmt.Printf("金币: %d\n", reward.Gold)

		// Apply rewards to player
		if reward.Exp > 0 {
			fmt.Printf("\n%s 获得了 %d 点经验", player.GetNameWithType(), reward.Exp)
			leveledUp := player.GainExp(reward.Exp)
			if leveledUp {
				fmt.Printf("\n%s 升级了！现在是 %d 级", player.GetNameWithType(), player.Level)
			}
		}
	}
}

func printBattleLog(log []battle.BattleLogEntry) {
	fmt.Println("\n战斗日志:")
	for _, entry := range log {
		fmt.Printf("[%d-%d] %s\n", entry.Round, entry.Turn, entry.Text)
	}
}
