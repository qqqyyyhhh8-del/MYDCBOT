package valueobject

// Nature 性格
type Nature string

const (
	NatureHardy   Nature = "勤奋" // 无修正
	NatureLonely  Nature = "怕寂寞" // +攻击 -防御
	NatureBrave   Nature = "勇敢" // +攻击 -速度
	NatureAdamant Nature = "固执" // +攻击 -特攻
	NatureNaughty Nature = "顽皮" // +攻击 -特防

	NatureBold    Nature = "大胆" // +防御 -攻击
	NatureDocile  Nature = "坦率" // 无修正
	NatureRelaxed Nature = "悠闲" // +防御 -速度
	NatureImpish  Nature = "淘气" // +防御 -特攻
	NatureLax     Nature = "乐天" // +防御 -特防

	NatureTimid  Nature = "胆小" // +速度 -攻击
	NatureHasty  Nature = "急躁" // +速度 -防御
	NatureSerious Nature = "认真" // 无修正
	NatureJolly  Nature = "爽朗" // +速度 -特攻
	NatureNaive  Nature = "天真" // +速度 -特防

	NatureModest Nature = "内敛" // +特攻 -攻击
	NatureMild   Nature = "慢吞吞" // +特攻 -防御
	NatureQuiet  Nature = "冷静" // +特攻 -速度
	NatureBashful Nature = "害羞" // 无修正
	NatureRash   Nature = "马虎" // +特攻 -特防

	NatureCalm    Nature = "温和" // +特防 -攻击
	NatureGentle  Nature = "温顺" // +特防 -防御
	NatureSassy   Nature = "自大" // +特防 -速度
	NatureCareful Nature = "慎重" // +特防 -特攻
	NatureQuirky  Nature = "浮躁" // 无修正
)

// NatureModifier 性格修正
type NatureModifier struct {
	Atk   float64
	Def   float64
	SpAtk float64
	SpDef float64
	Speed float64
}

// GetNatureModifier 获取性格修正值
func GetNatureModifier(nature Nature) NatureModifier {
	modifiers := map[Nature]NatureModifier{
		// 无修正
		NatureHardy:   {1.0, 1.0, 1.0, 1.0, 1.0},
		NatureDocile:  {1.0, 1.0, 1.0, 1.0, 1.0},
		NatureSerious: {1.0, 1.0, 1.0, 1.0, 1.0},
		NatureBashful: {1.0, 1.0, 1.0, 1.0, 1.0},
		NatureQuirky:  {1.0, 1.0, 1.0, 1.0, 1.0},

		// +攻击
		NatureLonely:  {1.1, 0.9, 1.0, 1.0, 1.0},
		NatureBrave:   {1.1, 1.0, 1.0, 1.0, 0.9},
		NatureAdamant: {1.1, 1.0, 0.9, 1.0, 1.0},
		NatureNaughty: {1.1, 1.0, 1.0, 0.9, 1.0},

		// +防御
		NatureBold:    {0.9, 1.1, 1.0, 1.0, 1.0},
		NatureRelaxed: {1.0, 1.1, 1.0, 1.0, 0.9},
		NatureImpish:  {1.0, 1.1, 0.9, 1.0, 1.0},
		NatureLax:     {1.0, 1.1, 1.0, 0.9, 1.0},

		// +速度
		NatureTimid: {0.9, 1.0, 1.0, 1.0, 1.1},
		NatureHasty: {1.0, 0.9, 1.0, 1.0, 1.1},
		NatureJolly: {1.0, 1.0, 0.9, 1.0, 1.1},
		NatureNaive: {1.0, 1.0, 1.0, 0.9, 1.1},

		// +特攻
		NatureModest: {0.9, 1.0, 1.1, 1.0, 1.0},
		NatureMild:   {1.0, 0.9, 1.1, 1.0, 1.0},
		NatureQuiet:  {1.0, 1.0, 1.1, 1.0, 0.9},
		NatureRash:   {1.0, 1.0, 1.1, 0.9, 1.0},

		// +特防
		NatureCalm:    {0.9, 1.0, 1.0, 1.1, 1.0},
		NatureGentle:  {1.0, 0.9, 1.0, 1.1, 1.0},
		NatureSassy:   {1.0, 1.0, 1.0, 1.1, 0.9},
		NatureCareful: {1.0, 1.0, 0.9, 1.1, 1.0},
	}

	if mod, ok := modifiers[nature]; ok {
		return mod
	}
	return NatureModifier{1.0, 1.0, 1.0, 1.0, 1.0}
}

// AllNatures 所有性格列表
var AllNatures = []Nature{
	NatureHardy, NatureLonely, NatureBrave, NatureAdamant, NatureNaughty,
	NatureBold, NatureDocile, NatureRelaxed, NatureImpish, NatureLax,
	NatureTimid, NatureHasty, NatureSerious, NatureJolly, NatureNaive,
	NatureModest, NatureMild, NatureQuiet, NatureBashful, NatureRash,
	NatureCalm, NatureGentle, NatureSassy, NatureCareful, NatureQuirky,
}

// GetNatureDescription 获取性格描述
func (n Nature) Description() string {
	mod := GetNatureModifier(n)
	if mod.Atk == 1.0 && mod.Def == 1.0 && mod.SpAtk == 1.0 && mod.SpDef == 1.0 && mod.Speed == 1.0 {
		return string(n) + " (无修正)"
	}

	desc := string(n) + " ("
	if mod.Atk > 1 {
		desc += "+攻击 "
	} else if mod.Atk < 1 {
		desc += "-攻击 "
	}
	if mod.Def > 1 {
		desc += "+防御 "
	} else if mod.Def < 1 {
		desc += "-防御 "
	}
	if mod.SpAtk > 1 {
		desc += "+特攻 "
	} else if mod.SpAtk < 1 {
		desc += "-特攻 "
	}
	if mod.SpDef > 1 {
		desc += "+特防 "
	} else if mod.SpDef < 1 {
		desc += "-特防 "
	}
	if mod.Speed > 1 {
		desc += "+速度"
	} else if mod.Speed < 1 {
		desc += "-速度"
	}
	desc += ")"
	return desc
}
