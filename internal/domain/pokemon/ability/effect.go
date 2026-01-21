package ability

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// TriggerType 特性触发时机
type TriggerType string

const (
	TriggerOnEntry        TriggerType = "on_entry"         // 出场时
	TriggerOnDamageCalc   TriggerType = "on_damage_calc"   // 伤害计算时
	TriggerOnBeingHit     TriggerType = "on_being_hit"     // 被攻击时
	TriggerOnTurnEnd      TriggerType = "on_turn_end"      // 回合结束时
	TriggerOnStatusApply  TriggerType = "on_status_apply"  // 状态施加时（用于免疫）
	TriggerOnStatChange   TriggerType = "on_stat_change"   // 能力变化时
	TriggerOnWeatherCheck TriggerType = "on_weather_check" // 天气检查时
	TriggerOnMoveUse      TriggerType = "on_move_use"      // 使用技能时
	TriggerOnKO           TriggerType = "on_ko"            // 击倒对手时
	TriggerOnSpeedCalc    TriggerType = "on_speed_calc"    // 速度计算时
	TriggerOnPriorityCalc TriggerType = "on_priority_calc" // 优先度计算时
)

// BattleContext 战斗上下文，用于特性效果处理
type BattleContext struct {
	Weather       valueobject.Weather // 当前天气
	Terrain       string              // 当前场地
	Turn          int                 // 当前回合
	IsDoubles     bool                // 是否双打
}

// Battler 战斗宝可梦接口（避免循环依赖）
type Battler interface {
	GetAbility() *valueobject.Ability
	GetTypes() []valueobject.PokeType
	GetCurrentHP() int
	GetMaxHP() int
	GetHPPercent() float64
	GetStatus() string
	IsAlive() bool
	ModifyStat(stat string, stages int) (int, bool)
	TakeDamage(damage int) int
	Heal(amount int) int
	SetStatus(status string)
	HasVolatile(status string) bool
	AddVolatile(status string)
	RemoveVolatile(status string)
	GetItem() *valueobject.Item
	IsItemConsumed() bool
	ConsumeItem()
}

// Move 技能接口
type Move interface {
	GetName() string
	GetType() valueobject.PokeType
	GetCategory() string
	GetPower() int
	GetPriority() int
	IsContact() bool
	IsBite() bool
	IsPunch() bool
	IsSound() bool
	IsBullet() bool
	IsRecoil() bool // 是否有反作用力（舍身撞、蛮力等）
	IsPulse() bool  // 是否为波动/波导技能
}

// DamageModifier 伤害修正结果
type DamageModifier struct {
	PowerMod      float64 // 威力修正
	AttackMod     float64 // 攻击修正
	DefenseMod    float64 // 防御修正
	DamageMod     float64 // 最终伤害修正
	STABMod       float64 // 本属性加成修正
	CritMod       float64 // 会心修正
	Immune        bool    // 是否免疫
	HealPercent   float64 // 吸收回复比例（如蓄电、储水）
	TypeOverride  *valueobject.PokeType // 属性覆盖
}

// NewDamageModifier 创建默认伤害修正
func NewDamageModifier() *DamageModifier {
	return &DamageModifier{
		PowerMod:   1.0,
		AttackMod:  1.0,
		DefenseMod: 1.0,
		DamageMod:  1.0,
		STABMod:    1.0,
		CritMod:    1.0,
		Immune:     false,
	}
}

// EntryResult 出场效果结果
type EntryResult struct {
	Messages     []string             // 消息
	WeatherSet   *valueobject.Weather // 设置天气
	StatChanges  map[string]int       // 对手能力变化
}

// HitResult 被击中效果结果
type HitResult struct {
	Messages       []string       // 消息
	DamageReduced  float64        // 伤害减免比例
	ContactEffect  string         // 接触效果（如麻痹）
	ContactChance  int            // 触发几率（百分比）
	RecoilDamage   int            // 反伤伤害
	StatChanges    map[string]int // 能力变化（对攻击方）
}

// TurnEndResult 回合结束效果结果
type TurnEndResult struct {
	Messages     []string       // 消息
	StatBoosts   map[string]int // 能力提升
	Healing      int            // 回复量（旧字段，保留兼容）
	Damage       int            // 伤害量（旧字段，保留兼容）
	HealAmount   int            // 回复量
	DamageAmount int            // 伤害量
	CureStatus   bool           // 是否治愈状态
	CureChance   int            // 治愈几率（百分比）
	NegatePoison bool           // 是否免疫中毒伤害
}

// StatusCheckResult 状态检查结果
type StatusCheckResult struct {
	Immune   bool   // 是否免疫
	Message  string // 消息
}

// SpeedModifier 速度修正结果
type SpeedModifier struct {
	Multiplier float64 // 速度倍率
}

// PriorityModifier 优先度修正结果
type PriorityModifier struct {
	Bonus     int  // 优先度加成
	Condition bool // 是否满足条件
}

// FormChangeResult 形态变化结果
type FormChangeResult struct {
	Triggered     bool                     // 是否触发形态变化
	NewFormID     int                      // 新形态的宝可梦ID（用于获取新数据）
	NewFormName   string                   // 新形态名称
	NewTypes      []valueobject.PokeType   // 新属性（nil表示不变）
	StatBoosts    map[string]int           // 能力值提升（绝对值，非等级）
	Messages      []string                 // 形态变化消息
	SpriteURL     string                   // 新精灵图URL
	RevertOnExit  bool                     // 退场时是否恢复原形态
	RevertOnFaint bool                     // 濒死时是否恢复原形态
}

// Effect 特性效果接口
type Effect interface {
	// GetAbilityID 获取对应的特性ID
	GetAbilityID() int
	
	// GetTriggers 获取触发时机
	GetTriggers() []TriggerType
	
	// OnEntry 出场时触发
	OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult
	
	// OnDamageCalc 伤害计算时触发（作为攻击方）
	OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier
	
	// OnDamageCalcDefender 伤害计算时触发（作为防御方）
	OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier
	
	// OnBeingHit 被攻击后触发
	OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult
	
	// OnTurnEnd 回合结束时触发
	OnTurnEnd(self Battler, ctx *BattleContext) *TurnEndResult
	
	// OnStatusApply 状态施加时检查（用于免疫）
	OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult
	
	// OnSpeedCalc 速度计算时触发
	OnSpeedCalc(self Battler, ctx *BattleContext) *SpeedModifier
	
	// OnPriorityCalc 优先度计算时触发
	OnPriorityCalc(self Battler, move Move, ctx *BattleContext) *PriorityModifier
	
	// OnKO 击倒对手时触发
	OnKO(self Battler, target Battler, ctx *BattleContext) *TurnEndResult
	
	// OnFormChange 击倒对手后检查形态变化（如羁绊进化）
	OnFormChange(self Battler, target Battler, ctx *BattleContext) *FormChangeResult
}

// BaseEffect 基础效果实现（提供默认空实现）
type BaseEffect struct {
	AbilityID int
	Triggers  []TriggerType
}

func (e *BaseEffect) GetAbilityID() int {
	return e.AbilityID
}

func (e *BaseEffect) GetTriggers() []TriggerType {
	return e.Triggers
}

func (e *BaseEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	return nil
}

func (e *BaseEffect) OnDamageCalcAttacker(self Battler, target Battler, move Move, ctx *BattleContext) *DamageModifier {
	return nil
}

func (e *BaseEffect) OnDamageCalcDefender(self Battler, attacker Battler, move Move, ctx *BattleContext) *DamageModifier {
	return nil
}

func (e *BaseEffect) OnBeingHit(self Battler, attacker Battler, move Move, damage int, ctx *BattleContext) *HitResult {
	return nil
}

func (e *BaseEffect) OnTurnEnd(self Battler, ctx *BattleContext) *TurnEndResult {
	return nil
}

func (e *BaseEffect) OnStatusApply(self Battler, status string, ctx *BattleContext) *StatusCheckResult {
	return nil
}

func (e *BaseEffect) OnSpeedCalc(self Battler, ctx *BattleContext) *SpeedModifier {
	return nil
}

func (e *BaseEffect) OnPriorityCalc(self Battler, move Move, ctx *BattleContext) *PriorityModifier {
	return nil
}

func (e *BaseEffect) OnKO(self Battler, target Battler, ctx *BattleContext) *TurnEndResult {
	return nil
}

func (e *BaseEffect) OnFormChange(self Battler, target Battler, ctx *BattleContext) *FormChangeResult {
	return nil
}
