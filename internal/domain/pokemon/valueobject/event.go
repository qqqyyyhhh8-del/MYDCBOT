package valueobject

// BattleEvent 对战事件类型
type BattleEvent string

const (
	// 出场相关
	EventOnEnter      BattleEvent = "on_enter"       // 宝可梦出场时
	EventOnSwitchOut  BattleEvent = "on_switch_out"  // 宝可梦换下时

	// 回合相关
	EventOnTurnStart  BattleEvent = "on_turn_start"  // 回合开始
	EventOnTurnEnd    BattleEvent = "on_turn_end"    // 回合结束

	// 攻击相关
	EventBeforeMove     BattleEvent = "before_move"      // 使用招式前（判断能否行动）
	EventOnMove         BattleEvent = "on_move"          // 使用招式时（变幻自如等）
	EventOnCalcPower    BattleEvent = "on_calc_power"    // 计算威力时
	EventOnCalcAttack   BattleEvent = "on_calc_attack"   // 计算攻击时
	EventOnCalcDefense  BattleEvent = "on_calc_defense"  // 计算防御时
	EventOnCalcSpeed    BattleEvent = "on_calc_speed"    // 计算速度时
	EventOnCalcAccuracy BattleEvent = "on_calc_accuracy" // 计算命中时
	EventOnCalcCrit     BattleEvent = "on_calc_crit"     // 计算会心时
	EventOnCalcDamage   BattleEvent = "on_calc_damage"   // 计算伤害时（最终修正）

	// 受击相关
	EventBeforeHit      BattleEvent = "before_hit"       // 被击中前（免疫判定）
	EventOnHit          BattleEvent = "on_hit"           // 被击中时
	EventAfterHit       BattleEvent = "after_hit"        // 被击中后（反击效果）
	EventOnTakeDamage   BattleEvent = "on_take_damage"   // 受到伤害时

	// 状态相关
	EventOnStatusChange BattleEvent = "on_status_change" // 状态变化时
	EventOnStatChange   BattleEvent = "on_stat_change"   // 能力变化时
	EventOnHPChange     BattleEvent = "on_hp_change"     // HP变化时

	// 天气相关
	EventOnWeatherChange BattleEvent = "on_weather_change" // 天气变化时
	EventOnWeatherDamage BattleEvent = "on_weather_damage" // 天气伤害时

	// 击倒相关
	EventOnKO        BattleEvent = "on_ko"         // 击倒对手时
	EventOnFaint     BattleEvent = "on_faint"      // 自己倒下时
)

// EventContext 事件上下文（传递数据）
type EventContext struct {
	// 基础信息
	Event      BattleEvent
	Source     interface{} // 事件发起者（*Battler）
	Target     interface{} // 事件目标（*Battler）
	
	// 招式相关
	Move       interface{} // 使用的招式（*Move）
	MoveType   PokeType    // 招式属性（可能被特性修改）
	IsContact  bool        // 是否接触招式
	
	// 数值修正
	Power      float64     // 威力修正
	Attack     float64     // 攻击修正
	Defense    float64     // 防御修正
	Speed      float64     // 速度修正
	Accuracy   float64     // 命中修正
	CritStage  int         // 会心等级修正
	Damage     float64     // 伤害修正
	
	// 状态
	Weather    Weather     // 当前天气
	NewStatus  string      // 新状态
	OldStatus  string      // 旧状态
	StatName   string      // 能力名称
	StatStages int         // 能力变化级数
	
	// 控制流
	Cancelled  bool        // 是否取消事件
	Absorbed   bool        // 是否被吸收（蓄电等）
	Immune     bool        // 是否免疫
	Messages   []string    // 生成的消息
}

// NewEventContext 创建事件上下文
func NewEventContext(event BattleEvent) *EventContext {
	return &EventContext{
		Event:    event,
		Power:    1.0,
		Attack:   1.0,
		Defense:  1.0,
		Speed:    1.0,
		Accuracy: 1.0,
		Damage:   1.0,
		Messages: make([]string, 0),
	}
}

// AddMessage 添加消息
func (c *EventContext) AddMessage(msg string) {
	c.Messages = append(c.Messages, msg)
}

// AbilityTrigger 特性触发条件
type AbilityTrigger struct {
	Event    BattleEvent // 触发事件
	Priority int         // 优先度（越高越先执行）
}

// MoveCategory 招式分类
type MoveCategory string

const (
	MoveCategoryPhysical MoveCategory = "物理"
	MoveCategorySpecial  MoveCategory = "特殊"
	MoveCategoryStatus   MoveCategory = "变化"
)

// MoveFlag 招式标签
type MoveFlag string

const (
	FlagContact   MoveFlag = "contact"   // 接触类
	FlagSound     MoveFlag = "sound"     // 声音类
	FlagPunch     MoveFlag = "punch"     // 拳类
	FlagBite      MoveFlag = "bite"      // 咬类
	FlagPulse     MoveFlag = "pulse"     // 波动类
	FlagBall      MoveFlag = "ball"      // 球类
	FlagPowder    MoveFlag = "powder"    // 粉末类
	FlagRecoil    MoveFlag = "recoil"    // 反作用力
	FlagHeal      MoveFlag = "heal"      // 回复类
	FlagProtect   MoveFlag = "protect"   // 守护类
	FlagPriority  MoveFlag = "priority"  // 先制类
	FlagMultiHit  MoveFlag = "multihit"  // 连续攻击
	FlagCharge    MoveFlag = "charge"    // 蓄力类
	FlagRecharge  MoveFlag = "recharge"  // 充能类
	FlagDance     MoveFlag = "dance"     // 舞蹈类
	FlagSlicing   MoveFlag = "slicing"   // 切割类
	FlagWind      MoveFlag = "wind"      // 风类
)
