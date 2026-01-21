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
	// 出场触发类
	r.Register(&IntimidateEffect{})
	r.Register(&DrizzleEffect{})
	r.Register(&DroughtEffect{})
	r.Register(&SandStreamEffect{})
	r.Register(&SnowWarningEffect{})
	r.Register(&PressureEffect{})
	r.Register(&UnnerveEffect{})

	// 计算修正类
	r.Register(&HugePowerEffect{})
	r.Register(&PurePowerEffect{})
	r.Register(&TechnicianEffect{})
	r.Register(&ToughClawsEffect{})
	r.Register(&StrongJawEffect{})
	r.Register(&AdaptabilityEffect{})
	r.Register(&SheerForceEffect{})
	r.Register(&OvergrowEffect{})
	r.Register(&BlazeEffect{})
	r.Register(&TorrentEffect{})
	r.Register(&ThickFatEffect{})
	r.Register(&LevitateEffect{})
	r.Register(&WonderGuardEffect{})
	r.Register(&MultiscaleEffect{})
	r.Register(&LightningRodEffect{})

	// 受击触发类
	r.Register(&StaticEffect{})
	r.Register(&CursedBodyEffect{})

	// 状态免疫类
	r.Register(&ImmunityEffect{})
	r.Register(&InnerFocusEffect{})

	// 回合结束类
	r.Register(&SpeedBoostEffect{})

	// 速度修正类
	r.Register(&SwiftSwimEffect{})
	r.Register(&ChlorophyllEffect{})
	r.Register(&SandRushEffect{})
	r.Register(&SlushRushEffect{})

	// 优先度修正类
	r.Register(&PranksterEffect{})
	r.Register(&GaleWingsEffect{})

	// 击倒触发类
	r.Register(&MoxieEffect{})
}
