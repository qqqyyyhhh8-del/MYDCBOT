package ability

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// ============================================
// 计算修正类特性（攻击方）
// ============================================

// HugePowerEffect 大力士特性
type HugePowerEffect struct {
	BaseEffect
}

func (e *HugePowerEffect) GetAbilityID() int {
	return 37
}

func (e *HugePowerEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *HugePowerEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetCategory() == "physical" {
		mod := NewDamageModifier()
		mod.AttackMod = 2.0
		return mod
	}
	return nil
}

// PurePowerEffect 瑜伽之力特性（与大力士相同）
type PurePowerEffect struct {
	BaseEffect
}

func (e *PurePowerEffect) GetAbilityID() int {
	return 74
}

func (e *PurePowerEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *PurePowerEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetCategory() == "physical" {
		mod := NewDamageModifier()
		mod.AttackMod = 2.0
		return mod
	}
	return nil
}

// TechnicianEffect 技术高手特性
type TechnicianEffect struct {
	BaseEffect
}

func (e *TechnicianEffect) GetAbilityID() int {
	return 101
}

func (e *TechnicianEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *TechnicianEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetPower() > 0 && move.GetPower() <= 60 {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// ToughClawsEffect 硬爪特性
type ToughClawsEffect struct {
	BaseEffect
}

func (e *ToughClawsEffect) GetAbilityID() int {
	return 181
}

func (e *ToughClawsEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *ToughClawsEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.IsContact() {
		mod := NewDamageModifier()
		mod.PowerMod = 1.3
		return mod
	}
	return nil
}

// StrongJawEffect 强壮之颚特性
type StrongJawEffect struct {
	BaseEffect
}

func (e *StrongJawEffect) GetAbilityID() int {
	return 173
}

func (e *StrongJawEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *StrongJawEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.IsBite() {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// AdaptabilityEffect 适应力特性
type AdaptabilityEffect struct {
	BaseEffect
}

func (e *AdaptabilityEffect) GetAbilityID() int {
	return 91
}

func (e *AdaptabilityEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *AdaptabilityEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	// 检查是否为本属性技能
	for _, t := range self.GetTypes() {
		if t == move.GetType() {
			mod := NewDamageModifier()
			mod.STABMod = 2.0 / 1.5 // 将STAB从1.5提升到2.0
			return mod
		}
	}
	return nil
}

// SheerForceEffect 强行特性
type SheerForceEffect struct {
	BaseEffect
}

func (e *SheerForceEffect) GetAbilityID() int {
	return 125
}

func (e *SheerForceEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SheerForceEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	// 强行特性：放弃追加效果，威力提升30%
	// 这里简化处理，假设所有有追加效果的技能都适用
	mod := NewDamageModifier()
	mod.PowerMod = 1.3
	return mod
}

// OvergrowEffect 茂盛特性
type OvergrowEffect struct {
	BaseEffect
}

func (e *OvergrowEffect) GetAbilityID() int {
	return 65
}

func (e *OvergrowEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *OvergrowEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeGrass && self.GetHPPercent() <= 33.3 {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// BlazeEffect 猛火特性
type BlazeEffect struct {
	BaseEffect
}

func (e *BlazeEffect) GetAbilityID() int {
	return 66
}

func (e *BlazeEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *BlazeEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeFire && self.GetHPPercent() <= 33.3 {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// TorrentEffect 激流特性
type TorrentEffect struct {
	BaseEffect
}

func (e *TorrentEffect) GetAbilityID() int {
	return 67
}

func (e *TorrentEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *TorrentEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeWater && self.GetHPPercent() <= 33.3 {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// ============================================
// 计算修正类特性（防御方）
// ============================================

// ThickFatEffect 厚脂肪特性
type ThickFatEffect struct {
	BaseEffect
}

func (e *ThickFatEffect) GetAbilityID() int {
	return 47
}

func (e *ThickFatEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *ThickFatEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeFire || move.GetType() == valueobject.TypeIce {
		mod := NewDamageModifier()
		mod.DamageMod = 0.5
		return mod
	}
	return nil
}

// LevitateEffect 飘浮特性
type LevitateEffect struct {
	BaseEffect
}

func (e *LevitateEffect) GetAbilityID() int {
	return 26
}

func (e *LevitateEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *LevitateEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeGround {
		mod := NewDamageModifier()
		mod.Immune = true
		return mod
	}
	return nil
}

// WonderGuardEffect 神奇守护特性
type WonderGuardEffect struct {
	BaseEffect
}

func (e *WonderGuardEffect) GetAbilityID() int {
	return 25
}

func (e *WonderGuardEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *WonderGuardEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	// 只有效果拔群的招式才能命中
	effectiveness := valueobject.GetEffectiveness(move.GetType(), self.GetTypes())
	if effectiveness <= 1.0 {
		mod := NewDamageModifier()
		mod.Immune = true
		return mod
	}
	return nil
}

// MultiscaleEffect 多重鳞片特性
type MultiscaleEffect struct {
	BaseEffect
}

func (e *MultiscaleEffect) GetAbilityID() int {
	return 136
}

func (e *MultiscaleEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *MultiscaleEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if self.GetCurrentHP() == self.GetMaxHP() {
		mod := NewDamageModifier()
		mod.DamageMod = 0.5
		return mod
	}
	return nil
}

// LightningRodEffect 避雷针特性
type LightningRodEffect struct {
	BaseEffect
}

func (e *LightningRodEffect) GetAbilityID() int {
	return 31
}

func (e *LightningRodEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *LightningRodEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeElectric {
		mod := NewDamageModifier()
		mod.Immune = true
		// 注意：特攻提升需要在其他地方处理
		return mod
	}
	return nil
}
