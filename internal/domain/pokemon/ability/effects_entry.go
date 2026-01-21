package ability

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// ============================================
// 出场触发类特性
// ============================================

// IntimidateEffect 威吓特性
type IntimidateEffect struct {
	BaseEffect
}

func (e *IntimidateEffect) GetAbilityID() int {
	return 22
}

func (e *IntimidateEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *IntimidateEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	if opponent == nil || !opponent.IsAlive() {
		return nil
	}
	// 检查对手是否有内在特性等免疫威吓
	if opponent.GetAbility() != nil {
		switch opponent.GetAbility().ID {
		case 39: // 精神力
			return &EntryResult{
				Messages: []string{"😤 对手的精神力阻止了威吓！"},
			}
		case 52: // 我行我素
			return &EntryResult{
				Messages: []string{"😤 对手的我行我素阻止了威吓！"},
			}
		}
	}
	return &EntryResult{
		Messages:    []string{"😨 威吓降低了对手的攻击！"},
		StatChanges: map[string]int{"atk": -1},
	}
}

// DrizzleEffect 降雨特性
type DrizzleEffect struct {
	BaseEffect
}

func (e *DrizzleEffect) GetAbilityID() int {
	return 2
}

func (e *DrizzleEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *DrizzleEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	rain := valueobject.WeatherRain
	return &EntryResult{
		Messages:   []string{"🌧️ 降雨开始下雨了！"},
		WeatherSet: &rain,
	}
}

// DroughtEffect 日照特性
type DroughtEffect struct {
	BaseEffect
}

func (e *DroughtEffect) GetAbilityID() int {
	return 70
}

func (e *DroughtEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *DroughtEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	sun := valueobject.WeatherSun
	return &EntryResult{
		Messages:   []string{"☀️ 日照变得非常晴朗！"},
		WeatherSet: &sun,
	}
}

// SandStreamEffect 扬沙特性
type SandStreamEffect struct {
	BaseEffect
}

func (e *SandStreamEffect) GetAbilityID() int {
	return 45
}

func (e *SandStreamEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *SandStreamEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	sand := valueobject.WeatherSand
	return &EntryResult{
		Messages:   []string{"🏜️ 扬沙掀起了沙暴！"},
		WeatherSet: &sand,
	}
}

// SnowWarningEffect 降雪特性
type SnowWarningEffect struct {
	BaseEffect
}

func (e *SnowWarningEffect) GetAbilityID() int {
	return 117
}

func (e *SnowWarningEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *SnowWarningEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	hail := valueobject.WeatherHail
	return &EntryResult{
		Messages:   []string{"❄️ 降雪开始下冰雹了！"},
		WeatherSet: &hail,
	}
}

// PressureEffect 压迫感特性
type PressureEffect struct {
	BaseEffect
}

func (e *PressureEffect) GetAbilityID() int {
	return 46
}

func (e *PressureEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *PressureEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	return &EntryResult{
		Messages: []string{"😰 压迫感让对手感到压力！"},
	}
}

// UnnerveEffect 紧张感特性
type UnnerveEffect struct {
	BaseEffect
}

func (e *UnnerveEffect) GetAbilityID() int {
	return 127
}

func (e *UnnerveEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *UnnerveEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	return &EntryResult{
		Messages: []string{"😰 紧张感让对手无法食用树果！"},
	}
}

// DownloadEffect 下载特性
type DownloadEffect struct {
	BaseEffect
}

func (e *DownloadEffect) GetAbilityID() int {
	return 88
}

func (e *DownloadEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *DownloadEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	// 比较对手防御和特防，选择较低的提升对应攻击
	// 这里简化处理，随机选择
	return &EntryResult{
		Messages:    []string{"📥 下载分析了对手的数据！"},
		StatChanges: map[string]int{"spatk": 1}, // 默认提升特攻
	}
}

// FriskEffect 察觉特性
type FriskEffect struct {
	BaseEffect
}

func (e *FriskEffect) GetAbilityID() int {
	return 119
}

func (e *FriskEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *FriskEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	if opponent != nil && opponent.GetItem() != nil {
		return &EntryResult{
			Messages: []string{"👁️ 察觉发现对手携带了 " + opponent.GetItem().Name + "！"},
		}
	}
	return &EntryResult{
		Messages: []string{"👁️ 察觉没有发现对手携带道具。"},
	}
}

// IntrepidSwordEffect 不挠之剑特性
type IntrepidSwordEffect struct {
	BaseEffect
}

func (e *IntrepidSwordEffect) GetAbilityID() int {
	return 234
}

func (e *IntrepidSwordEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *IntrepidSwordEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	return &EntryResult{
		Messages:    []string{"⚔️ 不挠之剑提升了攻击！"},
		StatChanges: map[string]int{"atk": 1},
	}
}

// DauntlessShieldEffect 不屈之盾特性
type DauntlessShieldEffect struct {
	BaseEffect
}

func (e *DauntlessShieldEffect) GetAbilityID() int {
	return 235
}

func (e *DauntlessShieldEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *DauntlessShieldEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	return &EntryResult{
		Messages:    []string{"🛡️ 不屈之盾提升了防御！"},
		StatChanges: map[string]int{"def": 1},
	}
}

// TraceEffect 复制特性
type TraceEffect struct {
	BaseEffect
}

func (e *TraceEffect) GetAbilityID() int {
	return 36
}

func (e *TraceEffect) GetTriggers() []TriggerType {
	return []TriggerType{TriggerOnEntry}
}

func (e *TraceEffect) OnEntry(self Battler, opponent Battler, ctx *BattleContext) *EntryResult {
	if opponent != nil && opponent.GetAbility() != nil {
		return &EntryResult{
			Messages: []string{"🔄 复制了对手的 " + opponent.GetAbility().Name + "！"},
		}
	}
	return nil
}
