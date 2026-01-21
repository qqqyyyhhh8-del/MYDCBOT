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

// VoltAbsorbEffect 蓄电特性
type VoltAbsorbEffect struct {
	BaseEffect
}

func (e *VoltAbsorbEffect) GetAbilityID() int {
	return 10
}

func (e *VoltAbsorbEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *VoltAbsorbEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeElectric {
		mod := NewDamageModifier()
		mod.Immune = true
		mod.HealPercent = 25 // 回复25% HP
		return mod
	}
	return nil
}

// WaterAbsorbEffect 储水特性
type WaterAbsorbEffect struct {
	BaseEffect
}

func (e *WaterAbsorbEffect) GetAbilityID() int {
	return 11
}

func (e *WaterAbsorbEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *WaterAbsorbEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeWater {
		mod := NewDamageModifier()
		mod.Immune = true
		mod.HealPercent = 25 // 回复25% HP
		return mod
	}
	return nil
}

// FlashFireEffect 引火特性
type FlashFireEffect struct {
	BaseEffect
}

func (e *FlashFireEffect) GetAbilityID() int {
	return 18
}

func (e *FlashFireEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *FlashFireEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeFire {
		mod := NewDamageModifier()
		mod.Immune = true
		// 注意：火属性威力提升需要在攻击时处理
		return mod
	}
	return nil
}

// IronFistEffect 铁拳特性
type IronFistEffect struct {
	BaseEffect
}

func (e *IronFistEffect) GetAbilityID() int {
	return 89
}

func (e *IronFistEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *IronFistEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.IsPunch() {
		mod := NewDamageModifier()
		mod.PowerMod = 1.2
		return mod
	}
	return nil
}

// SniperEffect 狙击手特性
type SniperEffect struct {
	BaseEffect
}

func (e *SniperEffect) GetAbilityID() int {
	return 97
}

func (e *SniperEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SniperEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	// 会心一击时威力额外提升（从1.5倍到2.25倍）
	mod := NewDamageModifier()
	mod.CritMod = 2.25 / 1.5
	return mod
}

// FurCoatEffect 毛皮大衣特性
type FurCoatEffect struct {
	BaseEffect
}

func (e *FurCoatEffect) GetAbilityID() int {
	return 169
}

func (e *FurCoatEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *FurCoatEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetCategory() == "physical" {
		mod := NewDamageModifier()
		mod.DamageMod = 0.5
		return mod
	}
	return nil
}

// SwarmEffect 虫之预感特性
type SwarmEffect struct {
	BaseEffect
}

func (e *SwarmEffect) GetAbilityID() int {
	return 68
}

func (e *SwarmEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SwarmEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeBug && self.GetHPPercent() <= 33.3 {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// SandForceEffect 沙之力特性
type SandForceEffect struct {
	BaseEffect
}

func (e *SandForceEffect) GetAbilityID() int {
	return 159
}

func (e *SandForceEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SandForceEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if ctx != nil && ctx.Weather == valueobject.WeatherSand {
		moveType := move.GetType()
		if moveType == valueobject.TypeRock || moveType == valueobject.TypeGround || moveType == valueobject.TypeSteel {
			mod := NewDamageModifier()
			mod.PowerMod = 1.3
			return mod
		}
	}
	return nil
}

// HeatproofEffect 耐热特性
type HeatproofEffect struct {
	BaseEffect
}

func (e *HeatproofEffect) GetAbilityID() int {
	return 85
}

func (e *HeatproofEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *HeatproofEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeFire {
		mod := NewDamageModifier()
		mod.DamageMod = 0.5
		return mod
	}
	return nil
}

// DrySkinEffect 干燥皮肤特性
type DrySkinEffect struct {
	BaseEffect
}

func (e *DrySkinEffect) GetAbilityID() int {
	return 87
}

func (e *DrySkinEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc, TriggerOnTurnEnd}
}

func (e *DrySkinEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeWater {
		mod := NewDamageModifier()
		mod.Immune = true
		mod.HealPercent = 25
		return mod
	}
	if move.GetType() == valueobject.TypeFire {
		mod := NewDamageModifier()
		mod.DamageMod = 1.25
		return mod
	}
	return nil
}

// StormDrainEffect 引水特性
type StormDrainEffect struct {
	BaseEffect
}

func (e *StormDrainEffect) GetAbilityID() int {
	return 114
}

func (e *StormDrainEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *StormDrainEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeWater {
		mod := NewDamageModifier()
		mod.Immune = true
		// 注意：特攻提升需要在其他地方处理
		return mod
	}
	return nil
}

// SapSipperEffect 食草特性
type SapSipperEffect struct {
	BaseEffect
}

func (e *SapSipperEffect) GetAbilityID() int {
	return 157
}

func (e *SapSipperEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SapSipperEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeGrass {
		mod := NewDamageModifier()
		mod.Immune = true
		// 注意：攻击提升需要在其他地方处理
		return mod
	}
	return nil
}

// MotorDriveEffect 电气引擎特性
type MotorDriveEffect struct {
	BaseEffect
}

func (e *MotorDriveEffect) GetAbilityID() int {
	return 78
}

func (e *MotorDriveEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *MotorDriveEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeElectric {
		mod := NewDamageModifier()
		mod.Immune = true
		// 注意：速度提升需要在其他地方处理
		return mod
	}
	return nil
}

// SolidRockEffect 坚硬岩石特性
type SolidRockEffect struct {
	BaseEffect
}

func (e *SolidRockEffect) GetAbilityID() int {
	return 116
}

func (e *SolidRockEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SolidRockEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	effectiveness := valueobject.GetEffectiveness(move.GetType(), self.GetTypes())
	if effectiveness > 1.0 {
		mod := NewDamageModifier()
		mod.DamageMod = 0.75
		return mod
	}
	return nil
}

// FilterEffect 过滤特性（与坚硬岩石相同）
type FilterEffect struct {
	BaseEffect
}

func (e *FilterEffect) GetAbilityID() int {
	return 111
}

func (e *FilterEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *FilterEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	effectiveness := valueobject.GetEffectiveness(move.GetType(), self.GetTypes())
	if effectiveness > 1.0 {
		mod := NewDamageModifier()
		mod.DamageMod = 0.75
		return mod
	}
	return nil
}

// PrismArmorEffect 棱镜装甲特性
type PrismArmorEffect struct {
	BaseEffect
}

func (e *PrismArmorEffect) GetAbilityID() int {
	return 232
}

func (e *PrismArmorEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *PrismArmorEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	effectiveness := valueobject.GetEffectiveness(move.GetType(), self.GetTypes())
	if effectiveness > 1.0 {
		mod := NewDamageModifier()
		mod.DamageMod = 0.75
		return mod
	}
	return nil
}

// TintedLensEffect 有色眼镜特性
type TintedLensEffect struct {
	BaseEffect
}

func (e *TintedLensEffect) GetAbilityID() int {
	return 110
}

func (e *TintedLensEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *TintedLensEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	effectiveness := valueobject.GetEffectiveness(move.GetType(), target.GetTypes())
	if effectiveness < 1.0 && effectiveness > 0 {
		mod := NewDamageModifier()
		mod.DamageMod = 2.0 // 效果不好变为正常威力
		return mod
	}
	return nil
}

// NeuroforceEffect 脑核之力特性
type NeuroforceEffect struct {
	BaseEffect
}

func (e *NeuroforceEffect) GetAbilityID() int {
	return 233
}

func (e *NeuroforceEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *NeuroforceEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	effectiveness := valueobject.GetEffectiveness(move.GetType(), target.GetTypes())
	if effectiveness > 1.0 {
		mod := NewDamageModifier()
		mod.DamageMod = 1.25
		return mod
	}
	return nil
}

// RecklessEffect 舍身特性
type RecklessEffect struct {
	BaseEffect
}

func (e *RecklessEffect) GetAbilityID() int {
	return 120
}

func (e *RecklessEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *RecklessEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.IsRecoil() {
		mod := NewDamageModifier()
		mod.PowerMod = 1.2
		return mod
	}
	return nil
}

// MegaLauncherEffect 超级发射器特性
type MegaLauncherEffect struct {
	BaseEffect
}

func (e *MegaLauncherEffect) GetAbilityID() int {
	return 178
}

func (e *MegaLauncherEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *MegaLauncherEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.IsPulse() {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}

// SteelworkerEffect 钢能力者特性
type SteelworkerEffect struct {
	BaseEffect
}

func (e *SteelworkerEffect) GetAbilityID() int {
	return 200
}

func (e *SteelworkerEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnDamageCalc}
}

func (e *SteelworkerEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	if move.GetType() == valueobject.TypeSteel {
		mod := NewDamageModifier()
		mod.PowerMod = 1.5
		return mod
	}
	return nil
}
