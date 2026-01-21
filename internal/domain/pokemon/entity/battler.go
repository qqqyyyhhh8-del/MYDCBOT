package entity

import (
	"math/rand"
	"time"

	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// Battler 对战中的宝可梦
type Battler struct {
	Build     *PokemonBuild              // 宝可梦配置
	Pokemon   *Pokemon                   // 宝可梦数据（快捷引用）
	Level     int                        // 等级
	CurrentHP int                        // 当前HP
	MaxHP     int                        // 最大HP
	Atk       int                        // 攻击
	Def       int                        // 防御
	SpAtk     int                        // 特攻
	SpDef     int                        // 特防
	Speed     int                        // 速度
	Moves     []*Move                    // 技能副本

	// 对战状态
	StatStages    StatStages             // 能力等级变化
	Status        StatusCondition        // 异常状态
	StatusTurns   int                    // 状态持续回合
	Volatile      []VolatileStatus       // 临时状态
	Types         []valueobject.PokeType // 当前属性
	Ability       *valueobject.Ability   // 当前特性
	Item          *valueobject.Item      // 当前道具
	ItemConsumed  bool                   // 道具是否已消耗
	LastMove      *Move                  // 上一次使用的技能
	LastMoveTurns int                    // 连续使用同技能的回合数
	Protected     bool                   // 是否处于守住状态
	Flinched      bool                   // 是否畏缩
	MustRecharge  bool                   // 下回合必须充能（如破坏光线后）
	IsCharging    bool                   // 正在蓄力中（如日光束）
	ChargingMove  *Move                  // 蓄力中的技能

	// 特殊系统状态
	IsMega         bool                 // 是否已超级进化
	IsDynamaxed    bool                 // 是否已极巨化
	DynamaxTurns   int                  // 极巨化剩余回合
	IsTerastalized bool                 // 是否已太晶化
	TeraType       valueobject.PokeType // 太晶属性
	UsedZMove      bool                 // 是否已使用Z招式
}

// StatStages 能力等级
type StatStages struct {
	Atk      int // -6 到 +6
	Def      int
	SpAtk    int
	SpDef    int
	Speed    int
	Accuracy int
	Evasion  int
}

// StatusCondition 异常状态
type StatusCondition string

const (
	StatusNone      StatusCondition = ""
	StatusPoison    StatusCondition = "中毒"
	StatusBadPoison StatusCondition = "剧毒"
	StatusBurn      StatusCondition = "灼伤"
	StatusParalyze  StatusCondition = "麻痹"
	StatusSleep     StatusCondition = "睡眠"
	StatusFreeze    StatusCondition = "冰冻"
)

// VolatileStatus 临时状态
type VolatileStatus string

const (
	VolatileConfusion   VolatileStatus = "混乱"
	VolatileAttraction  VolatileStatus = "着迷"
	VolatileTaunt       VolatileStatus = "挑衅"
	VolatileTorment     VolatileStatus = "无理取闹"
	VolatileDisable     VolatileStatus = "定身法"
	VolatileEncore      VolatileStatus = "再来一次"
	VolatileLeechSeed   VolatileStatus = "寄生种子"
	VolatileSubstitute  VolatileStatus = "替身"
	VolatileFocusEnergy VolatileStatus = "聚气"
)

// NewBattler 创建对战宝可梦（简化版，兼容旧代码）
func NewBattler(pokemon *Pokemon, level int) *Battler {
	build := NewPokemonBuild(pokemon)
	build.Level = level
	// 使用玩家选择的特性（如果已设置）
	if pokemon.SelectedAbility != nil {
		build.Ability = pokemon.SelectedAbility
	}
	if len(pokemon.LearnableMoves) > 0 {
		for i := 0; i < 4 && i < len(pokemon.LearnableMoves); i++ {
			build.AddMove(pokemon.LearnableMoves[i])
		}
	}
	return NewBattlerFromBuild(build)
}

// NewBattlerFromBuild 从配置创建对战宝可梦
func NewBattlerFromBuild(build *PokemonBuild) *Battler {
	b := &Battler{
		Build:      build,
		Pokemon:    build.Pokemon,
		Level:      build.Level,
		Types:      append([]valueobject.PokeType{}, build.Pokemon.Types...),
		Ability:    build.Ability,
		Item:       build.Item,
		TeraType:   build.TeraType,
		StatStages: StatStages{},
		Volatile:   make([]VolatileStatus, 0),
	}
	b.calculateStats()
	b.copyMoves()
	return b
}

// calculateStats 计算实际属性（完整公式）
func (b *Battler) calculateStats() {
	pokemon := b.Pokemon
	level := b.Level
	ivs := b.Build.IVs
	evs := b.Build.EVs
	nature := valueobject.GetNatureModifier(b.Build.Nature)

	// HP = floor((2 * Base + IV + floor(EV/4)) * Level / 100) + Level + 10
	b.MaxHP = ((2*pokemon.BaseHP+ivs.HP+evs.HP/4)*level)/100 + level + 10
	b.CurrentHP = b.MaxHP

	// 其他属性 = floor((floor((2 * Base + IV + floor(EV/4)) * Level / 100) + 5) * Nature)
	b.Atk = int(float64(((2*pokemon.BaseAtk+ivs.Atk+evs.Atk/4)*level)/100+5) * nature.Atk)
	b.Def = int(float64(((2*pokemon.BaseDef+ivs.Def+evs.Def/4)*level)/100+5) * nature.Def)
	b.SpAtk = int(float64(((2*pokemon.BaseSpAtk+ivs.SpAtk+evs.SpAtk/4)*level)/100+5) * nature.SpAtk)
	b.SpDef = int(float64(((2*pokemon.BaseSpDef+ivs.SpDef+evs.SpDef/4)*level)/100+5) * nature.SpDef)
	b.Speed = int(float64(((2*pokemon.BaseSpeed+ivs.Speed+evs.Speed/4)*level)/100+5) * nature.Speed)
}

// copyMoves 复制技能
func (b *Battler) copyMoves() {
	sourceMoves := b.Build.Moves
	if len(sourceMoves) == 0 {
		sourceMoves = b.Pokemon.LearnableMoves
	}
	b.Moves = make([]*Move, len(sourceMoves))
	for i, m := range sourceMoves {
		b.Moves[i] = &Move{
			Name:     m.Name,
			Type:     m.Type,
			Category: m.Category,
			Power:    m.Power,
			Accuracy: m.Accuracy,
			PP:       m.PP,
			MaxPP:    m.MaxPP,
		}
	}
}

// GetEffectiveAtk 获取有效攻击力（含能力等级）
func (b *Battler) GetEffectiveAtk() int {
	return applyStatStage(b.Atk, b.StatStages.Atk)
}

// GetEffectiveDef 获取有效防御力
func (b *Battler) GetEffectiveDef() int {
	return applyStatStage(b.Def, b.StatStages.Def)
}

// GetEffectiveSpAtk 获取有效特攻
func (b *Battler) GetEffectiveSpAtk() int {
	return applyStatStage(b.SpAtk, b.StatStages.SpAtk)
}

// GetEffectiveSpDef 获取有效特防
func (b *Battler) GetEffectiveSpDef() int {
	return applyStatStage(b.SpDef, b.StatStages.SpDef)
}

// GetEffectiveSpeed 获取有效速度
func (b *Battler) GetEffectiveSpeed() int {
	speed := applyStatStage(b.Speed, b.StatStages.Speed)
	if b.Status == StatusParalyze {
		speed = speed / 2
	}
	if b.Item != nil && b.Item.Name == "讲究围巾" && !b.ItemConsumed {
		speed = speed * 3 / 2
	}
	return speed
}

// applyStatStage 应用能力等级
func applyStatStage(base int, stage int) int {
	if stage > 6 {
		stage = 6
	}
	if stage < -6 {
		stage = -6
	}
	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

// ModifyStat 修改能力等级
func (b *Battler) ModifyStat(stat string, stages int) (int, bool) {
	var target *int
	switch stat {
	case "atk":
		target = &b.StatStages.Atk
	case "def":
		target = &b.StatStages.Def
	case "spatk":
		target = &b.StatStages.SpAtk
	case "spdef":
		target = &b.StatStages.SpDef
	case "speed":
		target = &b.StatStages.Speed
	case "accuracy":
		target = &b.StatStages.Accuracy
	case "evasion":
		target = &b.StatStages.Evasion
	default:
		return 0, false
	}

	oldValue := *target
	*target += stages
	if *target > 6 {
		*target = 6
	}
	if *target < -6 {
		*target = -6
	}
	return *target - oldValue, *target != oldValue
}

// IsAlive 是否存活
func (b *Battler) IsAlive() bool {
	return b.CurrentHP > 0
}

// TakeDamage 受到伤害
func (b *Battler) TakeDamage(damage int) int {
	if damage < 0 {
		damage = 0
	}
	if damage > b.CurrentHP {
		damage = b.CurrentHP
	}
	b.CurrentHP -= damage
	return damage
}

// Heal 恢复HP
func (b *Battler) Heal(amount int) int {
	if amount < 0 {
		amount = 0
	}
	healed := amount
	if b.CurrentHP+amount > b.MaxHP {
		healed = b.MaxHP - b.CurrentHP
	}
	b.CurrentHP += healed
	return healed
}

// GetHPPercent 获取HP百分比
func (b *Battler) GetHPPercent() float64 {
	return float64(b.CurrentHP) / float64(b.MaxHP) * 100
}

// DamageResult 伤害计算结果
type DamageResult struct {
	Damage        int
	Effectiveness float64
	Hit           bool
	Critical      bool
}

// CalculateDamage 计算伤害（完整公式）
func (b *Battler) CalculateDamage(move *Move, target *Battler) DamageResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := DamageResult{Hit: false}

	// 命中判定
	accuracy := move.Accuracy
	if accuracy > 0 {
		accMod := applyStatStage(100, b.StatStages.Accuracy-target.StatStages.Evasion)
		effectiveAcc := accuracy * accMod / 100
		if r.Intn(100) >= effectiveAcc {
			return result
		}
	}
	result.Hit = true

	// 变化技能不造成伤害
	if move.Category == CategoryStatus {
		result.Effectiveness = 1
		return result
	}

	// 选择攻击和防御属性
	var atk, def int
	if move.Category == CategoryPhysical {
		atk = b.GetEffectiveAtk()
		def = target.GetEffectiveDef()
		// 灼伤减攻击（除非有毅力特性）
		if b.Status == StatusBurn && (b.Ability == nil || b.Ability.Name != "毅力") {
			atk = atk / 2
		}
	} else {
		atk = b.GetEffectiveSpAtk()
		def = target.GetEffectiveSpDef()
	}

	// 道具加成
	if b.Item != nil && !b.ItemConsumed {
		switch b.Item.Name {
		case "讲究头带":
			if move.Category == CategoryPhysical {
				atk = atk * 3 / 2
			}
		case "讲究眼镜":
			if move.Category == CategorySpecial {
				atk = atk * 3 / 2
			}
		}
	}

	// 极巨化威力转换
	power := move.Power
	if b.IsDynamaxed {
		if power < 40 {
			power = 90
		} else if power < 50 {
			power = 100
		} else if power < 60 {
			power = 110
		} else if power < 70 {
			power = 120
		} else if power < 100 {
			power = 130
		} else if power < 140 {
			power = 140
		} else {
			power = 150
		}
	}

	// 基础伤害公式
	baseDamage := ((2*b.Level/5+2)*power*atk/def)/50 + 2

	// 属性克制
	defenseTypes := target.Types
	if target.IsTerastalized {
		defenseTypes = []valueobject.PokeType{target.TeraType}
	}
	result.Effectiveness = valueobject.GetEffectiveness(move.Type, defenseTypes)

	// 同属性加成 (STAB)
	stab := 1.0
	attackTypes := b.Types
	if b.IsTerastalized {
		attackTypes = append(attackTypes, b.TeraType)
	}
	for _, t := range attackTypes {
		if t == move.Type {
			stab = 1.5
			// 适应力特性
			if b.Ability != nil && b.Ability.Name == "适应力" {
				stab = 2.0
			}
			break
		}
	}

	// 随机因子 (85-100%)
	randomFactor := float64(r.Intn(16)+85) / 100.0

	// 会心一击判定
	critical := 1.0
	critStage := 0
	for _, v := range b.Volatile {
		if v == VolatileFocusEnergy {
			critStage += 2
		}
	}
	critChances := []int{24, 8, 2, 1}
	if critStage > 3 {
		critStage = 3
	}
	if r.Intn(critChances[critStage]) == 0 {
		critical = 1.5
		result.Critical = true
		// 狙击手特性
		if b.Ability != nil && b.Ability.Name == "狙击手" {
			critical = 2.25
		}
	}

	// 道具威力加成
	itemMod := 1.0
	if b.Item != nil && !b.ItemConsumed {
		switch b.Item.Name {
		case "生命宝珠":
			itemMod = 1.3
		case "达人带":
			if result.Effectiveness > 1 {
				itemMod = 1.2
			}
		}
	}

	// 最终伤害
	result.Damage = int(float64(baseDamage) * result.Effectiveness * stab * randomFactor * critical * itemMod)
	if result.Damage < 1 && result.Effectiveness > 0 {
		result.Damage = 1
	}

	return result
}

// TakeDamageWithItem 受到伤害（含道具效果）
func (b *Battler) TakeDamageWithItem(damage int) int {
	// 气势披带
	if b.Item != nil && b.Item.Name == "气势披带" && !b.ItemConsumed && b.CurrentHP == b.MaxHP && damage >= b.CurrentHP {
		damage = b.CurrentHP - 1
		b.ItemConsumed = true
	}
	return b.TakeDamage(damage)
}
