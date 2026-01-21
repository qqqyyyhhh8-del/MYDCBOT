package pokemon

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/user/dcminigames/internal/domain/pokemon/entity"
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
	"github.com/user/dcminigames/internal/infrastructure/persistence/memory"
	"github.com/user/dcminigames/internal/infrastructure/pokeapi"
)

// PokemonConfig 玩家配置中的宝可梦设置
type PokemonConfig struct {
	PokemonID   int
	Nature      valueobject.Nature
	AbilitySlot int                   // 0=第一特性, 1=第二特性, -1=隐藏特性
	MoveIndices []int                 // 选择的技能索引
	TeraType    valueobject.PokeType  // 太晶属性
}

// TeamPreset 配队预设
type TeamPreset struct {
	ID          string           // 预设ID
	Name        string           // 预设名称
	UserID      string           // 用户ID
	PokemonID   int              // 宝可梦ID
	PokemonName string           // 宝可梦名称（显示用）
	Nature      valueobject.Nature
	AbilitySlot int
	MoveIndices []int
	TeraType    valueobject.PokeType
}

// Handler 宝可梦对战应用层处理器
type Handler struct {
	repo        *memory.BattleRepository
	client      *pokeapi.Client
	configMu    sync.RWMutex
	configs     map[string]*PokemonConfig // key: channelID:playerID
	presetMu    sync.RWMutex
	presets     map[string][]*TeamPreset  // key: userID
}

// NewHandler 创建处理器
func NewHandler(repo *memory.BattleRepository) *Handler {
	return &Handler{
		repo:    repo,
		client:  pokeapi.NewClient(),
		configs: make(map[string]*PokemonConfig),
		presets: make(map[string][]*TeamPreset),
	}
}

// configKey 生成配置键
func configKey(channelID, playerID string) string {
	return channelID + ":" + playerID
}

// SetConfig 设置玩家配置
func (h *Handler) SetConfig(channelID, playerID string, config *PokemonConfig) {
	h.configMu.Lock()
	defer h.configMu.Unlock()
	h.configs[configKey(channelID, playerID)] = config
}

// GetConfig 获取玩家配置
func (h *Handler) GetConfig(channelID, playerID string) *PokemonConfig {
	h.configMu.RLock()
	defer h.configMu.RUnlock()
	return h.configs[configKey(channelID, playerID)]
}

// ClearConfig 清除玩家配置
func (h *Handler) ClearConfig(channelID, playerID string) {
	h.configMu.Lock()
	defer h.configMu.Unlock()
	delete(h.configs, configKey(channelID, playerID))
}

// CreateBattle 创建对战 (默认1v1)
func (h *Handler) CreateBattle(channelID, playerID, username string) (*entity.Battle, error) {
	return h.CreateBattleWithTeamSize(channelID, playerID, username, entity.TeamSize1v1)
}

// CreateBattleWithTeamSize 创建指定队伍大小的对战
func (h *Handler) CreateBattleWithTeamSize(channelID, playerID, username string, teamSize entity.TeamSize) (*entity.Battle, error) {
	if h.repo.Exists(channelID) {
		return nil, fmt.Errorf("该频道已有对战进行中")
	}
	battle := entity.NewBattleWithTeamSize(uuid.New().String(), channelID, teamSize)
	if err := battle.AddPlayer(playerID, username); err != nil {
		return nil, err
	}
	if err := h.repo.Save(battle); err != nil {
		return nil, err
	}
	return battle, nil
}

// CreateAIBattle 创建人机对战
func (h *Handler) CreateAIBattle(channelID, playerID, username string, teamSize entity.TeamSize) (*entity.Battle, error) {
	if h.repo.Exists(channelID) {
		return nil, fmt.Errorf("该频道已有对战进行中")
	}
	battle := entity.NewAIBattle(uuid.New().String(), channelID, teamSize)
	// 玩家加入
	if err := battle.AddPlayer(playerID, username); err != nil {
		return nil, err
	}
	// AI 加入
	if err := battle.AddPlayer(entity.AIPlayerID, entity.AIPlayerName); err != nil {
		return nil, err
	}
	// AI 自动选择宝可梦
	if err := h.aiSelectPokemon(battle, teamSize); err != nil {
		return nil, err
	}
	if err := h.repo.Save(battle); err != nil {
		return nil, err
	}
	return battle, nil
}

// aiSelectPokemon AI 自动选择宝可梦
func (h *Handler) aiSelectPokemon(battle *entity.Battle, teamSize entity.TeamSize) error {
	// 热门宝可梦 ID 列表（用于 AI 选择）
	popularPokemonIDs := []int{
		6,   // 喷火龙
		9,   // 水箭龟
		3,   // 妙蛙花
		25,  // 皮卡丘
		150, // 超梦
		149, // 快龙
		130, // 暴鲤龙
		143, // 卡比兽
		94,  // 耿鬼
		65,  // 胡地
		131, // 乘龙
		59,  // 风速狗
		38,  // 九尾
		68,  // 怪力
		76,  // 隆隆岩
	}

	// 随机打乱
	rand.Shuffle(len(popularPokemonIDs), func(i, j int) {
		popularPokemonIDs[i], popularPokemonIDs[j] = popularPokemonIDs[j], popularPokemonIDs[i]
	})

	count := int(teamSize)
	selected := 0

	for _, pokemonID := range popularPokemonIDs {
		if selected >= count {
			break
		}
		pokemon := pokeapi.GetPredefinedPokemon(pokemonID)
		if pokemon == nil {
			continue
		}
		// 随机选择性格
		natures := []valueobject.Nature{
			valueobject.NatureAdamant,
			valueobject.NatureJolly,
			valueobject.NatureModest,
			valueobject.NatureTimid,
			valueobject.NatureBold,
		}
		pokemon.Nature = natures[rand.Intn(len(natures))]

		// 选择前4个技能
		if len(pokemon.LearnableMoves) > 4 {
			pokemon.LearnableMoves = pokemon.LearnableMoves[:4]
		}

		if err := battle.SetPokemon(entity.AIPlayerID, pokemon, 50); err != nil {
			continue
		}
		selected++
	}

	if selected == 0 {
		return fmt.Errorf("AI 无法选择宝可梦")
	}
	return nil
}

// AIChooseAction AI 选择行动
func (h *Handler) AIChooseAction(battle *entity.Battle) *entity.BattleAction {
	aiPlayer := battle.GetAIPlayer()
	if aiPlayer == nil || aiPlayer.Pokemon == nil {
		return nil
	}

	humanPlayer := battle.GetHumanPlayer()
	if humanPlayer == nil || humanPlayer.Pokemon == nil {
		return nil
	}

	// 简单 AI 策略：选择伤害最高的技能
	bestMoveIdx := 0
	bestScore := 0.0

	for idx, move := range aiPlayer.Pokemon.Moves {
		if !move.CanUse() {
			continue
		}

		// 计算预估伤害分数
		score := float64(move.Power)

		// 属性克制加成
		effectiveness := valueobject.GetEffectiveness(move.Type, humanPlayer.Pokemon.Pokemon.Types)
		score *= effectiveness

		// STAB 加成
		for _, atkType := range aiPlayer.Pokemon.Pokemon.Types {
			if atkType == move.Type {
				score *= 1.5
				break
			}
		}

		// 添加随机因素避免太机械
		score *= (0.9 + rand.Float64()*0.2)

		if score > bestScore {
			bestScore = score
			bestMoveIdx = idx
		}
	}

	return &entity.BattleAction{
		Type:      entity.ActionMove,
		MoveIndex: bestMoveIdx,
	}
}

// ExecuteAITurn 执行 AI 回合（玩家行动后自动触发）
func (h *Handler) ExecuteAITurn(channelID string) ([]string, error) {
	battle, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}

	if !battle.IsAIBattle {
		return nil, nil
	}

	aiPlayer := battle.GetAIPlayer()
	if aiPlayer == nil {
		return nil, nil
	}

	// AI 选择行动
	if aiPlayer.Action == nil {
		aiAction := h.AIChooseAction(battle)
		if aiAction != nil {
			aiPlayer.Action = aiAction
		}
	}

	// 执行回合
	var logs []string
	if battle.BothActionsReady() {
		logs = battle.ExecuteTurn()
	}

	if err := h.repo.Save(battle); err != nil {
		return nil, err
	}

	return logs, nil
}

// JoinBattle 加入对战
func (h *Handler) JoinBattle(channelID, playerID, username string) error {
	battle, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return fmt.Errorf("没有进行中的对战")
	}
	if err := battle.AddPlayer(playerID, username); err != nil {
		return err
	}
	return h.repo.Save(battle)
}

// GetBattle 获取对战
func (h *Handler) GetBattle(channelID string) (*entity.Battle, error) {
	return h.repo.FindByChannelID(channelID)
}

// SelectPokemon 选择宝可梦（使用配置）
func (h *Handler) SelectPokemon(channelID, playerID string, pokemonID, level int) error {
	battle, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return err
	}

	pokemon := pokeapi.GetPredefinedPokemon(pokemonID)
	if pokemon == nil {
		return fmt.Errorf("未找到宝可梦")
	}

	// 应用玩家配置
	config := h.GetConfig(channelID, playerID)
	if config != nil {
		// 应用性格
		if config.Nature != "" {
			pokemon.Nature = config.Nature
		}
		// 应用特性
		if config.AbilitySlot >= 0 && config.AbilitySlot < len(pokemon.Abilities) {
			pokemon.SelectedAbility = &pokemon.Abilities[config.AbilitySlot]
		}
		// 应用技能选择
		if len(config.MoveIndices) > 0 {
			var selectedMoves []*entity.Move
			for _, idx := range config.MoveIndices {
				if idx >= 0 && idx < len(pokemon.LearnableMoves) {
					selectedMoves = append(selectedMoves, pokemon.LearnableMoves[idx])
				}
			}
			if len(selectedMoves) > 0 {
				pokemon.LearnableMoves = selectedMoves
			}
		}
		// 应用太晶属性
		if config.TeraType != "" {
			pokemon.TeraType = config.TeraType
		}
		// 清除配置
		h.ClearConfig(channelID, playerID)
	}

	if err := battle.SetPokemon(playerID, pokemon, level); err != nil {
		return err
	}

	return h.repo.Save(battle)
}

// UseMove 使用技能
func (h *Handler) UseMove(channelID, playerID string, moveIndex int) ([]string, error) {
	battle, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}

	action := &entity.BattleAction{
		Type:      entity.ActionMove,
		MoveIndex: moveIndex,
	}

	if err := battle.SetAction(playerID, action); err != nil {
		return nil, err
	}

	// 检查是否双方都已行动
	var logs []string
	if battle.BothActionsReady() {
		logs = battle.ExecuteTurn()
	}

	if err := h.repo.Save(battle); err != nil {
		return nil, err
	}

	return logs, nil
}

// Forfeit 认输
func (h *Handler) Forfeit(channelID, playerID string) ([]string, error) {
	battle, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}

	action := &entity.BattleAction{
		Type: entity.ActionForfeit,
	}

	if err := battle.SetAction(playerID, action); err != nil {
		return nil, err
	}

	// 立即执行（对手自动跳过）
	opponent := battle.GetOpponent(playerID)
	if opponent != nil && opponent.Action == nil {
		opponent.Action = &entity.BattleAction{Type: entity.ActionMove, MoveIndex: 0}
	}

	logs := battle.ExecuteTurn()

	if err := h.repo.Save(battle); err != nil {
		return nil, err
	}

	return logs, nil
}

// SwitchPokemon 换人
func (h *Handler) SwitchPokemon(channelID, playerID string, switchIndex int) ([]string, error) {
	battle, err := h.repo.FindByChannelID(channelID)
	if err != nil {
		return nil, err
	}

	player := battle.GetPlayer(playerID)
	if player == nil {
		return nil, fmt.Errorf("你不在对战中")
	}

	// 验证换人目标
	if switchIndex < 0 || switchIndex >= len(player.Team) {
		return nil, fmt.Errorf("无效的宝可梦")
	}
	target := player.Team[switchIndex]
	if !target.IsAlive() {
		return nil, fmt.Errorf("该宝可梦已倒下")
	}
	if target == player.Pokemon {
		return nil, fmt.Errorf("该宝可梦已在场上")
	}

	action := &entity.BattleAction{
		Type:        entity.ActionSwitch,
		SwitchIndex: switchIndex,
	}

	if err := battle.SetAction(playerID, action); err != nil {
		return nil, err
	}

	// 检查是否双方都已行动
	var logs []string
	if battle.BothActionsReady() {
		logs = battle.ExecuteTurn()
	}

	if err := h.repo.Save(battle); err != nil {
		return nil, err
	}

	return logs, nil
}

// EndBattle 结束对战
func (h *Handler) EndBattle(channelID string) error {
	return h.repo.Delete(channelID)
}

// GetAvailablePokemon 获取可选宝可梦列表
func (h *Handler) GetAvailablePokemon() []*entity.Pokemon {
	return pokeapi.GetAllPredefinedPokemon()
}

// GetPokemonByID 通过ID获取宝可梦
func (h *Handler) GetPokemonByID(id int) *entity.Pokemon {
	return pokeapi.GetPredefinedPokemon(id)
}

// GetSpriteURL 获取精灵图URL
func (h *Handler) GetSpriteURL(pokemonID int) string {
	return pokeapi.GetSpriteURL(pokemonID)
}

// IsPlayerInBattle 检查玩家是否在对战中
func (h *Handler) IsPlayerInBattle(playerID string) bool {
	_, err := h.repo.FindByPlayerID(playerID)
	return err == nil
}

// SavePreset 保存配队预设
func (h *Handler) SavePreset(userID, name string, config *PokemonConfig) (*TeamPreset, error) {
	pokemon := pokeapi.GetPredefinedPokemon(config.PokemonID)
	if pokemon == nil {
		return nil, fmt.Errorf("未找到宝可梦")
	}

	preset := &TeamPreset{
		ID:          uuid.New().String()[:8],
		Name:        name,
		UserID:      userID,
		PokemonID:   config.PokemonID,
		PokemonName: pokemon.Name,
		Nature:      config.Nature,
		AbilitySlot: config.AbilitySlot,
		MoveIndices: append([]int{}, config.MoveIndices...),
		TeraType:    config.TeraType,
	}

	h.presetMu.Lock()
	defer h.presetMu.Unlock()

	// 检查是否超过最大数量（最多10个预设）
	if len(h.presets[userID]) >= 10 {
		return nil, fmt.Errorf("预设数量已达上限（10个）")
	}

	h.presets[userID] = append(h.presets[userID], preset)
	return preset, nil
}

// GetPresets 获取用户所有预设
func (h *Handler) GetPresets(userID string) []*TeamPreset {
	h.presetMu.RLock()
	defer h.presetMu.RUnlock()
	return h.presets[userID]
}

// GetPreset 获取指定预设
func (h *Handler) GetPreset(userID, presetID string) *TeamPreset {
	h.presetMu.RLock()
	defer h.presetMu.RUnlock()
	for _, p := range h.presets[userID] {
		if p.ID == presetID {
			return p
		}
	}
	return nil
}

// DeletePreset 删除预设
func (h *Handler) DeletePreset(userID, presetID string) error {
	h.presetMu.Lock()
	defer h.presetMu.Unlock()

	presets := h.presets[userID]
	for i, p := range presets {
		if p.ID == presetID {
			h.presets[userID] = append(presets[:i], presets[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("预设不存在")
}

// LoadPresetToConfig 加载预设到当前配置
func (h *Handler) LoadPresetToConfig(channelID, userID, presetID string) error {
	preset := h.GetPreset(userID, presetID)
	if preset == nil {
		return fmt.Errorf("预设不存在")
	}

	config := &PokemonConfig{
		PokemonID:   preset.PokemonID,
		Nature:      preset.Nature,
		AbilitySlot: preset.AbilitySlot,
		MoveIndices: append([]int{}, preset.MoveIndices...),
		TeraType:    preset.TeraType,
	}
	h.SetConfig(channelID, userID, config)
	return nil
}
