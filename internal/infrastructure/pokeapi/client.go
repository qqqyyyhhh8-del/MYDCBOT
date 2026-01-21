package pokeapi

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/user/dcminigames/internal/domain/pokemon/entity"
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
)

const (
	// GitHub Raw CSV 基础 URL
	githubCSVBase = "https://raw.githubusercontent.com/PokeAPI/pokeapi/master/data/v2/csv"
	// 简体中文语言ID
	langZhHans = 12
	// 最大宝可梦ID（第9世代）
	maxPokemonID = 1025
)

// Client PokeAPI 客户端（使用 GitHub CSV 数据）
type Client struct {
	httpClient *http.Client
	cache      *DataCache
	loading    bool
	loadMu     sync.Mutex
}

// PokemonAbilityInfo 宝可梦特性信息
type PokemonAbilityInfo struct {
	AbilityID int
	IsHidden  bool
	Slot      int
}

// DataCache 数据缓存
type DataCache struct {
	mu sync.RWMutex
	// 宝可梦名称 map[id]name
	PokemonNames map[int]string
	// 宝可梦种类 map[id]genus
	PokemonGenus map[int]string
	// 宝可梦基础数据 map[id]*Pokemon
	Pokemon map[int]*entity.Pokemon
	// 技能名称 map[id]name
	MoveNames map[int]string
	// 技能数据 map[id]*Move
	Moves map[int]*entity.Move
	// 宝可梦可学技能 map[pokemonID][]moveID
	PokemonMoves map[int][]int
	// 属性名称 map[id]name
	TypeNames map[int]string
	// 特性名称 map[id]name
	AbilityNames map[int]string
	// 特性描述 map[id]description
	AbilityDescriptions map[int]string
	// 宝可梦特性 map[pokemonID][]abilityID (旧版兼容)
	PokemonAbilities map[int][]int
	// 宝可梦特性详情 map[pokemonID][]PokemonAbilityInfo
	PokemonAbilityInfos map[int][]PokemonAbilityInfo
	// 是否已加载
	Loaded bool
}

// NewClient 创建客户端
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		cache: &DataCache{
			PokemonNames:        make(map[int]string),
			PokemonGenus:        make(map[int]string),
			Pokemon:             make(map[int]*entity.Pokemon),
			MoveNames:           make(map[int]string),
			Moves:               make(map[int]*entity.Move),
			PokemonMoves:        make(map[int][]int),
			TypeNames:           make(map[int]string),
			AbilityNames:        make(map[int]string),
			AbilityDescriptions: make(map[int]string),
			PokemonAbilities:    make(map[int][]int),
			PokemonAbilityInfos: make(map[int][]PokemonAbilityInfo),
		},
	}
}

// 默认客户端实例
var defaultClient = NewClient()

// fetchCSV 获取 CSV 数据
func (c *Client) fetchCSV(ctx context.Context, filename string) ([][]string, error) {
	url := githubCSVBase + "/" + filename
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "DcMiniGames/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}

	// 跳过标题行
	if len(records) > 0 {
		records = records[1:]
	}
	return records, nil
}

// LoadAllData 加载所有数据
func (c *Client) LoadAllData(ctx context.Context) error {
	c.loadMu.Lock()
	if c.loading {
		c.loadMu.Unlock()
		return nil
	}
	c.loading = true
	c.loadMu.Unlock()

	defer func() {
		c.loadMu.Lock()
		c.loading = false
		c.loadMu.Unlock()
	}()

	c.cache.mu.Lock()
	if c.cache.Loaded {
		c.cache.mu.Unlock()
		return nil
	}
	c.cache.mu.Unlock()

	// 并行加载数据
	var wg sync.WaitGroup
	errChan := make(chan error, 6)

	// 加载属性名称
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.loadTypeNames(ctx); err != nil {
			errChan <- fmt.Errorf("加载属性名称失败: %w", err)
		}
	}()

	// 加载宝可梦名称
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.loadPokemonNames(ctx); err != nil {
			errChan <- fmt.Errorf("加载宝可梦名称失败: %w", err)
		}
	}()

	// 加载技能名称
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.loadMoveNames(ctx); err != nil {
			errChan <- fmt.Errorf("加载技能名称失败: %w", err)
		}
	}()

	// 加��技能数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.loadMoves(ctx); err != nil {
			errChan <- fmt.Errorf("加载技能数据失败: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	// 检查错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// 加载宝可梦基础数据（需要属性名称）
	if err := c.loadPokemonData(ctx); err != nil {
		return fmt.Errorf("加载宝可梦数据失败: %w", err)
	}

	// 加载宝可梦可学技能
	if err := c.loadPokemonMoves(ctx); err != nil {
		return fmt.Errorf("加载宝可梦技能失败: %w", err)
	}

	// 加载特性名称
	if err := c.loadAbilityNames(ctx); err != nil {
		return fmt.Errorf("加载特性名称失败: %w", err)
	}

	// 加载特性描述
	if err := c.loadAbilityDescriptions(ctx); err != nil {
		return fmt.Errorf("加载特性描述失败: %w", err)
	}

	// 加载宝可梦特性
	if err := c.loadPokemonAbilities(ctx); err != nil {
		return fmt.Errorf("加载宝可梦特性失败: %w", err)
	}

	c.cache.mu.Lock()
	c.cache.Loaded = true
	c.cache.mu.Unlock()

	return nil
}

// loadTypeNames 加载属性名称
func (c *Client) loadTypeNames(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "type_names.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		typeID, _ := strconv.Atoi(record[0])
		langID, _ := strconv.Atoi(record[1])
		name := record[2]

		if langID == langZhHans && typeID > 0 {
			c.cache.TypeNames[typeID] = name
		}
	}
	return nil
}

// loadPokemonNames 加载宝可梦名称
func (c *Client) loadPokemonNames(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "pokemon_species_names.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		pokemonID, _ := strconv.Atoi(record[0])
		langID, _ := strconv.Atoi(record[1])
		name := record[2]
		genus := record[3]

		if langID == langZhHans && pokemonID > 0 && pokemonID <= maxPokemonID {
			c.cache.PokemonNames[pokemonID] = name
			c.cache.PokemonGenus[pokemonID] = genus
		}
	}
	return nil
}

// loadMoveNames 加载技能名称
func (c *Client) loadMoveNames(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "move_names.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		moveID, _ := strconv.Atoi(record[0])
		langID, _ := strconv.Atoi(record[1])
		name := record[2]

		if langID == langZhHans && moveID > 0 {
			c.cache.MoveNames[moveID] = name
		}
	}
	return nil
}

// loadMoves 加载技能数据
func (c *Client) loadMoves(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "moves.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	// CSV格式: id,identifier,generation_id,type_id,power,pp,accuracy,priority,target_id,damage_class_id,...
	for _, record := range records {
		if len(record) < 10 {
			continue
		}
		moveID, _ := strconv.Atoi(record[0])
		typeID, _ := strconv.Atoi(record[3])
		power, _ := strconv.Atoi(record[4])
		pp, _ := strconv.Atoi(record[5])
		accuracy, _ := strconv.Atoi(record[6])
		priority, _ := strconv.Atoi(record[7])
		damageClass, _ := strconv.Atoi(record[9])

		if moveID <= 0 {
			continue
		}

		// 获取属性
		pokeType := typeIDToPokeType(typeID)

		// 获取分类
		var category entity.MoveCategory
		switch damageClass {
		case 2:
			category = entity.CategoryPhysical
		case 3:
			category = entity.CategorySpecial
		default:
			category = entity.CategoryStatus
		}

		// 判断是否需要充能（破坏光线等技能）
		rechargeRequired := isRechargeMove(moveID)

		c.cache.Moves[moveID] = &entity.Move{
			Type:             pokeType,
			Category:         category,
			Power:            power,
			Accuracy:         accuracy,
			PP:               pp,
			MaxPP:            pp,
			Priority:         priority,
			RechargeRequired: rechargeRequired,
		}
	}
	return nil
}

// loadPokemonData 加载宝可梦数据
func (c *Client) loadPokemonData(ctx context.Context) error {
	// 加载基础数据
	pokemonRecords, err := c.fetchCSV(ctx, "pokemon.csv")
	if err != nil {
		return err
	}

	// 加载种族值
	statsRecords, err := c.fetchCSV(ctx, "pokemon_stats.csv")
	if err != nil {
		return err
	}

	// 加载属性
	typesRecords, err := c.fetchCSV(ctx, "pokemon_types.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	// 解析基础数据
	for _, record := range pokemonRecords {
		if len(record) < 8 {
			continue
		}
		id, _ := strconv.Atoi(record[0])
		isDefault := record[7] == "1"

		// 只处理默认形态，且在范围内
		if id <= 0 || id > maxPokemonID || !isDefault {
			continue
		}

		name := c.cache.PokemonNames[id]
		if name == "" {
			name = record[1] // 使用英文名
		}

		c.cache.Pokemon[id] = &entity.Pokemon{
			ID:             id,
			Name:           name,
			Types:          []valueobject.PokeType{},
			LearnableMoves: []*entity.Move{},
			SpriteURL:      GetSpriteURL(id),
		}
	}

	// 解析种族值
	// CSV格式: pokemon_id,stat_id,base_stat,effort
	// stat_id: 1=HP, 2=攻击, 3=防御, 4=特攻, 5=特防, 6=速度
	for _, record := range statsRecords {
		if len(record) < 3 {
			continue
		}
		pokemonID, _ := strconv.Atoi(record[0])
		statID, _ := strconv.Atoi(record[1])
		baseStat, _ := strconv.Atoi(record[2])

		pokemon := c.cache.Pokemon[pokemonID]
		if pokemon == nil {
			continue
		}

		switch statID {
		case 1:
			pokemon.BaseHP = baseStat
		case 2:
			pokemon.BaseAtk = baseStat
		case 3:
			pokemon.BaseDef = baseStat
		case 4:
			pokemon.BaseSpAtk = baseStat
		case 5:
			pokemon.BaseSpDef = baseStat
		case 6:
			pokemon.BaseSpeed = baseStat
		}
	}

	// 解析属性
	// CSV格式: pokemon_id,type_id,slot
	for _, record := range typesRecords {
		if len(record) < 3 {
			continue
		}
		pokemonID, _ := strconv.Atoi(record[0])
		typeID, _ := strconv.Atoi(record[1])

		pokemon := c.cache.Pokemon[pokemonID]
		if pokemon == nil {
			continue
		}

		pokeType := typeIDToPokeType(typeID)
		if pokeType != "" {
			pokemon.Types = append(pokemon.Types, pokeType)
		}
	}

	return nil
}

// loadPokemonMoves 加载宝可梦可学技能
func (c *Client) loadPokemonMoves(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "pokemon_moves.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	// CSV格式: pokemon_id,version_group_id,move_id,pokemon_move_method_id,level,order
	// 只取最新版本组的技能（version_group_id 较大的）
	latestVersion := make(map[int]int) // pokemon_id -> max version_group_id
	moveSet := make(map[int]map[int]bool) // pokemon_id -> set of move_ids

	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		pokemonID, _ := strconv.Atoi(record[0])
		versionGroupID, _ := strconv.Atoi(record[1])
		moveID, _ := strconv.Atoi(record[2])

		if pokemonID <= 0 || pokemonID > maxPokemonID || moveID <= 0 {
			continue
		}

		// 记录最新版本
		if versionGroupID > latestVersion[pokemonID] {
			latestVersion[pokemonID] = versionGroupID
		}

		// 添加技能
		if moveSet[pokemonID] == nil {
			moveSet[pokemonID] = make(map[int]bool)
		}
		moveSet[pokemonID][moveID] = true
	}

	// 转换为切片并排序
	for pokemonID, moves := range moveSet {
		var moveIDs []int
		for moveID := range moves {
			moveIDs = append(moveIDs, moveID)
		}
		sort.Ints(moveIDs)
		c.cache.PokemonMoves[pokemonID] = moveIDs
	}

	return nil
}

// loadAbilityNames 加载特性名称
func (c *Client) loadAbilityNames(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "ability_names.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	// CSV格式: ability_id,local_language_id,name
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		abilityID, _ := strconv.Atoi(record[0])
		langID, _ := strconv.Atoi(record[1])
		name := record[2]

		if langID == langZhHans && abilityID > 0 {
			c.cache.AbilityNames[abilityID] = name
		}
	}
	return nil
}

// loadAbilityDescriptions 加载特性描述
func (c *Client) loadAbilityDescriptions(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "ability_flavor_text.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	// CSV格式: ability_id,version_group_id,language_id,flavor_text
	// 取最新版本组的描述
	latestVersion := make(map[int]int) // ability_id -> max version_group_id

	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		abilityID, _ := strconv.Atoi(record[0])
		versionGroupID, _ := strconv.Atoi(record[1])
		langID, _ := strconv.Atoi(record[2])
		flavorText := record[3]

		if langID != langZhHans || abilityID <= 0 {
			continue
		}

		// 只保留最新版本的描述
		if versionGroupID >= latestVersion[abilityID] {
			latestVersion[abilityID] = versionGroupID
			// 清理换行符
			flavorText = strings.ReplaceAll(flavorText, "\n", "")
			c.cache.AbilityDescriptions[abilityID] = flavorText
		}
	}
	return nil
}

// loadPokemonAbilities 加载宝可梦特性
func (c *Client) loadPokemonAbilities(ctx context.Context) error {
	records, err := c.fetchCSV(ctx, "pokemon_abilities.csv")
	if err != nil {
		return err
	}

	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()

	// CSV格式: pokemon_id,ability_id,is_hidden,slot
	abilitySet := make(map[int]map[int]bool)           // pokemon_id -> set of ability_ids (旧版兼容)
	abilityInfos := make(map[int][]PokemonAbilityInfo) // pokemon_id -> []PokemonAbilityInfo

	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		pokemonID, _ := strconv.Atoi(record[0])
		abilityID, _ := strconv.Atoi(record[1])
		isHidden := record[2] == "1"
		slot, _ := strconv.Atoi(record[3])

		if pokemonID <= 0 || pokemonID > maxPokemonID || abilityID <= 0 {
			continue
		}

		// 旧版兼容
		if abilitySet[pokemonID] == nil {
			abilitySet[pokemonID] = make(map[int]bool)
		}
		abilitySet[pokemonID][abilityID] = true

		// 新版详情
		abilityInfos[pokemonID] = append(abilityInfos[pokemonID], PokemonAbilityInfo{
			AbilityID: abilityID,
			IsHidden:  isHidden,
			Slot:      slot,
		})
	}

	// 转换为切片并排序
	for pokemonID, abilities := range abilitySet {
		var abilityIDs []int
		for abilityID := range abilities {
			abilityIDs = append(abilityIDs, abilityID)
		}
		sort.Ints(abilityIDs)
		c.cache.PokemonAbilities[pokemonID] = abilityIDs
	}

	// 按 slot 排序特性详情
	for pokemonID, infos := range abilityInfos {
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Slot < infos[j].Slot
		})
		c.cache.PokemonAbilityInfos[pokemonID] = infos
	}

	return nil
}

// typeIDToPokeType 将类型ID转换为PokeType
func typeIDToPokeType(typeID int) valueobject.PokeType {
	switch typeID {
	case 1:
		return valueobject.TypeNormal
	case 2:
		return valueobject.TypeFighting
	case 3:
		return valueobject.TypeFlying
	case 4:
		return valueobject.TypePoison
	case 5:
		return valueobject.TypeGround
	case 6:
		return valueobject.TypeRock
	case 7:
		return valueobject.TypeBug
	case 8:
		return valueobject.TypeGhost
	case 9:
		return valueobject.TypeSteel
	case 10:
		return valueobject.TypeFire
	case 11:
		return valueobject.TypeWater
	case 12:
		return valueobject.TypeGrass
	case 13:
		return valueobject.TypeElectric
	case 14:
		return valueobject.TypePsychic
	case 15:
		return valueobject.TypeIce
	case 16:
		return valueobject.TypeDragon
	case 17:
		return valueobject.TypeDark
	case 18:
		return valueobject.TypeFairy
	default:
		return ""
	}
}

// isRechargeMove 判断是否为需要充能的技能
func isRechargeMove(moveID int) bool {
	// 需要充能的技能ID列表
	// 63: 破坏光线 (Hyper Beam)
	// 76: 日光束 (Solar Beam) - 蓄力技能，这里暂不处理
	// 143: 天空攻击 (Sky Attack) - 蓄力技能
	// 200: 愤怒门 (Outrage) - 连续技能
	// 304: 过热 (Overheat) - 不需要充能，只是降低特攻
	// 416: 终极冲击 (Giga Impact)
	// 800: 流星群 (Meteor Assault)
	// 801: 永恒之光 (Eternabeam)
	rechargeMoves := map[int]bool{
		63:  true, // 破坏光线
		416: true, // 终极冲击
		800: true, // 流星群
		801: true, // 永恒之光
		891: true, // 巨兽斩 (Behemoth Blade) - 不需要充能
		892: true, // 巨兽弹 (Behemoth Bash) - 不需要充能
	}
	// 实际需要充能的只有这几个
	rechargeMoves = map[int]bool{
		63:  true, // 破坏光线
		416: true, // 终极冲击
		800: true, // 流星突击
		801: true, // 无尽光束
	}
	return rechargeMoves[moveID]
}

// GetSpriteURL 获取宝可梦精灵图URL
func GetSpriteURL(pokemonID int) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/%d.gif", pokemonID)
}

// GetSpriteURLBack 获取宝可梦背面精灵图URL
func GetSpriteURLBack(pokemonID int) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/back/%d.gif", pokemonID)
}

// GetPokemonTypeString 获取属性字符串
func GetPokemonTypeString(types []valueobject.PokeType) string {
	var typeStrs []string
	for _, t := range types {
		typeStrs = append(typeStrs, string(t))
	}
	return strings.Join(typeStrs, "/")
}

// ===== 对外接口（使用默认客户端） =====

// EnsureDataLoaded 确保数据已加载
func EnsureDataLoaded() error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	return defaultClient.LoadAllData(ctx)
}

// IsDataLoaded 检查数据是否已加载
func IsDataLoaded() bool {
	defaultClient.cache.mu.RLock()
	defer defaultClient.cache.mu.RUnlock()
	return defaultClient.cache.Loaded
}

// GetPredefinedPokemon 获取宝可梦（兼容旧接口）
func GetPredefinedPokemon(id int) *entity.Pokemon {
	if !IsDataLoaded() {
		if err := EnsureDataLoaded(); err != nil {
			return nil
		}
	}

	defaultClient.cache.mu.RLock()
	pokemon := defaultClient.cache.Pokemon[id]
	defaultClient.cache.mu.RUnlock()

	if pokemon == nil {
		return nil
	}

	// 返回副本并填充技能
	return copyPokemonWithMoves(pokemon)
}

// GetAllPredefinedPokemon 获取所有宝可梦
func GetAllPredefinedPokemon() []*entity.Pokemon {
	if !IsDataLoaded() {
		if err := EnsureDataLoaded(); err != nil {
			return nil
		}
	}

	defaultClient.cache.mu.RLock()
	defer defaultClient.cache.mu.RUnlock()

	result := make([]*entity.Pokemon, 0, len(defaultClient.cache.Pokemon))
	for _, p := range defaultClient.cache.Pokemon {
		result = append(result, copyPokemonWithMoves(p))
	}

	// 按ID排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

// SearchPredefinedPokemon 搜索宝可梦
func SearchPredefinedPokemon(keyword string) []*entity.Pokemon {
	if !IsDataLoaded() {
		if err := EnsureDataLoaded(); err != nil {
			return nil
		}
	}

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	defaultClient.cache.mu.RLock()
	defer defaultClient.cache.mu.RUnlock()

	var results []*entity.Pokemon

	// 尝试按ID搜索
	if id, err := strconv.Atoi(keyword); err == nil && id > 0 {
		if p := defaultClient.cache.Pokemon[id]; p != nil {
			results = append(results, copyPokemonWithMoves(p))
			return results
		}
	}

	// 按名称搜索
	for _, p := range defaultClient.cache.Pokemon {
		if strings.Contains(p.Name, keyword) {
			results = append(results, copyPokemonWithMoves(p))
		}
	}

	// 按ID排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	// 限制结果数量
	if len(results) > 25 {
		results = results[:25]
	}

	return results
}

// GetPokemonByKeyword 通过关键词获取宝可梦
func GetPokemonByKeyword(keyword string) *entity.Pokemon {
	results := SearchPredefinedPokemon(keyword)
	if len(results) > 0 {
		return results[0]
	}
	return nil
}

// copyPokemonWithMoves 复制宝可梦并填充技能和特性
func copyPokemonWithMoves(p *entity.Pokemon) *entity.Pokemon {
	newP := &entity.Pokemon{
		ID:        p.ID,
		Name:      p.Name,
		Types:     append([]valueobject.PokeType{}, p.Types...),
		BaseHP:    p.BaseHP,
		BaseAtk:   p.BaseAtk,
		BaseDef:   p.BaseDef,
		BaseSpAtk: p.BaseSpAtk,
		BaseSpDef: p.BaseSpDef,
		BaseSpeed: p.BaseSpeed,
		SpriteURL: p.SpriteURL,
	}

	// 填充可学技能
	moveIDs := defaultClient.cache.PokemonMoves[p.ID]
	for _, moveID := range moveIDs {
		moveData := defaultClient.cache.Moves[moveID]
		moveName := defaultClient.cache.MoveNames[moveID]
		if moveData != nil && moveName != "" {
			move := &entity.Move{
				Name:             moveName,
				Type:             moveData.Type,
				Category:         moveData.Category,
				Power:            moveData.Power,
				Accuracy:         moveData.Accuracy,
				PP:               moveData.PP,
				MaxPP:            moveData.MaxPP,
				Priority:         moveData.Priority,
				RechargeRequired: moveData.RechargeRequired,
			}
			newP.LearnableMoves = append(newP.LearnableMoves, move)
		}
	}

	// 如果没有技能，添加默认技能
	if len(newP.LearnableMoves) == 0 {
		newP.LearnableMoves = append(newP.LearnableMoves, &entity.Move{
			Name:     "撞击",
			Type:     valueobject.TypeNormal,
			Category: entity.CategoryPhysical,
			Power:    40,
			Accuracy: 100,
			PP:       35,
			MaxPP:    35,
		})
	}

	// 填充特性（使用详细信息）
	abilityInfos := defaultClient.cache.PokemonAbilityInfos[p.ID]
	for _, info := range abilityInfos {
		abilityName := defaultClient.cache.AbilityNames[info.AbilityID]
		abilityDesc := defaultClient.cache.AbilityDescriptions[info.AbilityID]
		if abilityName != "" {
			ability := valueobject.Ability{
				ID:          info.AbilityID,
				Name:        abilityName,
				Description: abilityDesc,
				IsHidden:    info.IsHidden,
			}
			if info.IsHidden {
				newP.HiddenAbility = &ability
			} else {
				newP.Abilities = append(newP.Abilities, ability)
			}
		}
	}

	// 如果没有特性，添加默认特性
	if len(newP.Abilities) == 0 {
		newP.Abilities = append(newP.Abilities, valueobject.Ability{
			ID:   1,
			Name: "恶臭",
		})
	}

	return newP
}

// GetTotalPokemonCount 获取宝可梦总数
func GetTotalPokemonCount() int {
	if !IsDataLoaded() {
		return 0
	}
	defaultClient.cache.mu.RLock()
	defer defaultClient.cache.mu.RUnlock()
	return len(defaultClient.cache.Pokemon)
}