package valueobject

// PokeType 宝可梦属性类型
type PokeType string

const (
	TypeNormal   PokeType = "一般"
	TypeFire     PokeType = "火"
	TypeWater    PokeType = "水"
	TypeElectric PokeType = "电"
	TypeGrass    PokeType = "草"
	TypeIce      PokeType = "冰"
	TypeFighting PokeType = "格斗"
	TypePoison   PokeType = "毒"
	TypeGround   PokeType = "地面"
	TypeFlying   PokeType = "飞行"
	TypePsychic  PokeType = "超能力"
	TypeBug      PokeType = "虫"
	TypeRock     PokeType = "岩石"
	TypeGhost    PokeType = "幽灵"
	TypeDragon   PokeType = "龙"
	TypeDark     PokeType = "恶"
	TypeSteel    PokeType = "钢"
	TypeFairy    PokeType = "妖精"
)

// TypeEffectiveness 属性克制表
var TypeEffectiveness = map[PokeType]map[PokeType]float64{
	TypeNormal: {
		TypeRock:  0.5,
		TypeGhost: 0,
		TypeSteel: 0.5,
	},
	TypeFire: {
		TypeFire:  0.5,
		TypeWater: 0.5,
		TypeGrass: 2,
		TypeIce:   2,
		TypeBug:   2,
		TypeRock:  0.5,
		TypeDragon: 0.5,
		TypeSteel: 2,
	},
	TypeWater: {
		TypeFire:   2,
		TypeWater:  0.5,
		TypeGrass:  0.5,
		TypeGround: 2,
		TypeRock:   2,
		TypeDragon: 0.5,
	},
	TypeElectric: {
		TypeWater:    2,
		TypeElectric: 0.5,
		TypeGrass:    0.5,
		TypeGround:   0,
		TypeFlying:   2,
		TypeDragon:   0.5,
	},
	TypeGrass: {
		TypeFire:   0.5,
		TypeWater:  2,
		TypeGrass:  0.5,
		TypePoison: 0.5,
		TypeGround: 2,
		TypeFlying: 0.5,
		TypeBug:    0.5,
		TypeRock:   2,
		TypeDragon: 0.5,
		TypeSteel:  0.5,
	},
	TypeIce: {
		TypeFire:  0.5,
		TypeWater: 0.5,
		TypeGrass: 2,
		TypeIce:   0.5,
		TypeGround: 2,
		TypeFlying: 2,
		TypeDragon: 2,
		TypeSteel:  0.5,
	},
	TypeFighting: {
		TypeNormal:  2,
		TypeIce:     2,
		TypePoison:  0.5,
		TypeFlying:  0.5,
		TypePsychic: 0.5,
		TypeBug:     0.5,
		TypeRock:    2,
		TypeGhost:   0,
		TypeDark:    2,
		TypeSteel:   2,
		TypeFairy:   0.5,
	},
	TypePoison: {
		TypeGrass:  2,
		TypePoison: 0.5,
		TypeGround: 0.5,
		TypeRock:   0.5,
		TypeGhost:  0.5,
		TypeSteel:  0,
		TypeFairy:  2,
	},
	TypeGround: {
		TypeFire:     2,
		TypeElectric: 2,
		TypeGrass:    0.5,
		TypePoison:   2,
		TypeFlying:   0,
		TypeBug:      0.5,
		TypeRock:     2,
		TypeSteel:    2,
	},
	TypeFlying: {
		TypeElectric: 0.5,
		TypeGrass:    2,
		TypeFighting: 2,
		TypeBug:      2,
		TypeRock:     0.5,
		TypeSteel:    0.5,
	},
	TypePsychic: {
		TypeFighting: 2,
		TypePoison:   2,
		TypePsychic:  0.5,
		TypeDark:     0,
		TypeSteel:    0.5,
	},
	TypeBug: {
		TypeFire:     0.5,
		TypeGrass:    2,
		TypeFighting: 0.5,
		TypePoison:   0.5,
		TypeFlying:   0.5,
		TypePsychic:  2,
		TypeGhost:    0.5,
		TypeDark:     2,
		TypeSteel:    0.5,
		TypeFairy:    0.5,
	},
	TypeRock: {
		TypeFire:     2,
		TypeIce:      2,
		TypeFighting: 0.5,
		TypeGround:   0.5,
		TypeFlying:   2,
		TypeBug:      2,
		TypeSteel:    0.5,
	},
	TypeGhost: {
		TypeNormal:  0,
		TypePsychic: 2,
		TypeGhost:   2,
		TypeDark:    0.5,
	},
	TypeDragon: {
		TypeDragon: 2,
		TypeSteel:  0.5,
		TypeFairy:  0,
	},
	TypeDark: {
		TypeFighting: 0.5,
		TypePsychic:  2,
		TypeGhost:    2,
		TypeDark:     0.5,
		TypeFairy:    0.5,
	},
	TypeSteel: {
		TypeFire:     0.5,
		TypeWater:    0.5,
		TypeElectric: 0.5,
		TypeIce:      2,
		TypeRock:     2,
		TypeSteel:    0.5,
		TypeFairy:    2,
	},
	TypeFairy: {
		TypeFire:     0.5,
		TypeFighting: 2,
		TypePoison:   0.5,
		TypeDragon:   2,
		TypeDark:     2,
		TypeSteel:    0.5,
	},
}

// GetEffectiveness 获取属性克制倍率
func GetEffectiveness(attackType PokeType, defenseTypes []PokeType) float64 {
	multiplier := 1.0
	attackEffects, ok := TypeEffectiveness[attackType]
	if !ok {
		return multiplier
	}
	for _, defType := range defenseTypes {
		if effect, exists := attackEffects[defType]; exists {
			multiplier *= effect
		}
	}
	return multiplier
}

// PokeTypeFromString 从字符串转换为属性类型
func PokeTypeFromString(s string) PokeType {
	typeMap := map[string]PokeType{
		"一般": TypeNormal, "普通": TypeNormal, "Normal": TypeNormal,
		"火": TypeFire, "Fire": TypeFire,
		"水": TypeWater, "Water": TypeWater,
		"电": TypeElectric, "電": TypeElectric, "Electric": TypeElectric,
		"草": TypeGrass, "Grass": TypeGrass,
		"冰": TypeIce, "Ice": TypeIce,
		"格斗": TypeFighting, "格鬥": TypeFighting, "Fighting": TypeFighting,
		"毒": TypePoison, "Poison": TypePoison,
		"地面": TypeGround, "Ground": TypeGround,
		"飞行": TypeFlying, "飛行": TypeFlying, "Flying": TypeFlying,
		"超能力": TypePsychic, "Psychic": TypePsychic,
		"虫": TypeBug, "蟲": TypeBug, "Bug": TypeBug,
		"岩石": TypeRock, "Rock": TypeRock,
		"幽灵": TypeGhost, "幽靈": TypeGhost, "Ghost": TypeGhost,
		"龙": TypeDragon, "龍": TypeDragon, "Dragon": TypeDragon,
		"恶": TypeDark, "惡": TypeDark, "Dark": TypeDark,
		"钢": TypeSteel, "鋼": TypeSteel, "Steel": TypeSteel,
		"妖精": TypeFairy, "Fairy": TypeFairy,
	}
	if t, ok := typeMap[s]; ok {
		return t
	}
	return TypeNormal
}
