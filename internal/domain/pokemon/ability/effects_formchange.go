package ability

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// ============================================
// 形态变化类特性
// ============================================

// BattleBondEffect 羁绊变身特性 (甲贺忍蛙专属)
// 击倒对手后变身为小智版甲贺忍蛙
type BattleBondEffect struct {
	BaseEffect
}

func (e *BattleBondEffect) GetAbilityID() int {
	return 210 // Battle Bond
}

func (e *BattleBondEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnKO}
}

func (e *BattleBondEffect) OnFormChange(self Battler, target Battler, ctx *BattleContext) *FormChangeResult {
	// 只有甲贺忍蛙可以触发 (ID: 658)
	// 检查是否已经变身过
	if self.HasVolatile("battle_bond_transformed") {
		return nil
	}

	return &FormChangeResult{
		Triggered:   true,
		NewFormID:   10116, // 小智版甲贺忍蛙的形态ID
		NewFormName: "甲贺忍蛙(小智版)",
		NewTypes:    []valueobject.PokeType{valueobject.TypeWater, valueobject.TypeDark},
		StatBoosts: map[string]int{
			"atk":   50,
			"spatk": 50,
			"speed": 10,
		},
		Messages: []string{
			"🌟 甲贺忍蛙与训练师的羁绊达到了顶点！",
			"✨ 甲贺忍蛙变身为小智版甲贺忍蛙！",
		},
		SpriteURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/10116.gif",
		RevertOnExit:  true,
		RevertOnFaint: true,
	}
}

// ZenModeEffect 达摩模式特性 (达摩狒狒专属)
// HP<=50%时变成达摩模式
type ZenModeEffect struct {
	BaseEffect
}

func (e *ZenModeEffect) GetAbilityID() int {
	return 161 // Zen Mode
}

func (e *ZenModeEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnTurnEnd}
}

func (e *ZenModeEffect) OnFormChange(self Battler, target Battler, ctx *BattleContext) *FormChangeResult {
	// HP<=50%时变成达摩模式
	if self.GetHPPercent() <= 50 && !self.HasVolatile("zen_mode") {
		return &FormChangeResult{
			Triggered:   true,
			NewFormID:   10017, // 达摩模式
			NewFormName: "达摩狒狒(达摩模式)",
			NewTypes:    []valueobject.PokeType{valueobject.TypeFire, valueobject.TypePsychic},
			StatBoosts: map[string]int{
				"atk":   -60,
				"spatk": 90,
				"speed": 55,
			},
			Messages:      []string{"🧘 达摩狒狒进入了达摩模式！"},
			SpriteURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/555-zen.gif",
			RevertOnExit:  true,
			RevertOnFaint: true,
		}
	}
	return nil
}

// PowerConstructEffect 群聚变形特性 (基格尔德专属)
// HP<=50%时变成完全体形态
type PowerConstructEffect struct {
	BaseEffect
}

func (e *PowerConstructEffect) GetAbilityID() int {
	return 211 // Power Construct
}

func (e *PowerConstructEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnTurnEnd}
}

func (e *PowerConstructEffect) OnFormChange(self Battler, target Battler, ctx *BattleContext) *FormChangeResult {
	// HP<=50%时变成完全体形态
	if self.GetHPPercent() <= 50 && !self.HasVolatile("power_construct_complete") {
		return &FormChangeResult{
			Triggered:   true,
			NewFormID:   10118, // 完全体形态
			NewFormName: "基格尔德(完全体)",
			NewTypes:    []valueobject.PokeType{valueobject.TypeDragon, valueobject.TypeGround},
			StatBoosts: map[string]int{
				"hp": 108, // 完全体HP大幅提升
			},
			Messages: []string{
				"🐉 基格尔德召集了所有细胞！",
				"✨ 基格尔德变成了完全体形态！",
			},
			SpriteURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/718-complete.gif",
			RevertOnExit:  true,
			RevertOnFaint: false, // 完全体不会因濒死恢复
		}
	}
	return nil
}

// StanceChangeEffect 战斗切换特性 (坚盾剑怪专属)
// 使用攻击技能时变成剑形态，使用王者盾牌时变成盾形态
type StanceChangeEffect struct {
	BaseEffect
}

func (e *StanceChangeEffect) GetAbilityID() int {
	return 176 // Stance Change
}

func (e *StanceChangeEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnMoveUse}
}

func (e *StanceChangeEffect) OnFormChange(self Battler, target Battler, ctx *BattleContext) *FormChangeResult {
	return nil
}

// GetBladeFormChange 获取剑形态变化数据
func (e *StanceChangeEffect) GetBladeFormChange() *FormChangeResult {
	return &FormChangeResult{
		Triggered:   true,
		NewFormID:   10026, // 剑形态
		NewFormName: "坚盾剑怪(剑形态)",
		NewTypes:    []valueobject.PokeType{valueobject.TypeSteel, valueobject.TypeGhost},
		StatBoosts: map[string]int{
			"atk":   100,
			"def":   -100,
			"spatk": 100,
			"spdef": -100,
		},
		Messages:      []string{"⚔️ 坚盾剑怪变成了剑形态！"},
		SpriteURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/681-blade.gif",
		RevertOnExit:  true,
		RevertOnFaint: true,
	}
}

// GetShieldFormChange 获取盾形态变化数据
func (e *StanceChangeEffect) GetShieldFormChange() *FormChangeResult {
	return &FormChangeResult{
		Triggered:   true,
		NewFormID:   681, // 盾形态（默认）
		NewFormName: "坚盾剑怪(盾形态)",
		NewTypes:    []valueobject.PokeType{valueobject.TypeSteel, valueobject.TypeGhost},
		StatBoosts: map[string]int{
			"atk":   -100,
			"def":   100,
			"spatk": -100,
			"spdef": 100,
		},
		Messages:      []string{"🛡️ 坚盾剑怪变成了盾形态！"},
		SpriteURL:     "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/681.gif",
		RevertOnExit:  true,
		RevertOnFaint: true,
	}
}
