package valueobject

// Ability 特性
type Ability struct {
	ID          int
	Name        string
	Description string
	IsHidden    bool // 是否为隐藏特性
}

// 常用特性定义
var (
	AbilityOvergrow     = Ability{65, "茂盛", "HP低时草属性招式威力提升", false}
	AbilityBlaze        = Ability{66, "猛火", "HP低时火属性招式威力提升", false}
	AbilityTorrent      = Ability{67, "激流", "HP低时水属性招式威力提升", false}
	AbilityStatic       = Ability{9, "静电", "接触时可能使对手麻痹", false}
	AbilityLightningRod = Ability{31, "避雷针", "吸引电属性招式并提升特攻", true}
	AbilityLevitate     = Ability{26, "飘浮", "免疫地面属性招式", false}
	AbilityIntimdate    = Ability{22, "威吓", "出场时降低对手攻击", false}
	AbilityMoxie        = Ability{153, "自信过剩", "击倒对手后攻击提升", true}
	AbilityMultiscale   = Ability{136, "多重鳞片", "HP满时受到伤害减半", true}
	AbilityInnerFocus   = Ability{39, "精神力", "不会畏缩", false}
	AbilityPressure     = Ability{46, "压迫感", "对手消耗更多PP", false}
	AbilityUnnerve      = Ability{127, "紧张感", "对手无法食用树果", true}
	AbilityCursedBody   = Ability{130, "诅咒之躯", "被攻击时可能封印对手招式", false}
	AbilityImmunity     = Ability{17, "免疫", "不会中毒", false}
	AbilityThickFat     = Ability{47, "厚脂肪", "火和冰属性伤害减半", true}
	AbilityGluttony     = Ability{82, "贪吃鬼", "提前食用树果", true}
	AbilitySwiftSwim    = Ability{33, "悠游自如", "雨天速度翻倍", false}
	AbilityChlorophyll  = Ability{34, "叶绿素", "晴天速度翻倍", false}
	AbilitySandRush     = Ability{146, "拨沙", "沙暴时速度翻倍", true}
	AbilitySlushRush    = Ability{202, "拨雪", "冰雹时速度翻倍", true}
	AbilityDrizzle      = Ability{2, "降雨", "出场时召唤雨天", false}
	AbilityDrought      = Ability{70, "日照", "出场时召唤晴天", false}
	AbilitySandStream   = Ability{45, "扬沙", "出场时召唤沙暴", false}
	AbilitySnowWarning  = Ability{117, "降雪", "出场时召唤冰雹", false}
	AbilityProtean      = Ability{168, "变幻自如", "使用招式前变为该属性", true}
	AbilityLibero       = Ability{236, "自由者", "使用招式前变为该属性", true}
	AbilityHugePower    = Ability{37, "大力士", "攻击翻倍", false}
	AbilityPurePower    = Ability{74, "瑜伽之力", "攻击翻倍", false}
	AbilitySpeedBoost   = Ability{3, "加速", "每回合速度提升", false}
	AbilityWonderGuard  = Ability{25, "神奇守护", "只会被效果拔群的招式击中", false}
	AbilityMagicGuard   = Ability{98, "魔法防守", "只受到攻击招式的伤害", false}
	AbilityMagicBounce  = Ability{156, "魔法镜", "反弹变化招式", true}
	AbilityPrankster    = Ability{158, "恶作剧之心", "变化招式优先度+1", false}
	AbilityGaleWings    = Ability{177, "疾风之翼", "HP满时飞行招式优先度+1", true}
	AbilityToughClaws   = Ability{181, "硬爪", "接触招式威力提升30%", false}
	AbilityStrongJaw    = Ability{173, "强壮之颚", "咬类招式威力提升50%", false}
	AbilitySheerForce   = Ability{125, "强行", "放弃追加效果提升威力", true}
	AbilityTechnician   = Ability{101, "技术高手", "威力60以下招式威力提升50%", false}
	AbilityAdaptability = Ability{91, "适应力", "本属性加成变为2倍", false}
)

// AbilityMap 特性ID映射
var AbilityMap = map[int]Ability{
	65:  AbilityOvergrow,
	66:  AbilityBlaze,
	67:  AbilityTorrent,
	9:   AbilityStatic,
	31:  AbilityLightningRod,
	26:  AbilityLevitate,
	22:  AbilityIntimdate,
	153: AbilityMoxie,
	136: AbilityMultiscale,
	39:  AbilityInnerFocus,
	46:  AbilityPressure,
	127: AbilityUnnerve,
	130: AbilityCursedBody,
	17:  AbilityImmunity,
	47:  AbilityThickFat,
	82:  AbilityGluttony,
	33:  AbilitySwiftSwim,
	34:  AbilityChlorophyll,
	146: AbilitySandRush,
	202: AbilitySlushRush,
	2:   AbilityDrizzle,
	70:  AbilityDrought,
	45:  AbilitySandStream,
	117: AbilitySnowWarning,
	168: AbilityProtean,
	236: AbilityLibero,
	37:  AbilityHugePower,
	74:  AbilityPurePower,
	3:   AbilitySpeedBoost,
	25:  AbilityWonderGuard,
	98:  AbilityMagicGuard,
	156: AbilityMagicBounce,
	158: AbilityPrankster,
	177: AbilityGaleWings,
	181: AbilityToughClaws,
	173: AbilityStrongJaw,
	125: AbilitySheerForce,
	101: AbilityTechnician,
	91:  AbilityAdaptability,
}

// GetAbilityByID 通过ID获取特性
func GetAbilityByID(id int) *Ability {
	if ability, ok := AbilityMap[id]; ok {
		return &ability
	}
	return nil
}

// GetAbilityByName 通过名称获取特性
func GetAbilityByName(name string) *Ability {
	for _, ability := range AbilityMap {
		if ability.Name == name {
			return &ability
		}
	}
	return nil
}
