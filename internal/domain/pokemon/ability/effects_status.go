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
