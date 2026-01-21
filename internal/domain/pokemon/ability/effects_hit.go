package ability

import (
	"math/rand"
	"time"
)

// ============================================
// 受击触发类特性
// ============================================

// StaticEffect 静电特性
type StaticEffect struct {
	BaseEffect
}

func (e *StaticEffect) GetAbilityID() int {
	return 9
}

func (e *StaticEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *StaticEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() && attacker.GetStatus() == "" {
		return &HitResult{
			ContactEffect: "麻痹",
			ContactChance: 30,
		}
	}
	return nil
}

// CursedBodyEffect 诅咒之躯特性
type CursedBodyEffect struct {
	BaseEffect
}

func (e *CursedBodyEffect) GetAbilityID() int {
	return 130
}

func (e *CursedBodyEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *CursedBodyEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(100) < 30 {
		return &HitResult{
			Messages: []string{"👻 诅咒之躯封印了 " + move.GetName() + "！"},
		}
	}
	return nil
}

// StenchEffect 恶臭特性 - 攻击时有10%几率使对手畏缩
// 注意：恶臭是攻击方特性，这里通过伤害计算触发器实现
type StenchEffect struct {
	BaseEffect
}

func (e *StenchEffect) GetAbilityID() int {
	return 1
}

func (e *StenchEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *StenchEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	// 恶臭特性：攻击时有10%几率使对手畏缩
	// 畏缩效果需要在战斗逻辑中处理，这里只返回标记
	return nil
}

// PoisonPointEffect 毒刺特性
type PoisonPointEffect struct {
	BaseEffect
}

func (e *PoisonPointEffect) GetAbilityID() int {
	return 38
}

func (e *PoisonPointEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *PoisonPointEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() && attacker.GetStatus() == "" {
		return &HitResult{
			ContactEffect: "中毒",
			ContactChance: 30,
		}
	}
	return nil
}

// FlameBodyEffect 火焰之躯特性
type FlameBodyEffect struct {
	BaseEffect
}

func (e *FlameBodyEffect) GetAbilityID() int {
	return 49
}

func (e *FlameBodyEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *FlameBodyEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() && attacker.GetStatus() == "" {
		return &HitResult{
			ContactEffect: "灼伤",
			ContactChance: 30,
		}
	}
	return nil
}

// RoughSkinEffect 粗糙皮肤特性
type RoughSkinEffect struct {
	BaseEffect
}

func (e *RoughSkinEffect) GetAbilityID() int {
	return 24
}

func (e *RoughSkinEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *RoughSkinEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() {
		recoilDamage := attacker.GetMaxHP() / 8
		if recoilDamage < 1 {
			recoilDamage = 1
		}
		return &HitResult{
			RecoilDamage: recoilDamage,
			Messages:     []string{"🦔 粗糙皮肤反弹了伤害！"},
		}
	}
	return nil
}

// EffectSporeEffect 孢子特性
type EffectSporeEffect struct {
	BaseEffect
}

func (e *EffectSporeEffect) GetAbilityID() int {
	return 27
}

func (e *EffectSporeEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *EffectSporeEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() && attacker.GetStatus() == "" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rolls := r.Intn(100)
		if rolls < 10 {
			return &HitResult{
				ContactEffect: "中毒",
				ContactChance: 100,
			}
		} else if rolls < 20 {
			return &HitResult{
				ContactEffect: "麻痹",
				ContactChance: 100,
			}
		} else if rolls < 30 {
			return &HitResult{
				ContactEffect: "睡眠",
				ContactChance: 100,
			}
		}
	}
	return nil
}

// IronBarbsEffect 铁刺特性
type IronBarbsEffect struct {
	BaseEffect
}

func (e *IronBarbsEffect) GetAbilityID() int {
	return 160
}

func (e *IronBarbsEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *IronBarbsEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() {
		recoilDamage := attacker.GetMaxHP() / 8
		if recoilDamage < 1 {
			recoilDamage = 1
		}
		return &HitResult{
			RecoilDamage: recoilDamage,
			Messages:     []string{"🔩 铁刺反弹了伤害！"},
		}
	}
	return nil
}

// CuteCharmEffect 迷人之躯特性
type CuteCharmEffect struct {
	BaseEffect
}

func (e *CuteCharmEffect) GetAbilityID() int {
	return 56
}

func (e *CuteCharmEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *CuteCharmEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() {
		return &HitResult{
			ContactEffect: "着迷",
			ContactChance: 30,
		}
	}
	return nil
}

// GooeyEffect 黏滑特性
type GooeyEffect struct {
	BaseEffect
}

func (e *GooeyEffect) GetAbilityID() int {
	return 183
}

func (e *GooeyEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *GooeyEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() {
		return &HitResult{
			StatChanges: map[string]int{"speed": -1},
			Messages:    []string{"🐌 黏滑降低了对手的速度！"},
		}
	}
	return nil
}

// TanglingHairEffect 卷发特性
type TanglingHairEffect struct {
	BaseEffect
}

func (e *TanglingHairEffect) GetAbilityID() int {
	return 221
}

func (e *TanglingHairEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *TanglingHairEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() {
		return &HitResult{
			StatChanges: map[string]int{"speed": -1},
			Messages:    []string{"💇 卷发降低了对手的速度！"},
		}
	}
	return nil
}

// MummyEffect 木乃伊特性
type MummyEffect struct {
	BaseEffect
}

func (e *MummyEffect) GetAbilityID() int {
	return 152
}

func (e *MummyEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnBeingHit}
}

func (e *MummyEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	if move.IsContact() {
		return &HitResult{
			Messages: []string{"🧟 木乃伊将对手的特性变为木乃伊！"},
		}
	}
	return nil
}
