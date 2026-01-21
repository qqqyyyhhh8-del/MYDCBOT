package entity

import (
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// ============================================
// Battler 接口方法实现（用于特性系统）
// ============================================

// GetAbility 获取特性
func (b *Battler) GetAbility() *valueobject.Ability {
	return b.Ability
}

// GetTypes 获取当前属性
func (b *Battler) GetTypes() []valueobject.PokeType {
	return b.Types
}

// GetCurrentHP 获取当前HP
func (b *Battler) GetCurrentHP() int {
	return b.CurrentHP
}

// GetMaxHP 获取最大HP
func (b *Battler) GetMaxHP() int {
	return b.MaxHP
}

// GetStatus 获取异常状态
func (b *Battler) GetStatus() string {
	return string(b.Status)
}

// SetStatus 设置异常状态
func (b *Battler) SetStatus(status string) {
	b.Status = StatusCondition(status)
}

// HasVolatile 检查是否有临时状态
func (b *Battler) HasVolatile(status string) bool {
	for _, v := range b.Volatile {
		if string(v) == status {
			return true
		}
	}
	return false
}

// AddVolatile 添加临时状态
func (b *Battler) AddVolatile(status string) {
	if !b.HasVolatile(status) {
		b.Volatile = append(b.Volatile, VolatileStatus(status))
	}
}

// RemoveVolatile 移除临时状态
func (b *Battler) RemoveVolatile(status string) {
	for i, v := range b.Volatile {
		if string(v) == status {
			b.Volatile = append(b.Volatile[:i], b.Volatile[i+1:]...)
			return
		}
	}
}

// GetItem 获取道具
func (b *Battler) GetItem() *valueobject.Item {
	return b.Item
}

// IsItemConsumed 道具是否已消耗
func (b *Battler) IsItemConsumed() bool {
	return b.ItemConsumed
}

// ConsumeItem 消耗道具
func (b *Battler) ConsumeItem() {
	b.ItemConsumed = true
}

// ============================================
// Move 接口方法实现（用于特性系统）
// ============================================

// MoveAdapter Move适配器，实现ability.Move接口
type MoveAdapter struct {
	Move *Move
}

// NewMoveAdapter 创建Move适配器
func NewMoveAdapter(m *Move) *MoveAdapter {
	return &MoveAdapter{Move: m}
}

// GetName 获取技能名称
func (m *MoveAdapter) GetName() string {
	return m.Move.Name
}

// GetType 获取技能属性
func (m *MoveAdapter) GetType() valueobject.PokeType {
	return m.Move.Type
}

// GetCategory 获取技能分类
func (m *MoveAdapter) GetCategory() string {
	switch m.Move.Category {
	case CategoryPhysical:
		return "physical"
	case CategorySpecial:
		return "special"
	case CategoryStatus:
		return "status"
	}
	return ""
}

// GetPower 获取技能威力
func (m *MoveAdapter) GetPower() int {
	return m.Move.Power
}

// GetPriority 获取技能优先度
func (m *MoveAdapter) GetPriority() int {
	return m.Move.Priority
}

// IsContact 是否为接触技能
func (m *MoveAdapter) IsContact() bool {
	return m.Move.MakesContact
}

// IsBite 是否为咬类技能
func (m *MoveAdapter) IsBite() bool {
	// 咬类技能列表
	biteMovesMap := map[string]bool{
		"咬住":   true,
		"嘎吱嘎吱": true,
		"火焰牙":  true,
		"雷电牙":  true,
		"冰冻牙":  true,
		"剧毒牙":  true,
		"强力牙":  true,
	}
	return biteMovesMap[m.Move.Name]
}

// IsPunch 是否为拳类技能
func (m *MoveAdapter) IsPunch() bool {
	// 拳类技能列表
	punchMovesMap := map[string]bool{
		"音速拳":  true,
		"火焰拳":  true,
		"雷电拳":  true,
		"冰冻拳":  true,
		"吸取拳":  true,
		"真气拳":  true,
		"子弹拳":  true,
		"彗星拳":  true,
		"影子拳":  true,
		"天空上升拳": true,
		"蓄能拳":  true,
	}
	return punchMovesMap[m.Move.Name]
}

// IsSound 是否为声音技能
func (m *MoveAdapter) IsSound() bool {
	soundMovesMap := map[string]bool{
		"吼叫":   true,
		"唱歌":   true,
		"超音波":  true,
		"刺耳声":  true,
		"金属音":  true,
		"虫鸣":   true,
		"回声":   true,
		"轮唱":   true,
		"爆音波":  true,
		"精神冲击": true,
	}
	return soundMovesMap[m.Move.Name]
}

// IsBullet 是否为子弹/球类技能
func (m *MoveAdapter) IsBullet() bool {
	bulletMovesMap := map[string]bool{
		"种子机关枪": true,
		"气象球":   true,
		"暗影球":   true,
		"能量球":   true,
		"电球":    true,
		"淤泥炸弹":  true,
		"污泥波":   true,
		"岩石炮":   true,
		"磁力炸弹":  true,
		"花粉团":   true,
	}
	return bulletMovesMap[m.Move.Name]
}
