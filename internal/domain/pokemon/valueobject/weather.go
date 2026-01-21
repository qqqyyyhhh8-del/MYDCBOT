package valueobject

// Weather 天气类型
type Weather string

const (
	WeatherNone    Weather = ""
	WeatherSun     Weather = "晴天"
	WeatherRain    Weather = "雨天"
	WeatherSand    Weather = "沙暴"
	WeatherHail    Weather = "冰雹"
	WeatherHarshSun Weather = "大日照" // 原始固拉多
	WeatherHeavyRain Weather = "大雨"  // 原始盖欧卡
	WeatherStrongWinds Weather = "乱流" // 裂空座
)

// WeatherState 天气状态
type WeatherState struct {
	Current   Weather
	TurnsLeft int // 剩余回合数，0表示永久
}

// NewWeatherState 创建天气状态
func NewWeatherState() *WeatherState {
	return &WeatherState{
		Current:   WeatherNone,
		TurnsLeft: 0,
	}
}

// SetWeather 设置天气
func (w *WeatherState) SetWeather(weather Weather, turns int) {
	// 原始天气无法被覆盖
	if w.Current == WeatherHarshSun || w.Current == WeatherHeavyRain || w.Current == WeatherStrongWinds {
		if weather != WeatherNone {
			return
		}
	}
	w.Current = weather
	w.TurnsLeft = turns
}

// Tick 天气回合流逝
func (w *WeatherState) Tick() bool {
	if w.TurnsLeft > 0 {
		w.TurnsLeft--
		if w.TurnsLeft == 0 {
			w.Current = WeatherNone
			return true // 天气结束
		}
	}
	return false
}

// IsActive 天气是否激活
func (w *WeatherState) IsActive() bool {
	return w.Current != WeatherNone
}

// IsSunny 是否晴天
func (w *WeatherState) IsSunny() bool {
	return w.Current == WeatherSun || w.Current == WeatherHarshSun
}

// IsRainy 是否雨天
func (w *WeatherState) IsRainy() bool {
	return w.Current == WeatherRain || w.Current == WeatherHeavyRain
}

// IsSandy 是否沙暴
func (w *WeatherState) IsSandy() bool {
	return w.Current == WeatherSand
}

// IsHailing 是否冰雹
func (w *WeatherState) IsHailing() bool {
	return w.Current == WeatherHail
}

// GetFireModifier 获取火系招式修正
func (w *WeatherState) GetFireModifier() float64 {
	switch w.Current {
	case WeatherSun, WeatherHarshSun:
		return 1.5
	case WeatherRain:
		return 0.5
	case WeatherHeavyRain:
		return 0 // 大雨中火系无效
	default:
		return 1.0
	}
}

// GetWaterModifier 获取水系招式修正
func (w *WeatherState) GetWaterModifier() float64 {
	switch w.Current {
	case WeatherRain, WeatherHeavyRain:
		return 1.5
	case WeatherSun:
		return 0.5
	case WeatherHarshSun:
		return 0 // 大日照中水系无效
	default:
		return 1.0
	}
}

// GetWeatherDamageTypes 获取会受到天气伤害的属性（需要排除的）
func (w *WeatherState) GetWeatherDamageExemptTypes() []PokeType {
	switch w.Current {
	case WeatherSand:
		return []PokeType{TypeRock, TypeGround, TypeSteel}
	case WeatherHail:
		return []PokeType{TypeIce}
	default:
		return nil
	}
}

// CausesWeatherDamage 是否造成天气伤害
func (w *WeatherState) CausesWeatherDamage() bool {
	return w.Current == WeatherSand || w.Current == WeatherHail
}

// GetWeatherName 获取天气名称
func (w *WeatherState) GetWeatherName() string {
	switch w.Current {
	case WeatherSun:
		return "☀️ 日照强烈"
	case WeatherRain:
		return "🌧️ 下起了雨"
	case WeatherSand:
		return "🏜️ 沙暴肆虐"
	case WeatherHail:
		return "❄️ 冰雹来袭"
	case WeatherHarshSun:
		return "🔥 强烈的日照"
	case WeatherHeavyRain:
		return "⛈️ 暴风雨"
	case WeatherStrongWinds:
		return "🌪️ 乱流"
	default:
		return ""
	}
}
