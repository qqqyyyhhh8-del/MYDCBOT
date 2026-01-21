package entity

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// Pokemon 宝可梦实体（图鉴数据）
type Pokemon struct {
	ID              int                      // 全国图鉴编号
	Name            string                   // 名称
	Types           []valueobject.PokeType   // 属性（最多2个）
	BaseHP          int                      // 基础HP
	BaseAtk         int                      // 基础攻击
	BaseDef         int                      // 基础防御
	BaseSpAtk       int                      // 基础特攻
	BaseSpDef       int                      // 基础特防
	BaseSpeed       int                      // 基础速度
	Abilities       []valueobject.Ability    // 可用特性
	HiddenAbility   *valueobject.Ability     // 隐藏特性
	SelectedAbility *valueobject.Ability     // 选择的特性
	Nature          valueobject.Nature       // 性格
	TeraType        valueobject.PokeType     // 太晶属性
	LearnableMoves  []*Move                  // 可学习技能
	SpriteURL       string                   // 精灵图URL
	CanMegaEvolve   bool                     // 是否可超级进化
	MegaStoneID     int                      // 超级石ID
	CanGigantamax   bool                     // 是否可极巨化
}

// PokemonBuild 宝可梦配置（玩家自定义）
type PokemonBuild struct {
	Pokemon       *Pokemon                 // 基础宝可梦
	Nickname      string                   // 昵称
	Level         int                      // 等级
	Ability       *valueobject.Ability     // 选择的特性
	Nature        valueobject.Nature       // 性格
	Item          *valueobject.Item        // 携带道具
	Moves         []*Move                  // 选择的技能（最多4个）
	TeraType      valueobject.PokeType     // 太晶属性
	IVs           Stats                    // 个体值 (0-31)
	EVs           Stats                    // 努力值 (0-252, 总计510)
	Shiny         bool                     // 是否闪光
	Gender        Gender                   // 性别
}

// Stats 六维属性
type Stats struct {
	HP    int
	Atk   int
	Def   int
	SpAtk int
	SpDef int
	Speed int
}

// Gender 性别
type Gender string

const (
	GenderMale    Gender = "♂"
	GenderFemale  Gender = "♀"
	GenderUnknown Gender = "-"
)

// NewPokemonBuild 创建宝可梦配置
func NewPokemonBuild(pokemon *Pokemon) *PokemonBuild {
	build := &PokemonBuild{
		Pokemon:  pokemon,
		Level:    50,
		Nature:   valueobject.NatureHardy,
		TeraType: pokemon.Types[0],
		IVs:      Stats{31, 31, 31, 31, 31, 31}, // 默认6V
		EVs:      Stats{0, 0, 0, 0, 0, 0},
		Moves:    make([]*Move, 0, 4),
		Gender:   GenderUnknown,
	}
	// 默认选择第一个特性
	if len(pokemon.Abilities) > 0 {
		build.Ability = &pokemon.Abilities[0]
	}
	return build
}

// SetIVs 设置个体值
func (b *PokemonBuild) SetIVs(hp, atk, def, spAtk, spDef, speed int) {
	b.IVs = Stats{
		HP:    clamp(hp, 0, 31),
		Atk:   clamp(atk, 0, 31),
		Def:   clamp(def, 0, 31),
		SpAtk: clamp(spAtk, 0, 31),
		SpDef: clamp(spDef, 0, 31),
		Speed: clamp(speed, 0, 31),
	}
}

// SetEVs 设置努力值
func (b *PokemonBuild) SetEVs(hp, atk, def, spAtk, spDef, speed int) {
	// 单项最大252，总计最大510
	total := hp + atk + def + spAtk + spDef + speed
	if total > 510 {
		return // 超过上限不设置
	}
	b.EVs = Stats{
		HP:    clamp(hp, 0, 252),
		Atk:   clamp(atk, 0, 252),
		Def:   clamp(def, 0, 252),
		SpAtk: clamp(spAtk, 0, 252),
		SpDef: clamp(spDef, 0, 252),
		Speed: clamp(speed, 0, 252),
	}
}

// AddMove 添加技能
func (b *PokemonBuild) AddMove(move *Move) bool {
	if len(b.Moves) >= 4 {
		return false
	}
	b.Moves = append(b.Moves, move)
	return true
}

// SetMoves 设置技能
func (b *PokemonBuild) SetMoves(moves []*Move) {
	if len(moves) > 4 {
		moves = moves[:4]
	}
	b.Moves = moves
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// NewPokemon 创建宝可梦
func NewPokemon(id int, name string, types []valueobject.PokeType) *Pokemon {
	return &Pokemon{
		ID:             id,
		Name:           name,
		Types:          types,
		LearnableMoves: make([]*Move, 0),
	}
}

// SetBaseStats 设置基础属性
func (p *Pokemon) SetBaseStats(hp, atk, def, spAtk, spDef, speed int) {
	p.BaseHP = hp
	p.BaseAtk = atk
	p.BaseDef = def
	p.BaseSpAtk = spAtk
	p.BaseSpDef = spDef
	p.BaseSpeed = speed
}

// AddLearnableMove 添加可学习技能
func (p *Pokemon) AddLearnableMove(move *Move) {
	p.LearnableMoves = append(p.LearnableMoves, move)
}

// GetSpriteURL 获取精灵图URL
func (p *Pokemon) GetSpriteURL() string {
	if p.SpriteURL != "" {
		return p.SpriteURL
	}
	return p.GenerateSpriteURL()
}

// GenerateSpriteURL 生成精灵图URL
func (p *Pokemon) GenerateSpriteURL() string {
	return "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/" + 
		   itoa(p.ID) + ".gif"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// Move 技能
type Move struct {
	Name     string                 // 技能名称
	Type     valueobject.PokeType   // 技能属性
	Category MoveCategory           // 技能分类
	Power    int                    // 威力
	Accuracy int                    // 命中率
	PP       int                    // PP值
	MaxPP    int                    // 最大PP
	Priority int                    // 优先度（-7 到 +5）
	
	// 技能效果
	RechargeRequired bool           // 使用后需要充能（如破坏光线）
	ChargeRequired   bool           // 使用前需要蓄力（如日光束）
	MakesContact     bool           // 是否为接触技能
	EffectChance     int            // 追加效果触发概率（0-100）
}

// MoveCategory 技能分类
type MoveCategory string

const (
	CategoryPhysical MoveCategory = "物理"
	CategorySpecial  MoveCategory = "特殊"
	CategoryStatus   MoveCategory = "变化"
)

// NewMove 创建技能
func NewMove(name string, pokeType valueobject.PokeType, category MoveCategory, power, accuracy, pp int) *Move {
	return &Move{
		Name:     name,
		Type:     pokeType,
		Category: category,
		Power:    power,
		Accuracy: accuracy,
		PP:       pp,
		MaxPP:    pp,
	}
}

// CanUse 是否可以使用
func (m *Move) CanUse() bool {
	return m.PP > 0
}

// Use 使用技能
func (m *Move) Use() bool {
	if m.PP <= 0 {
		return false
	}
	m.PP--
	return true
}

// Restore 恢复PP
func (m *Move) Restore() {
	m.PP = m.MaxPP
}
