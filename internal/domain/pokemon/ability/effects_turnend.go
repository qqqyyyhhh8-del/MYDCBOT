package ability

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// ============================================
// 回合结束类特性
// ============================================

// SpeedBoostEffect 加速特性
type SpeedBoostEffect struct {
	BaseEffect
}

func (e *SpeedBoostEffect) GetAbilityID() int {
	return 3
}

func (e *SpeedBoostEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnTurnEnd}
}

func (e *SpeedBoostEffect) OnTurnEnd(self Battler, ctx *BattleContext) *TurnEndResult {
	return &TurnEndResult{
		Messages:   []string{"⚡ 加速提升了速度！"},
		StatBoosts: map[string]int{"speed": 1},
	}
}

// ============================================
// 速度修正类特性
// ============================================

// SwiftSwimEffect 悠游自如特性
type SwiftSwimEffect struct {
	BaseEffect
}

func (e *SwiftSwimEffect) GetAbilityID() int {
	return 33
}

func (e *SwiftSwimEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnSpeedCalc}
}

func (e *SwiftSwimEffect) OnSpeedCalc(self Battler, ctx *BattleContext) *SpeedModifier {
	if ctx != nil && ctx.Weather == valueobject.WeatherRain {
		return &SpeedModifier{Multiplier: 2.0}
	}
	return nil
}

// ChlorophyllEffect 叶绿素特性
type ChlorophyllEffect struct {
	BaseEffect
}

func (e *ChlorophyllEffect) GetAbilityID() int {
	return 34
}

func (e *ChlorophyllEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnSpeedCalc}
}

func (e *ChlorophyllEffect) OnSpeedCalc(self Battler, ctx *BattleContext) *SpeedModifier {
	if ctx != nil && ctx.Weather == valueobject.WeatherSun {
		return &SpeedModifier{Multiplier: 2.0}
	}
	return nil
}

// SandRushEffect 拨沙特性
type SandRushEffect struct {
	BaseEffect
}

func (e *SandRushEffect) GetAbilityID() int {
	return 146
}

func (e *SandRushEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnSpeedCalc}
}

func (e *SandRushEffect) OnSpeedCalc(self Battler, ctx *BattleContext) *SpeedModifier {
	if ctx != nil && ctx.Weather == valueobject.WeatherSand {
		return &SpeedModifier{Multiplier: 2.0}
	}
	return nil
}

// SlushRushEffect 拨雪特性
type SlushRushEffect struct {
	BaseEffect
}

func (e *SlushRushEffect) GetAbilityID() int {
	return 202
}

func (e *SlushRushEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnSpeedCalc}
}

func (e *SlushRushEffect) OnSpeedCalc(self Battler, ctx *BattleContext) *SpeedModifier {
	if ctx != nil && ctx.Weather == valueobject.WeatherHail {
		return &SpeedModifier{Multiplier: 2.0}
	}
	return nil
}

// ============================================
// 优先度修正类特性
// ============================================

// PranksterEffect 恶作剧之心特性
type PranksterEffect struct {
	BaseEffect
}

func (e *PranksterEffect) GetAbilityID() int {
	return 158
}

func (e *PranksterEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnPriorityCalc}
}

func (e *PranksterEffect) OnPriorityCalc(self Battler, move Move, ctx *BattleContext) *PriorityModifier {
	if move.GetCategory() == "status" {
		return &PriorityModifier{
			Bonus:     1,
			Condition: true,
		}
	}
	return nil
}

// GaleWingsEffect 疾风之翼特性
type GaleWingsEffect struct {
	BaseEffect
}

func (e *GaleWingsEffect) GetAbilityID() int {
	return 177
}

func (e *GaleWingsEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnPriorityCalc}
}

func (e *GaleWingsEffect) OnPriorityCalc(self Battler, move Move, ctx *BattleContext) *PriorityModifier {
	if move.GetType() == valueobject.TypeFlying && self.GetCurrentHP() == self.GetMaxHP() {
		return &PriorityModifier{
			Bonus:     1,
			Condition: true,
		}
	}
	return nil
}

// ============================================
// 击倒触发类特性
// ============================================

// MoxieEffect 自信过剩特性
type MoxieEffect struct {
	BaseEffect
}

func (e *MoxieEffect) GetAbilityID() int {
	return 153
}

func (e *MoxieEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnKO}
}

func (e *MoxieEffect) OnKO(self Battler, target Battler, ctx *BattleContext) *TurnEndResult {
	return &TurnEndResult{
		Messages:   []string{"💪 自信过剩提升了攻击！"},
		StatBoosts: map[string]int{"atk": 1},
	}
}
