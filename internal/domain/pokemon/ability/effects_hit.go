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
