package ability

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// Service 特性效果服务
type Service struct {
	registry *Registry
}

// NewService 创建特性效果服务
func NewService() *Service {
	return &Service{
		registry: GetRegistry(),
	}
}

// TriggerEntry 触发出场特性
func (s *Service) TriggerEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnEntry(self, opponent, ctx)
}

// ApplyAttackerDamageMods 应用攻击方伤害修正
func (s *Service) ApplyAttackerDamageMods(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnDamageCalcAttacker(self, target, move, ctx)
}

// ApplyDefenderDamageMods 应用防御方伤害修正
func (s *Service) ApplyDefenderDamageMods(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnDamageCalcDefender(self, attacker, move, ctx)
}

// TriggerBeingHit 触发被击中特性
func (s *Service) TriggerBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnBeingHit(self, attacker, move, damage, ctx)
}

// TriggerTurnEnd 触发回合结束特性
func (s *Service) TriggerTurnEnd(self Battler, ctx *BattleContext) *TurnEndResult {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnTurnEnd(self, ctx)
}

// CheckStatusImmunity 检查状态免疫
func (s *Service) CheckStatusImmunity(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnStatusApply(self, status, ctx)
}

// GetSpeedModifier 获取速度修正
func (s *Service) GetSpeedModifier(self Battler, ctx *BattleContext) *SpeedModifier {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnSpeedCalc(self, ctx)
}

// GetPriorityModifier 获取优先度修正
func (s *Service) GetPriorityModifier(self Battler, move Move, ctx *BattleContext) *PriorityModifier {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnPriorityCalc(self, move, ctx)
}

// TriggerKO 触发击倒特性
func (s *Service) TriggerKO(self Battler, target Battler, ctx *BattleContext) *TurnEndResult {
	ability := self.GetAbility()
	if ability == nil {
		return nil
	}

	effect := s.registry.Get(ability.ID)
	if effect == nil {
		return nil
	}

	return effect.OnKO(self, target, ctx)
}

// CalculateDamageWithAbilities 计算包含特性效果的伤害
// 返回：最终伤害倍率，是否免疫，消息列表
func (s *Service) CalculateDamageWithAbilities(
	attacker Battler,
	defender Battler,
	move Move,
	ctx *BattleContext,
) (powerMod, atkMod, defMod, damageMod, stabMod float64, immune bool, messages []string) {
	powerMod = 1.0
	atkMod = 1.0
	defMod = 1.0
	damageMod = 1.0
	stabMod = 1.0
	immune = false
	messages = make([]string, 0)

	// 攻击方特性
	atkAbilityMod := s.ApplyAttackerDamageMods(attacker, defender, move, ctx)
	if atkAbilityMod != nil {
		powerMod *= atkAbilityMod.PowerMod
		atkMod *= atkAbilityMod.AttackMod
		damageMod *= atkAbilityMod.DamageMod
		stabMod *= atkAbilityMod.STABMod
	}

	// 防御方特性
	defAbilityMod := s.ApplyDefenderDamageMods(defender, attacker, move, ctx)
	if defAbilityMod != nil {
		if defAbilityMod.Immune {
			immune = true
			if defender.GetAbility() != nil {
				messages = append(messages, "🛡️ "+defender.GetAbility().Name+"使攻击无效！")
			}
			return
		}
		defMod *= defAbilityMod.DefenseMod
		damageMod *= defAbilityMod.DamageMod
	}

	return
}

// GetEffectiveSpeed 获取包含特性效果的有效速度
func (s *Service) GetEffectiveSpeed(self Battler, baseSpeed int, ctx *BattleContext) int {
	speed := baseSpeed

	speedMod := s.GetSpeedModifier(self, ctx)
	if speedMod != nil {
		speed = int(float64(speed) * speedMod.Multiplier)
	}

	return speed
}

// GetEffectivePriority 获取包含特性效果的有效优先度
func (s *Service) GetEffectivePriority(self Battler, move Move, basePriority int, ctx *BattleContext) int {
	priority := basePriority

	priorityMod := s.GetPriorityModifier(self, move, ctx)
	if priorityMod != nil && priorityMod.Condition {
		priority += priorityMod.Bonus
	}

	return priority
}

// ProcessEntryAbility 处理出场特性并返回需要应用的效果
func (s *Service) ProcessEntryAbility(self Battler, opponent Battler, ctx *BattleContext) (messages []string, weather *valueobject.Weather, opponentStatChanges map[string]int) {
	messages = make([]string, 0)

	result := s.TriggerEntry(self, opponent, ctx)
	if result == nil {
		return
	}

	messages = append(messages, result.Messages...)
	weather = result.WeatherSet
	opponentStatChanges = result.StatChanges

	return
}

// ProcessTurnEndAbility 处理回合结束特性
func (s *Service) ProcessTurnEndAbility(self Battler, ctx *BattleContext) (messages []string, statBoosts map[string]int, healing int, damage int) {
	messages = make([]string, 0)

	result := s.TriggerTurnEnd(self, ctx)
	if result == nil {
		return
	}

	messages = append(messages, result.Messages...)
	statBoosts = result.StatBoosts
	healing = result.Healing
	damage = result.Damage

	return
}
