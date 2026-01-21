package valueobject

// Item 道具
type Item struct {
	ID          int
	Name        string
	Description string
	Category    ItemCategory
}

// ItemCategory 道具分类
type ItemCategory string

const (
	ItemCategoryHeld      ItemCategory = "held"       // 携带道具
	ItemCategoryBerry     ItemCategory = "berry"      // 树果
	ItemCategoryMega      ItemCategory = "mega"       // 超级石
	ItemCategoryZCrystal  ItemCategory = "z_crystal"  // Z纯晶
	ItemCategoryChoice    ItemCategory = "choice"     // 讲究系列
	ItemCategoryOrb       ItemCategory = "orb"        // 宝珠系列
	ItemCategoryPlate     ItemCategory = "plate"      // 石板
	ItemCategoryMemory    ItemCategory = "memory"     // 存储碟
	ItemCategoryBoost     ItemCategory = "boost"      // 强化道具
)

// 常用道具定义
var (
	// 讲究系列
	ItemChoiceBand   = Item{220, "讲究头带", "攻击x1.5但只能用一个招式", ItemCategoryChoice}
	ItemChoiceSpecs  = Item{297, "讲究眼镜", "特攻x1.5但只能用一个招式", ItemCategoryChoice}
	ItemChoiceScarf  = Item{287, "讲究围巾", "速度x1.5但只能用一个招式", ItemCategoryChoice}

	// 强化道具
	ItemLifeOrb       = Item{270, "生命宝珠", "招式威力x1.3但损失HP", ItemCategoryBoost}
	ItemExpertBelt    = Item{268, "达人带", "效果拔群时威力x1.2", ItemCategoryBoost}
	ItemMuscleBand    = Item{266, "力量头带", "物理招式威力x1.1", ItemCategoryBoost}
	ItemWiseGlasses   = Item{267, "博识眼镜", "特殊招式威力x1.1", ItemCategoryBoost}
	ItemMetronome     = Item{277, "节拍器", "连续使用同招式威力提升", ItemCategoryBoost}

	// 防御道具
	ItemFocusSash     = Item{275, "气势披带", "HP满时必定留1HP", ItemCategoryHeld}
	ItemFocusBand     = Item{230, "气势头带", "10%几率留1HP", ItemCategoryHeld}
	ItemAssaultVest   = Item{640, "突击背心", "特防x1.5但无法用变化招式", ItemCategoryHeld}
	ItemEviolite      = Item{538, "进化奇石", "未进化宝可梦双防x1.5", ItemCategoryHeld}
	ItemRockyHelmet   = Item{540, "凸凸头盔", "被接触时对手损失1/6HP", ItemCategoryHeld}
	ItemLeftovers     = Item{234, "吃剩的东西", "每回合回复1/16HP", ItemCategoryHeld}
	ItemBlackSludge   = Item{281, "黑色污泥", "毒系回复1/16HP否则损失", ItemCategoryHeld}
	ItemSitrusBerry   = Item{158, "文柚果", "HP低于50%时回复25%", ItemCategoryBerry}

	// 属性强化
	ItemTypeBoostFire    = Item{271, "木炭", "火属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostWater   = Item{243, "神秘水滴", "水属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostElec    = Item{242, "磁铁", "电属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostGrass   = Item{237, "奇迹种子", "草属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostIce     = Item{238, "不融冰", "冰属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostFight   = Item{241, "黑带", "格斗招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostPoison  = Item{245, "毒针", "毒属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostGround  = Item{247, "柔软沙子", "地面招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostFlying  = Item{244, "锐利鸟嘴", "飞行招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostPsychic = Item{248, "弯曲汤匙", "超能力招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostBug     = Item{246, "银粉", "虫属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostRock    = Item{ite(249), "硬石头", "岩石招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostGhost   = Item{250, "诅咒护符", "幽灵招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostDragon  = Item{252, "龙之牙", "龙属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostDark    = Item{251, "黑色眼镜", "恶属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostSteel   = Item{253, "金属膜", "钢属性招式威力x1.2", ItemCategoryHeld}
	ItemTypeBoostFairy   = Item{644, "妖精羽毛", "妖精招式威力x1.2", ItemCategoryHeld}

	// 速度道具
	ItemQuickClaw    = Item{217, "先制之爪", "20%几率先制", ItemCategoryHeld}
	ItemIronBall     = Item{278, "黑铁球", "速度减半且飘浮失效", ItemCategoryHeld}
	ItemLaggingTail  = Item{279, "后攻之尾", "必定后攻", ItemCategoryHeld}

	// 特殊道具
	ItemHeavyDutyBoots = Item{1120, "厚底靴", "免疫入场伤害", ItemCategoryHeld}
	ItemSafetyGoggles  = Item{650, "防尘护目镜", "免疫天气伤害和粉末", ItemCategoryHeld}
	ItemAirBalloon     = Item{541, "气球", "免疫地面直到被攻击", ItemCategoryHeld}
	ItemRedCard        = Item{542, "红牌", "被攻击时强制对手交换", ItemCategoryHeld}
	ItemEjectButton    = Item{547, "逃脱按键", "被攻击时可交换", ItemCategoryHeld}
	ItemShedShell      = Item{295, "脱壳", "可无视束缚交换", ItemCategoryHeld}

	// 宝珠系列
	ItemAdamantOrb  = Item{135, "金刚宝珠", "帝牙卢卡龙钢招式x1.2", ItemCategoryOrb}
	ItemLustrousOrb = Item{136, "白玉宝珠", "帕路奇亚龙水招式x1.2", ItemCategoryOrb}
	ItemGriseousOrb = Item{112, "白金宝珠", "骑拉帝纳龙幽灵招式x1.2", ItemCategoryOrb}
	ItemSoulDew     = Item{225, "心之水滴", "拉帝特攻特防x1.5", ItemCategoryOrb}
)

func ite(n int) int { return n }

// ItemMap 道具ID映射
var ItemMap = map[int]Item{
	220: ItemChoiceBand,
	297: ItemChoiceSpecs,
	287: ItemChoiceScarf,
	270: ItemLifeOrb,
	268: ItemExpertBelt,
	275: ItemFocusSash,
	640: ItemAssaultVest,
	538: ItemEviolite,
	540: ItemRockyHelmet,
	234: ItemLeftovers,
	281: ItemBlackSludge,
	158: ItemSitrusBerry,
	217: ItemQuickClaw,
	1120: ItemHeavyDutyBoots,
	541: ItemAirBalloon,
}

// GetItemByID 通过ID获取道具
func GetItemByID(id int) *Item {
	if item, ok := ItemMap[id]; ok {
		return &item
	}
	return nil
}

// GetItemByName 通过名称获取道具
func GetItemByName(name string) *Item {
	for _, item := range ItemMap {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

// CommonHeldItems 常用携带道具列表
var CommonHeldItems = []Item{
	ItemChoiceBand, ItemChoiceSpecs, ItemChoiceScarf,
	ItemLifeOrb, ItemExpertBelt,
	ItemFocusSash, ItemAssaultVest, ItemEviolite,
	ItemRockyHelmet, ItemLeftovers, ItemBlackSludge,
	ItemSitrusBerry, ItemHeavyDutyBoots, ItemAirBalloon,
}
