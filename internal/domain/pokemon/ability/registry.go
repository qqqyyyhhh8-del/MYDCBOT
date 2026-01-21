package ability

import (
	"sync"
)

// Registry 特性效果注册表
type Registry struct {
	effects map[int]Effect
	mu      sync.RWMutex
}

var (
	globalRegistry *Registry
	once           sync.Once
)

// GetRegistry 获取全局注册表
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			effects: make(map[int]Effect),
		}
		// 注册所有特性效果
		registerAllEffects(globalRegistry)
	})
	return globalRegistry
}

// Register 注册特性效果
func (r *Registry) Register(effect Effect) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.effects[effect.GetAbilityID()] = effect
}

// Get 获取特性效果
func (r *Registry) Get(abilityID int) Effect {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.effects[abilityID]
}

// Has 检查是否有特性效果
func (r *Registry) Has(abilityID int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.effects[abilityID]
	return ok
}

// GetAll 获取所有已注册的特性效果
func (r *Registry) GetAll() []Effect {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Effect, 0, len(r.effects))
	for _, e := range r.effects {
		result = append(result, e)
	}
	return result
}

// registerAllEffects 注册所有特性效果
func registerAllEffects(r *Registry) {
	// ============================================
	// 出场触发类
	// ============================================
	r.Register(&IntimidateEffect{})      // 22 威吓
	r.Register(&DrizzleEffect{})         // 2 降雨
	r.Register(&DroughtEffect{})         // 70 日照
	r.Register(&SandStreamEffect{})      // 45 扬沙
	r.Register(&SnowWarningEffect{})     // 117 降雪
	r.Register(&PressureEffect{})        // 46 压迫感
	r.Register(&UnnerveEffect{})         // 127 紧张感
	r.Register(&DownloadEffect{})        // 88 下载
	r.Register(&FriskEffect{})           // 119 察觉
	r.Register(&IntrepidSwordEffect{})   // 234 不挠之剑
	r.Register(&DauntlessShieldEffect{}) // 235 不屈之盾
	r.Register(&TraceEffect{})           // 36 复制

	// ============================================
	// 计算修正类（攻击方）
	// ============================================
	r.Register(&HugePowerEffect{})      // 37 大力士
	r.Register(&PurePowerEffect{})      // 74 瑜伽之力
	r.Register(&TechnicianEffect{})     // 101 技术高手
	r.Register(&ToughClawsEffect{})     // 181 硬爪
	r.Register(&StrongJawEffect{})      // 173 强壮之颚
	r.Register(&AdaptabilityEffect{})   // 91 适应力
	r.Register(&SheerForceEffect{})     // 125 强行
	r.Register(&OvergrowEffect{})       // 65 茂盛
	r.Register(&BlazeEffect{})          // 66 猛火
	r.Register(&TorrentEffect{})        // 67 激流
	r.Register(&SwarmEffect{})          // 68 虫之预感
	r.Register(&IronFistEffect{})       // 89 铁拳
	r.Register(&SniperEffect{})         // 97 狙击手
	r.Register(&SandForceEffect{})      // 159 沙之力
	r.Register(&TintedLensEffect{})     // 110 有色眼镜
	r.Register(&NeuroforceEffect{})     // 233 脑核之力
	r.Register(&RecklessEffect{})       // 120 舍身
	r.Register(&MegaLauncherEffect{})   // 178 超级发射器
	r.Register(&SteelworkerEffect{})    // 200 钢能力者

	// ============================================
	// 计算修正类（防御方）
	// ============================================
	r.Register(&ThickFatEffect{})       // 47 厚脂肪
	r.Register(&LevitateEffect{})       // 26 飘浮
	r.Register(&WonderGuardEffect{})    // 25 神奇守护
	r.Register(&MultiscaleEffect{})     // 136 多重鳞片
	r.Register(&LightningRodEffect{})   // 31 避雷针
	r.Register(&VoltAbsorbEffect{})     // 10 蓄电
	r.Register(&WaterAbsorbEffect{})    // 11 储水
	r.Register(&FlashFireEffect{})      // 18 引火
	r.Register(&FurCoatEffect{})        // 169 毛皮大衣
	r.Register(&HeatproofEffect{})      // 85 耐热
	r.Register(&DrySkinEffect{})        // 87 干燥皮肤
	r.Register(&StormDrainEffect{})     // 114 引水
	r.Register(&SapSipperEffect{})      // 157 食草
	r.Register(&MotorDriveEffect{})     // 78 电气引擎
	r.Register(&SolidRockEffect{})      // 116 坚硬岩石
	r.Register(&FilterEffect{})         // 111 过滤
	r.Register(&PrismArmorEffect{})     // 232 棱镜装甲

	// ============================================
	// 受击触发类
	// ============================================
	r.Register(&StaticEffect{})         // 9 静电
	r.Register(&CursedBodyEffect{})     // 130 诅咒之躯
	r.Register(&StenchEffect{})         // 1 恶臭
	r.Register(&PoisonPointEffect{})    // 38 毒刺
	r.Register(&FlameBodyEffect{})      // 49 火焰之躯
	r.Register(&RoughSkinEffect{})      // 24 粗糙皮肤
	r.Register(&EffectSporeEffect{})    // 27 孢子
	r.Register(&IronBarbsEffect{})      // 160 铁刺
	r.Register(&CuteCharmEffect{})      // 56 迷人之躯
	r.Register(&GooeyEffect{})          // 183 黏滑
	r.Register(&TanglingHairEffect{})   // 221 卷发
	r.Register(&MummyEffect{})          // 152 木乃伊

	// ============================================
	// 状态免疫类
	// ============================================
	r.Register(&ImmunityEffect{})       // 17 免疫
	r.Register(&InnerFocusEffect{})     // 39 精神力
	r.Register(&LimberEffect{})         // 7 柔软
	r.Register(&InsomniaEffect{})       // 15 不眠
	r.Register(&VitalSpiritEffect{})    // 72 干劲
	r.Register(&MagmaArmorEffect{})     // 40 熔岩铠甲
	r.Register(&WaterVeilEffect{})      // 41 水幕
	r.Register(&OwnTempoEffect{})       // 20 我行我素
	r.Register(&ObliviousEffect{})      // 12 迟钝

	// ============================================
	// 回合结束类
	// ============================================
	r.Register(&SpeedBoostEffect{})     // 3 加速
	r.Register(&RainDishEffect{})       // 44 雨盘
	r.Register(&IceBodyEffect{})        // 115 冰冻之躯
	r.Register(&ShedSkinEffect{})       // 61 蜕皮
	r.Register(&PoisonHealEffect{})     // 90 毒疗
	r.Register(&SolarPowerEffect{})     // 94 太阳之力

	// ============================================
	// 速度修正类
	// ============================================
	r.Register(&SwiftSwimEffect{})      // 33 悠游自如
	r.Register(&ChlorophyllEffect{})    // 34 叶绿素
	r.Register(&SandRushEffect{})       // 146 拨沙
	r.Register(&SlushRushEffect{})      // 202 拨雪

	// ============================================
	// 优先度修正类
	// ============================================
	r.Register(&PranksterEffect{})      // 158 恶作剧之心
	r.Register(&GaleWingsEffect{})      // 177 疾风之翼

	// ============================================
	// 击倒触发类
	// ============================================
	r.Register(&MoxieEffect{})          // 153 自信过剩
	r.Register(&BeastBoostEffect{})     // 224 异兽提升
	r.Register(&SoulHeartEffect{})      // 220 魂心

	// ============================================
	// 形态变化类
	// ============================================
	r.Register(&BattleBondEffect{})     // 210 羁绊变身
	r.Register(&ZenModeEffect{})        // 161 达摩模式
	r.Register(&PowerConstructEffect{}) // 211 群聚变形
	r.Register(&StanceChangeEffect{})   // 176 战斗切换
}
