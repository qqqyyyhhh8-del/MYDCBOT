package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/user/dcminigames/internal/domain/pokemon/ability"
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

// BattleState 对战状态
type BattleState string

const (
	BattleStateWaiting  BattleState = "waiting"  // 等待对手
	BattleStateChoosing BattleState = "choosing" // 选择宝可梦
	BattleStateBattling BattleState = "battling" // 对战中
	BattleStateFinished BattleState = "finished" // 已结束
)

// TeamSize 队伍大小（对战模式）
type TeamSize int

const (
	TeamSize1v1 TeamSize = 1 // 单挑模式
	TeamSize3v3 TeamSize = 3 // 3v3 单打
	TeamSize6v6 TeamSize = 6 // 6v6 单打
)

// GetDisplayName 获取模式显示名称
func (t TeamSize) GetDisplayName() string {
	switch t {
	case TeamSize1v1:
		return "单挑 (1v1)"
	case TeamSize3v3:
		return "3v3 单打"
	case TeamSize6v6:
		return "6v6 单打"
	default:
		return fmt.Sprintf("%dv%d", t, t)
	}
}

// AI 玩家常量
const (
	AIPlayerID   = "AI_PLAYER"
	AIPlayerName = "🤖 AI训练师"
)

// Battle 对战实体
type Battle struct {
	ID             string
	ChannelID      string
	Player1        *BattlePlayer
	Player2        *BattlePlayer
	CurrentTurn    int
	State          BattleState
	Winner         *BattlePlayer
	Logs           []string
	CreatedAt      time.Time
	TeamSize       TeamSize              // 队伍大小
	IsAIBattle     bool                  // 是否为人机对战
	Weather        valueobject.Weather   // 当前天气
	WeatherTurns   int                   // 天气剩余回合
	Terrain        string                // 当前场地
	TerrainTurns   int                   // 场地剩余回合
	AbilityService *ability.Service      // 特性服务
}

// BattlePlayer 对战玩家
type BattlePlayer struct {
	ID            string
	Username      string
	Pokemon       *Battler   // 当前出战的宝可梦
	Team          []*Battler // 宝可梦队伍
	ActiveIndex   int        // 当前出战宝可梦在队伍中的索引
	Ready         bool
	Action        *BattleAction
	SelectingSlot int // 当前正在选择的队伍槽位 (0-5)
}

// HasSwitchableTeamMember 检查是否有可换上场的队友
func (p *BattlePlayer) HasSwitchableTeamMember() bool {
	for idx, battler := range p.Team {
		if idx != p.ActiveIndex && battler.IsAlive() {
			return true
		}
	}
	return false
}

// GetActivePokemon 获取当前出战的宝可梦
func (p *BattlePlayer) GetActivePokemon() *Battler {
	if p.ActiveIndex >= 0 && p.ActiveIndex < len(p.Team) {
		return p.Team[p.ActiveIndex]
	}
	return p.Pokemon
}

// HasAlivePokemon 检查是否还有存活的宝可梦
func (p *BattlePlayer) HasAlivePokemon() bool {
	for _, battler := range p.Team {
		if battler.IsAlive() {
			return true
		}
	}
	return false
}

// BattleAction 对战行动
type BattleAction struct {
	Type        ActionType
	MoveIndex   int // 技能索引 (ActionMove 时使用)
	SwitchIndex int // 换人目标索引 (ActionSwitch 时使用)
}

// ActionType 行动类型
type ActionType string

const (
	ActionMove    ActionType = "move"
	ActionSwitch  ActionType = "switch"
	ActionForfeit ActionType = "forfeit"
)

// NewBattle 创建对战 (默认1v1)
func NewBattle(id, channelID string) *Battle {
	return NewBattleWithTeamSize(id, channelID, TeamSize1v1)
}

// NewBattleWithTeamSize 创建指定队伍大小的对战
func NewBattleWithTeamSize(id, channelID string, teamSize TeamSize) *Battle {
	return &Battle{
		ID:             id,
		ChannelID:      channelID,
		CurrentTurn:    1,
		State:          BattleStateWaiting,
		Logs:           make([]string, 0),
		CreatedAt:      time.Now(),
		TeamSize:       teamSize,
		IsAIBattle:     false,
		Weather:        valueobject.WeatherNone,
		AbilityService: ability.NewService(),
	}
}

// NewAIBattle 创建人机对战
func NewAIBattle(id, channelID string, teamSize TeamSize) *Battle {
	return &Battle{
		ID:             id,
		ChannelID:      channelID,
		CurrentTurn:    1,
		State:          BattleStateWaiting,
		Logs:           make([]string, 0),
		CreatedAt:      time.Now(),
		TeamSize:       teamSize,
		IsAIBattle:     true,
		Weather:        valueobject.WeatherNone,
		AbilityService: ability.NewService(),
	}
}

// IsAIPlayer 检查是否为 AI 玩家
func (b *Battle) IsAIPlayer(playerID string) bool {
	return playerID == AIPlayerID
}

// GetAIPlayer 获取 AI 玩家
func (b *Battle) GetAIPlayer() *BattlePlayer {
	if !b.IsAIBattle {
		return nil
	}
	if b.Player1 != nil && b.Player1.ID == AIPlayerID {
		return b.Player1
	}
	if b.Player2 != nil && b.Player2.ID == AIPlayerID {
		return b.Player2
	}
	return nil
}

// GetHumanPlayer 获取人类玩家
func (b *Battle) GetHumanPlayer() *BattlePlayer {
	if !b.IsAIBattle {
		return nil
	}
	if b.Player1 != nil && b.Player1.ID != AIPlayerID {
		return b.Player1
	}
	if b.Player2 != nil && b.Player2.ID != AIPlayerID {
		return b.Player2
	}
	return nil
}

// AddPlayer 添加玩家
func (b *Battle) AddPlayer(playerID, username string) error {
	if b.State != BattleStateWaiting && b.State != BattleStateChoosing {
		return errors.New("对战已开始")
	}
	if b.Player1 != nil && b.Player1.ID == playerID {
		return errors.New("你已在对战中")
	}
	if b.Player2 != nil && b.Player2.ID == playerID {
		return errors.New("你已在对战中")
	}

	player := &BattlePlayer{
		ID:            playerID,
		Username:      username,
		Team:          make([]*Battler, 0, int(b.TeamSize)),
		Ready:         false,
		SelectingSlot: 0,
	}

	if b.Player1 == nil {
		b.Player1 = player
	} else if b.Player2 == nil {
		b.Player2 = player
		b.State = BattleStateChoosing
	} else {
		return errors.New("对战已满员")
	}
	return nil
}

// SetPokemon 设置玩家的宝可梦（添加到队伍）
func (b *Battle) SetPokemon(playerID string, pokemon *Pokemon, level int) error {
	if b.State != BattleStateChoosing {
		return errors.New("当前不能选择宝可梦")
	}

	player := b.GetPlayer(playerID)
	if player == nil {
		return errors.New("你不在对战中")
	}

	// 检查队伍是否已满
	if len(player.Team) >= int(b.TeamSize) {
		return errors.New("队伍已满")
	}

	// 检查是否已经选择了相同的宝可梦（种族条款）
	for _, battler := range player.Team {
		if battler.Pokemon.ID == pokemon.ID {
			return errors.New("不能选择重复的宝可梦")
		}
	}

	battler := NewBattler(pokemon, level)
	player.Team = append(player.Team, battler)
	player.SelectingSlot++

	// 第一只宝可梦自动设为当前出战
	if len(player.Team) == 1 {
		player.Pokemon = battler
	}

	// 检查队伍是否已满
	if len(player.Team) >= int(b.TeamSize) {
		player.Ready = true
	}

	// 检查是否双方都准备好了
	if b.Player1 != nil && b.Player1.Ready && b.Player2 != nil && b.Player2.Ready {
		b.State = BattleStateBattling
		b.Logs = append(b.Logs, "⚔️ 对战开始！")
	}

	return nil
}

// GetRemainingSlots 获取剩余可选槽位数
func (b *Battle) GetRemainingSlots(playerID string) int {
	player := b.GetPlayer(playerID)
	if player == nil {
		return 0
	}
	return int(b.TeamSize) - len(player.Team)
}

// GetTeamStatus 获取队伍状态描述
func (p *BattlePlayer) GetTeamStatus() string {
	if len(p.Team) == 0 {
		return "未选择"
	}
	var names []string
	for _, battler := range p.Team {
		if battler.IsAlive() {
			names = append(names, battler.Pokemon.Name)
		} else {
			names = append(names, "💀"+battler.Pokemon.Name)
		}
	}
	return fmt.Sprintf("%d只: %s", len(p.Team), joinStrings(names, ", "))
}

// GetAliveCount 获取存活宝可梦数量
func (p *BattlePlayer) GetAliveCount() int {
	count := 0
	for _, battler := range p.Team {
		if battler.IsAlive() {
			count++
		}
	}
	return count
}

// HasAlive 是否还有存活的宝可梦
func (p *BattlePlayer) HasAlive() bool {
	return p.GetAliveCount() > 0
}

// GetNextAlive 获取下一只存活的宝可梦
func (p *BattlePlayer) GetNextAlive() *Battler {
	for _, battler := range p.Team {
		if battler.IsAlive() && battler != p.Pokemon {
			return battler
		}
	}
	return nil
}

// joinStrings 连接字符串切片
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// GetPlayer 获取玩家
func (b *Battle) GetPlayer(playerID string) *BattlePlayer {
	if b.Player1 != nil && b.Player1.ID == playerID {
		return b.Player1
	}
	if b.Player2 != nil && b.Player2.ID == playerID {
		return b.Player2
	}
	return nil
}

// GetOpponent 获取对手
func (b *Battle) GetOpponent(playerID string) *BattlePlayer {
	if b.Player1 != nil && b.Player1.ID == playerID {
		return b.Player2
	}
	if b.Player2 != nil && b.Player2.ID == playerID {
		return b.Player1
	}
	return nil
}

// SetAction 设置行动
func (b *Battle) SetAction(playerID string, action *BattleAction) error {
	if b.State != BattleStateBattling {
		return errors.New("对战未开始")
	}

	player := b.GetPlayer(playerID)
	if player == nil {
		return errors.New("你不在对战中")
	}

	if action.Type == ActionMove {
		if action.MoveIndex < 0 || action.MoveIndex >= len(player.Pokemon.Moves) {
			return errors.New("无效的技能")
		}
		move := player.Pokemon.Moves[action.MoveIndex]
		if !move.CanUse() {
			return errors.New("PP不足")
		}
	}

	player.Action = action
	return nil
}

// BothActionsReady 双方都已选择行动
func (b *Battle) BothActionsReady() bool {
	return b.Player1 != nil && b.Player1.Action != nil &&
		b.Player2 != nil && b.Player2.Action != nil
}

// ExecuteTurn 执行回合
func (b *Battle) ExecuteTurn() []string {
	if !b.BothActionsReady() {
		return nil
	}

	logs := make([]string, 0)
	logs = append(logs, "")
	logs = append(logs, "━━━━━━━━━━━━━━━━")
	logs = append(logs, "📍 **回合 "+itoa(b.CurrentTurn)+"**")

	// 检查认输
	if b.Player1.Action.Type == ActionForfeit {
		b.Winner = b.Player2
		b.State = BattleStateFinished
		logs = append(logs, "🏳️ "+b.Player1.Username+" 认输了！")
		logs = append(logs, "🏆 "+b.Player2.Username+" 获胜！")
		b.Logs = append(b.Logs, logs...)
		return logs
	}
	if b.Player2.Action.Type == ActionForfeit {
		b.Winner = b.Player1
		b.State = BattleStateFinished
		logs = append(logs, "🏳️ "+b.Player2.Username+" 认输了！")
		logs = append(logs, "🏆 "+b.Player1.Username+" 获胜！")
		b.Logs = append(b.Logs, logs...)
		return logs
	}

	// 处理换人（换人优先于攻击）
	if b.Player1.Action.Type == ActionSwitch {
		switchLogs := b.executeSwitch(b.Player1)
		logs = append(logs, switchLogs...)
	}
	if b.Player2.Action.Type == ActionSwitch {
		switchLogs := b.executeSwitch(b.Player2)
		logs = append(logs, switchLogs...)
	}

	// 确定行动顺序（优先度 > 速度）
	p1Priority := 0
	p2Priority := 0
	
	// 获取技能优先度
	if b.Player1.Action.Type == ActionMove && b.Player1.Action.MoveIndex < len(b.Player1.Pokemon.Moves) {
		p1Priority = b.Player1.Pokemon.Moves[b.Player1.Action.MoveIndex].Priority
	}
	if b.Player2.Action.Type == ActionMove && b.Player2.Action.MoveIndex < len(b.Player2.Pokemon.Moves) {
		p2Priority = b.Player2.Pokemon.Moves[b.Player2.Action.MoveIndex].Priority
	}
	
	first, second := b.Player1, b.Player2
	// 先比较优先度，再比较有效速度
	if p2Priority > p1Priority {
		first, second = b.Player2, b.Player1
	} else if p2Priority == p1Priority && b.Player2.Pokemon.GetEffectiveSpeed() > b.Player1.Pokemon.GetEffectiveSpeed() {
		first, second = b.Player2, b.Player1
	}

	// 先手行动
	if first.Action.Type == ActionMove {
		actionLogs := b.executeAction(first, second)
		logs = append(logs, actionLogs...)
	}

	// 检查后手宝可梦是否倒下
	if !second.Pokemon.IsAlive() {
		logs = append(logs, "💀 "+second.Pokemon.Pokemon.Name+" 倒下了！")
		// 检查是否还有存活的宝可梦
		if !second.HasAlive() {
			b.Winner = first
			b.State = BattleStateFinished
			logs = append(logs, "🏆 "+first.Username+" 获胜！")
			b.Logs = append(b.Logs, logs...)
			b.clearActions()
			return logs
		}
		// 自动换上下一只宝可梦
		nextPokemon := second.GetNextAlive()
		if nextPokemon != nil {
			second.Pokemon = nextPokemon
			logs = append(logs, "🔄 "+second.Username+" 派出了 "+nextPokemon.Pokemon.Name+"！")
		}
	}

	// 后手行动（如果还存活）
	if second.Pokemon.IsAlive() && second.Action.Type == ActionMove {
		actionLogs := b.executeAction(second, first)
		logs = append(logs, actionLogs...)
	}

	// 检查先手宝可梦是否倒下
	if !first.Pokemon.IsAlive() {
		logs = append(logs, "💀 "+first.Pokemon.Pokemon.Name+" 倒下了！")
		// 检查是否还有存活的宝可梦
		if !first.HasAlive() {
			b.Winner = second
			b.State = BattleStateFinished
			logs = append(logs, "🏆 "+second.Username+" 获胜！")
			b.Logs = append(b.Logs, logs...)
			b.clearActions()
			return logs
		}
		// 自动换上下一只宝可梦
		nextPokemon := first.GetNextAlive()
		if nextPokemon != nil {
			first.Pokemon = nextPokemon
			logs = append(logs, "🔄 "+first.Username+" 派出了 "+nextPokemon.Pokemon.Name+"！")
		}
	}

	// 回合结束特性触发
	turnEndLogs := b.TriggerTurnEndAbilities()
	logs = append(logs, turnEndLogs...)

	b.CurrentTurn++
	b.Logs = append(b.Logs, logs...)
	b.clearActions()
	return logs
}

// executeSwitch 执行换人
func (b *Battle) executeSwitch(player *BattlePlayer) []string {
	logs := make([]string, 0)
	if player.Action.SwitchIndex < 0 || player.Action.SwitchIndex >= len(player.Team) {
		return logs
	}
	newPokemon := player.Team[player.Action.SwitchIndex]
	if !newPokemon.IsAlive() || newPokemon == player.Pokemon {
		return logs
	}
	oldName := player.Pokemon.Pokemon.Name
	player.Pokemon = newPokemon
	player.ActiveIndex = player.Action.SwitchIndex
	logs = append(logs, "🔄 "+player.Username+" 收回了 "+oldName+"，派出了 "+newPokemon.Pokemon.Name+"！")

	// 触发出场特性
	var opponent *Battler
	if player == b.Player1 && b.Player2 != nil {
		opponent = b.Player2.Pokemon
	} else if player == b.Player2 && b.Player1 != nil {
		opponent = b.Player1.Pokemon
	}
	if opponent != nil {
		entryLogs := b.TriggerEntryAbility(newPokemon, opponent)
		logs = append(logs, entryLogs...)
	}

	return logs
}

// executeAction 执行单个行动
func (b *Battle) executeAction(attacker, defender *BattlePlayer) []string {
	logs := make([]string, 0)

	// 检查是否需要充能（如破坏光线后的回合）
	if attacker.Pokemon.MustRecharge {
		logs = append(logs, "⏳ "+attacker.Pokemon.Pokemon.Name+" 正在充能，无法行动！")
		attacker.Pokemon.MustRecharge = false
		return logs
	}

	if attacker.Action.Type != ActionMove {
		return logs
	}

	move := attacker.Pokemon.Moves[attacker.Action.MoveIndex]
	move.Use()

	logs = append(logs, "▶️ "+attacker.Pokemon.Pokemon.Name+" 使用了 **"+move.Name+"**！")

	result := attacker.Pokemon.CalculateDamage(move, defender.Pokemon)

	if !result.Hit {
		logs = append(logs, "❌ 但是没有命中！")
		return logs
	}

	if move.Category == CategoryStatus {
		logs = append(logs, "✨ 效果发动了！")
		return logs
	}

	// 应用特性伤害修正
	if b.AbilityService != nil {
		ctx := b.GetBattleContext()
		moveAdapter := NewMoveAdapter(move)
		_, _, _, _, _, immune, abilityMsgs := b.AbilityService.CalculateDamageWithAbilities(
			attacker.Pokemon, defender.Pokemon, moveAdapter, ctx)
		logs = append(logs, abilityMsgs...)
		if immune {
			return logs
		}
	}

	if result.Critical {
		logs = append(logs, "💥 会心一击！")
	}

	defender.Pokemon.TakeDamageWithItem(result.Damage)

	// 属性克制提示
	if result.Effectiveness > 1 {
		logs = append(logs, "💥 效果拔群！")
	} else if result.Effectiveness < 1 && result.Effectiveness > 0 {
		logs = append(logs, "🛡️ 效果不佳...")
	} else if result.Effectiveness == 0 {
		logs = append(logs, "⚫ 没有效果...")
		return logs
	}

	logs = append(logs, "💔 造成了 **"+itoa(result.Damage)+"** 点伤害！")
	logs = append(logs, "❤️ "+defender.Pokemon.Pokemon.Name+" HP: "+itoa(defender.Pokemon.CurrentHP)+"/"+itoa(defender.Pokemon.MaxHP))

	// 触发受击特���（如静电、粗糙皮肤等）
	if b.AbilityService != nil && defender.Pokemon.IsAlive() {
		ctx := b.GetBattleContext()
		moveAdapter := NewMoveAdapter(move)
		hitResult := b.AbilityService.TriggerBeingHit(defender.Pokemon, attacker.Pokemon, moveAdapter, result.Damage, ctx)
		if hitResult != nil {
			logs = append(logs, hitResult.Messages...)
			// 处理接触效果（如麻痹、中毒）
			if hitResult.ContactEffect != "" && hitResult.ContactChance > 0 {
				if randInt(100) < hitResult.ContactChance {
					if attacker.Pokemon.GetStatus() == "" {
						attacker.Pokemon.SetStatus(hitResult.ContactEffect)
						logs = append(logs, "⚡ "+attacker.Pokemon.Pokemon.Name+" 陷入了"+hitResult.ContactEffect+"状态！")
					}
				}
			}
			// 处理反伤（如粗糙皮肤、铁刺）
			if hitResult.RecoilDamage > 0 {
				attacker.Pokemon.TakeDamage(hitResult.RecoilDamage)
				logs = append(logs, "💥 "+attacker.Pokemon.Pokemon.Name+" 受到了反伤！")
			}
			// 处理能力变化（如黏滑降速）
			if hitResult.StatChanges != nil {
				for stat, stages := range hitResult.StatChanges {
					if newStage, changed := attacker.Pokemon.ModifyStat(stat, stages); changed {
						if stages < 0 {
							logs = append(logs, "📉 "+attacker.Pokemon.Pokemon.Name+" 的"+getStatName(stat)+"下降了！(现在: "+itoa(newStage)+"级)")
						}
					}
				}
			}
		}
	}

	// 检查击倒触发特性（如自信过剩、异兽提升）
	if b.AbilityService != nil && !defender.Pokemon.IsAlive() {
		ctx := b.GetBattleContext()
		koResult := b.AbilityService.TriggerKO(attacker.Pokemon, defender.Pokemon, ctx)
		if koResult != nil {
			logs = append(logs, koResult.Messages...)
			if koResult.StatBoosts != nil {
				for stat, stages := range koResult.StatBoosts {
					if newStage, changed := attacker.Pokemon.ModifyStat(stat, stages); changed {
						logs = append(logs, "📈 "+attacker.Pokemon.Pokemon.Name+" 的"+getStatName(stat)+"提升了！(现在: "+itoa(newStage)+"级)")
					}
				}
			}
		}
	}

	// 检查技能是否需要充能（如破坏光线）
	if move.RechargeRequired {
		attacker.Pokemon.MustRecharge = true
	}

	return logs
}

// clearActions 清除行动
func (b *Battle) clearActions() {
	if b.Player1 != nil {
		b.Player1.Action = nil
	}
	if b.Player2 != nil {
		b.Player2.Action = nil
	}
}

// IsPlayerTurn 检查是否轮到该玩家
func (b *Battle) IsPlayerTurn(playerID string) bool {
	player := b.GetPlayer(playerID)
	return player != nil && player.Action == nil
}

// GetBattleStatus 获取对战状态描述
func (b *Battle) GetBattleStatus() string {
	if b.State == BattleStateWaiting {
		return "等待对手加入..."
	}
	if b.State == BattleStateChoosing {
		status := "选择宝可梦阶段\n"
		if b.Player1 != nil {
			if b.Player1.Ready {
				status += "✅ " + b.Player1.Username + " 已准备\n"
			} else {
				status += "⏳ " + b.Player1.Username + " 选择中...\n"
			}
		}
		if b.Player2 != nil {
			if b.Player2.Ready {
				status += "✅ " + b.Player2.Username + " 已准备"
			} else {
				status += "⏳ " + b.Player2.Username + " 选择中..."
			}
		}
		return status
	}
	if b.State == BattleStateFinished {
		return "对战已结束"
	}
	return "对战进行中"
}

// ============================================
// 特性系统集成
// ============================================

// GetBattleContext 获取战斗上下文（用于特性系统）
func (b *Battle) GetBattleContext() *ability.BattleContext {
	return &ability.BattleContext{
		Weather:   b.Weather,
		Terrain:   b.Terrain,
		Turn:      b.CurrentTurn,
		IsDoubles: false,
	}
}

// TriggerEntryAbility 触发出场特性
func (b *Battle) TriggerEntryAbility(self *Battler, opponent *Battler) []string {
	logs := make([]string, 0)
	if b.AbilityService == nil || self.Ability == nil {
		return logs
	}

	ctx := b.GetBattleContext()
	messages, weather, statChanges := b.AbilityService.ProcessEntryAbility(self, opponent, ctx)

	logs = append(logs, messages...)

	// 设置天气
	if weather != nil {
		b.Weather = *weather
		b.WeatherTurns = 5
	}

	// 应用对手能力变化
	if statChanges != nil {
		for stat, stages := range statChanges {
			if newStage, changed := opponent.ModifyStat(stat, stages); changed {
				if stages < 0 {
					logs = append(logs, "📉 "+opponent.Pokemon.Name+" 的"+getStatName(stat)+"下降了！(现在: "+itoa(newStage)+"级)")
				} else {
					logs = append(logs, "📈 "+opponent.Pokemon.Name+" 的"+getStatName(stat)+"提升了！(现在: "+itoa(newStage)+"级)")
				}
			}
		}
	}

	return logs
}

// TriggerTurnEndAbilities 触发回合结束特性
func (b *Battle) TriggerTurnEndAbilities() []string {
	logs := make([]string, 0)
	if b.AbilityService == nil {
		return logs
	}

	ctx := b.GetBattleContext()

	// 处理天气伤害/回复
	if b.Weather != valueobject.WeatherNone {
		logs = append(logs, b.processWeatherEffects()...)
		b.WeatherTurns--
		if b.WeatherTurns <= 0 {
			logs = append(logs, "☀️ 天气恢复正常了。")
			b.Weather = valueobject.WeatherNone
		}
	}

	// 玩家1的回合结束特性
	if b.Player1 != nil && b.Player1.Pokemon != nil && b.Player1.Pokemon.IsAlive() {
		messages, statBoosts, healing, damage := b.AbilityService.ProcessTurnEndAbility(b.Player1.Pokemon, ctx)
		logs = append(logs, messages...)
		if healing > 0 {
			b.Player1.Pokemon.Heal(healing)
		}
		if damage > 0 {
			b.Player1.Pokemon.TakeDamage(damage)
		}
		if statBoosts != nil {
			for stat, stages := range statBoosts {
				b.Player1.Pokemon.ModifyStat(stat, stages)
			}
		}
	}

	// 玩家2的回合结束特性
	if b.Player2 != nil && b.Player2.Pokemon != nil && b.Player2.Pokemon.IsAlive() {
		messages, statBoosts, healing, damage := b.AbilityService.ProcessTurnEndAbility(b.Player2.Pokemon, ctx)
		logs = append(logs, messages...)
		if healing > 0 {
			b.Player2.Pokemon.Heal(healing)
		}
		if damage > 0 {
			b.Player2.Pokemon.TakeDamage(damage)
		}
		if statBoosts != nil {
			for stat, stages := range statBoosts {
				b.Player2.Pokemon.ModifyStat(stat, stages)
			}
		}
	}

	return logs
}

// processWeatherEffects 处理天气效果
func (b *Battle) processWeatherEffects() []string {
	logs := make([]string, 0)

	processPokemon := func(pokemon *Battler) {
		if pokemon == nil || !pokemon.IsAlive() {
			return
		}

		// 检查是否免疫天气伤害
		immune := false
		for _, t := range pokemon.Types {
			switch b.Weather {
			case valueobject.WeatherSand:
				if t == valueobject.TypeRock || t == valueobject.TypeGround || t == valueobject.TypeSteel {
					immune = true
				}
			case valueobject.WeatherHail:
				if t == valueobject.TypeIce {
					immune = true
				}
			}
		}

		if !immune {
			switch b.Weather {
			case valueobject.WeatherSand:
				damage := pokemon.MaxHP / 16
				if damage < 1 {
					damage = 1
				}
				pokemon.TakeDamage(damage)
				logs = append(logs, "🏜️ "+pokemon.Pokemon.Name+" 受到了沙暴伤害！")
			case valueobject.WeatherHail:
				damage := pokemon.MaxHP / 16
				if damage < 1 {
					damage = 1
				}
				pokemon.TakeDamage(damage)
				logs = append(logs, "🌨️ "+pokemon.Pokemon.Name+" 受到了冰雹伤害！")
			}
		}
	}

	if b.Player1 != nil {
		processPokemon(b.Player1.Pokemon)
	}
	if b.Player2 != nil {
		processPokemon(b.Player2.Pokemon)
	}

	return logs
}

// getStatName 获取能力名称
func getStatName(stat string) string {
	names := map[string]string{
		"attack":    "攻击",
		"defense":   "防御",
		"spattack":  "特攻",
		"spdefense": "特防",
		"speed":     "速度",
		"accuracy":  "命中",
		"evasion":   "闪避",
	}
	if name, ok := names[stat]; ok {
		return name
	}
	return stat
}

// randInt 生成 0 到 max-1 的随机整数
func randInt(max int) int {
	if max <= 0 {
		return 0
	}
	return int(time.Now().UnixNano() % int64(max))
}
