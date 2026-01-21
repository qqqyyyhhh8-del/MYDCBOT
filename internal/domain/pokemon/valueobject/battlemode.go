package valueobject

// BattleMode 对战模式
type BattleMode string

const (
	// 基础模式
	ModeSingle     BattleMode = "single"      // 单打
	ModeDouble     BattleMode = "double"      // 双打（2人各控制2只）
	ModeMulti      BattleMode = "multi"       // 多人双打（4人各控制1只）
	ModeRotation   BattleMode = "rotation"    // 轮盘对战
	ModeTriple     BattleMode = "triple"      // 三打（已废弃但可支持）
	ModeBattleRoyal BattleMode = "royal"      // 皇家对战（4人混战）
)

// BattleRule 对战规则
type BattleRule string

const (
	RuleNoRule       BattleRule = "no_rule"       // 无规则
	RuleFlat         BattleRule = "flat"          // 平坦规则（等级压制到50）
	RuleNoRestrict   BattleRule = "no_restrict"   // 无限制
	RuleLittleCup    BattleRule = "little_cup"    // 小精灵杯（Lv5未进化）
	RuleOU           BattleRule = "ou"            // OU规则
	RuleUbers        BattleRule = "ubers"         // Ubers规则
	RuleVGC          BattleRule = "vgc"           // VGC官方规则
	RuleBSS          BattleRule = "bss"           // Battle Stadium Singles
	RuleBSD          BattleRule = "bsd"           // Battle Stadium Doubles
)

// FunRule 娱乐规则（可叠加）
type FunRule string

const (
	FunRuleNone           FunRule = "none"            // 无
	FunRuleTypeReverse    FunRule = "type_reverse"    // 属性反转
	FunRuleRandomMoves    FunRule = "random_moves"    // 随机技能
	FunRuleRandomAbility  FunRule = "random_ability"  // 随机特性
	FunRuleWonderLauncher FunRule = "wonder_launcher" // 奇迹发射器
	FunRuleInverseBattle  FunRule = "inverse"         // 逆属性对战
	FunRuleRandomPokemon  FunRule = "random_pokemon"  // 随机宝可梦
	FunRuleSameType       FunRule = "same_type"       // 同属性队伍
	FunRuleMonotype       FunRule = "monotype"        // 单属性限制
)

// GimmickSystem 特殊系统
type GimmickSystem string

const (
	GimmickNone        GimmickSystem = "none"          // 无
	GimmickMegaEvo     GimmickSystem = "mega"          // 超级进化
	GimmickZMove       GimmickSystem = "z_move"        // Z招式
	GimmickDynamax     GimmickSystem = "dynamax"       // 极巨化
	GimmickTerastal    GimmickSystem = "terastal"      // 太晶化
	GimmickAll         GimmickSystem = "all"           // 全部可用
)

// BattleConfig 对战配置
type BattleConfig struct {
	Mode           BattleMode      // 对战模式
	Rule           BattleRule      // 对战规则
	FunRules       []FunRule       // 娱乐规则（可多选）
	Gimmick        GimmickSystem   // 特殊系统
	LevelCap       int             // 等级上限
	TeamSize       int             // 队伍大小
	BringCount     int             // 实际出战数量
	TimerEnabled   bool            // 是否启用计时器
	TurnTimeLimit  int             // 每回合时间限制（秒）
	TotalTimeLimit int             // 总时间限制（秒）
	ItemClause     bool            // 道具条款（同队不能重复道具）
	SpeciesClause  bool            // 种族条款（同队不能重复宝可梦）
	SleepClause    bool            // 睡眠条款
	OHKOClause     bool            // 一击必杀条款
	EvasionClause  bool            // 闪避条款
}

// DefaultSingleConfig 默认单打配置
func DefaultSingleConfig() *BattleConfig {
	return &BattleConfig{
		Mode:          ModeSingle,
		Rule:          RuleFlat,
		FunRules:      []FunRule{},
		Gimmick:       GimmickNone,
		LevelCap:      50,
		TeamSize:      6,
		BringCount:    3,
		TimerEnabled:  false,
		ItemClause:    true,
		SpeciesClause: true,
		SleepClause:   true,
		OHKOClause:    true,
		EvasionClause: true,
	}
}

// DefaultDoubleConfig 默认双打配置
func DefaultDoubleConfig() *BattleConfig {
	return &BattleConfig{
		Mode:          ModeDouble,
		Rule:          RuleFlat,
		FunRules:      []FunRule{},
		Gimmick:       GimmickNone,
		LevelCap:      50,
		TeamSize:      6,
		BringCount:    4,
		TimerEnabled:  false,
		ItemClause:    true,
		SpeciesClause: true,
		SleepClause:   true,
		OHKOClause:    true,
		EvasionClause: true,
	}
}

// QuickBattleConfig 快速对战配置（无规则）
func QuickBattleConfig() *BattleConfig {
	return &BattleConfig{
		Mode:          ModeSingle,
		Rule:          RuleNoRule,
		FunRules:      []FunRule{},
		Gimmick:       GimmickAll,
		LevelCap:      100,
		TeamSize:      1,
		BringCount:    1,
		TimerEnabled:  false,
		ItemClause:    false,
		SpeciesClause: false,
		SleepClause:   false,
		OHKOClause:    false,
		EvasionClause: false,
	}
}

// VGCConfig VGC规则配置
func VGCConfig() *BattleConfig {
	return &BattleConfig{
		Mode:           ModeDouble,
		Rule:           RuleVGC,
		FunRules:       []FunRule{},
		Gimmick:        GimmickTerastal,
		LevelCap:       50,
		TeamSize:       6,
		BringCount:     4,
		TimerEnabled:   true,
		TurnTimeLimit:  45,
		TotalTimeLimit: 900,
		ItemClause:     true,
		SpeciesClause:  true,
		SleepClause:    false,
		OHKOClause:     false,
		EvasionClause:  false,
	}
}

// GetModeDisplayName 获取模式显示名称
func (m BattleMode) DisplayName() string {
	names := map[BattleMode]string{
		ModeSingle:      "单打",
		ModeDouble:      "双打",
		ModeMulti:       "多人双打",
		ModeRotation:    "轮盘对战",
		ModeTriple:      "三打",
		ModeBattleRoyal: "皇家对战",
	}
	if name, ok := names[m]; ok {
		return name
	}
	return string(m)
}

// GetRuleDisplayName 获取规则显示名称
func (r BattleRule) DisplayName() string {
	names := map[BattleRule]string{
		RuleNoRule:     "无规则",
		RuleFlat:       "平坦规则",
		RuleNoRestrict: "无限制",
		RuleLittleCup:  "小精灵杯",
		RuleOU:         "OU",
		RuleUbers:      "Ubers",
		RuleVGC:        "VGC",
		RuleBSS:        "单打竞技场",
		RuleBSD:        "双打竞技场",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return string(r)
}

// GetGimmickDisplayName 获取系统显示名称
func (g GimmickSystem) DisplayName() string {
	names := map[GimmickSystem]string{
		GimmickNone:     "无",
		GimmickMegaEvo:  "超级进化",
		GimmickZMove:    "Z招式",
		GimmickDynamax:  "极巨化",
		GimmickTerastal: "太晶化",
		GimmickAll:      "全部可用",
	}
	if name, ok := names[g]; ok {
		return name
	}
	return string(g)
}
