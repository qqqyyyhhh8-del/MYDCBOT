package ability

// ============================================
// 状态免疫类特性
// ============================================

// ImmunityEffect 免疫特性
type ImmunityEffect struct {
	BaseEffect
}

func (e *ImmunityEffect) GetAbilityID() int {
	return 17
}

func (e *ImmunityEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *ImmunityEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "中毒" || status == "剧毒" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 免疫特性阻止了中毒！",
		}
	}
	return nil
}

// InnerFocusEffect 精神力特性
type InnerFocusEffect struct {
	BaseEffect
}

func (e *InnerFocusEffect) GetAbilityID() int {
	return 39
}

func (e *InnerFocusEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *InnerFocusEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "畏缩" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 精神力阻止了畏缩！",
		}
	}
	return nil
}

// LimberEffect 柔软特性
type LimberEffect struct {
	BaseEffect
}

func (e *LimberEffect) GetAbilityID() int {
	return 7
}

func (e *LimberEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *LimberEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "麻痹" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 柔软特性阻止了麻痹！",
		}
	}
	return nil
}

// InsomniaEffect 不眠特性
type InsomniaEffect struct {
	BaseEffect
}

func (e *InsomniaEffect) GetAbilityID() int {
	return 15
}

func (e *InsomniaEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *InsomniaEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "睡眠" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 不眠特性阻止了睡眠！",
		}
	}
	return nil
}

// VitalSpiritEffect 干劲特性
type VitalSpiritEffect struct {
	BaseEffect
}

func (e *VitalSpiritEffect) GetAbilityID() int {
	return 72
}

func (e *VitalSpiritEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *VitalSpiritEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "睡眠" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 干劲特性阻止了睡眠！",
		}
	}
	return nil
}

// MagmaArmorEffect 熔岩铠甲特性
type MagmaArmorEffect struct {
	BaseEffect
}

func (e *MagmaArmorEffect) GetAbilityID() int {
	return 40
}

func (e *MagmaArmorEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *MagmaArmorEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "冰冻" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 熔岩铠甲阻止了冰冻！",
		}
	}
	return nil
}

// WaterVeilEffect 水幕特性
type WaterVeilEffect struct {
	BaseEffect
}

func (e *WaterVeilEffect) GetAbilityID() int {
	return 41
}

func (e *WaterVeilEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *WaterVeilEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "灼伤" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 水幕特性阻止了灼伤！",
		}
	}
	return nil
}

// OwnTempoEffect 我行我素特性
type OwnTempoEffect struct {
	BaseEffect
}

func (e *OwnTempoEffect) GetAbilityID() int {
	return 20
}

func (e *OwnTempoEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *OwnTempoEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "混乱" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 我行我素阻止了混乱！",
		}
	}
	return nil
}

// ObliviousEffect 迟钝特性
type ObliviousEffect struct {
	BaseEffect
}

func (e *ObliviousEffect) GetAbilityID() int {
	return 12
}

func (e *ObliviousEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnStatusApply}
}

func (e *ObliviousEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	if status == "着迷" || status == "挑衅" {
		return &StatusCheckResult{
			Immune:  true,
			Message: "🛡️ 迟钝特性阻止了" + status + "！",
		}
	}
	return nil
}
