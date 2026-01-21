package commands

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	pokemon_app "github.com/user/dcminigames/internal/application/pokemon"
	"github.com/user/dcminigames/internal/domain/pokemon/entity"
	"github.com/user/dcminigames/internal/domain/pokemon/valueobject"
	"github.com/user/dcminigames/internal/infrastructure/discord"
	"github.com/user/dcminigames/internal/infrastructure/pokeapi"
)

// PokemonCommands 宝可梦对战命令处理器
type PokemonCommands struct {
	bot     *discord.Bot
	handler *pokemon_app.Handler
}

// NewPokemonCommands 创建命令处理器
func NewPokemonCommands(bot *discord.Bot, handler *pokemon_app.Handler) *PokemonCommands {
	return &PokemonCommands{bot: bot, handler: handler}
}

// Commands 返回斜杠命令定义
func (c *PokemonCommands) Commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "pokemon",
			Description: "宝可梦对战",
		},
	}
}

// HandleInteraction 处理交互
func (c *PokemonCommands) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		if data.Name == "pokemon" {
			c.showPanel(i)
		}
	} else if i.Type == discordgo.InteractionMessageComponent {
		c.handleComponent(i)
	} else if i.Type == discordgo.InteractionModalSubmit {
		data := i.ModalSubmitData()
		if data.CustomID == "pkm:search_modal" {
			c.handleSearchModal(i)
		} else if data.CustomID == "pkm:savepreset_modal" {
			c.handleSavePresetSubmit(i)
		} else if strings.HasPrefix(data.CustomID, "pkm:searchmove_modal:") {
			parts := strings.Split(data.CustomID, ":")
			if len(parts) >= 3 {
				c.handleSearchMoveModalSubmit(i, parts[2])
			}
		}
	}
}

// showPanel 显示主面板
func (c *PokemonCommands) showPanel(i *discordgo.InteractionCreate) {
	channelID := i.ChannelID
	userID := i.Member.User.ID

	battle, err := c.handler.GetBattle(channelID)
	var embed *discordgo.MessageEmbed
	var components []discordgo.MessageComponent

	if err != nil {
		// 没有进行中的对战
		embed = &discordgo.MessageEmbed{
			Title:       "⚔️ 宝可梦对战",
			Description: "当前没有进行中的对战\n选择对战模式创建新对战：\n\n**🎮 PVP 对战**\n• 单挑 (1v1) / 3v3 / 6v6\n\n**🤖 人机对战** (Debug 模式)\n• 与 AI 训练师对战，方便调试",
			Color:       0xFFCB05,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/25.gif",
			},
		}
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "⚔️ 单挑 (1v1)", Style: discordgo.SuccessButton, CustomID: "pkm:create:1"},
					discordgo.Button{Label: "⚔️ 3v3 单打", Style: discordgo.PrimaryButton, CustomID: "pkm:create:3"},
					discordgo.Button{Label: "⚔️ 6v6 单打", Style: discordgo.DangerButton, CustomID: "pkm:create:6"},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "🤖 人机 1v1", Style: discordgo.SecondaryButton, CustomID: "pkm:ai:1"},
					discordgo.Button{Label: "🤖 人机 3v3", Style: discordgo.SecondaryButton, CustomID: "pkm:ai:3"},
					discordgo.Button{Label: "🤖 人机 6v6", Style: discordgo.SecondaryButton, CustomID: "pkm:ai:6"},
				},
			},
		}
	} else {
		embed, components = c.buildBattlePanel(battle, userID)
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, components, true)
}

// buildBattlePanel 构建对战面板
func (c *PokemonCommands) buildBattlePanel(battle *entity.Battle, userID string) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	var embed *discordgo.MessageEmbed
	var components []discordgo.MessageComponent

	player := battle.GetPlayer(userID)
	isInBattle := player != nil
	isHost := battle.Player1 != nil && battle.Player1.ID == userID

	switch battle.State {
	case entity.BattleStateWaiting:
		modeName := battle.TeamSize.GetDisplayName()
		embed = &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("⚔️ 宝可梦对战 - %s", modeName),
			Description: fmt.Sprintf("对战ID: `%s`\n模式: **%s**\n\n**玩家1:** %s\n**玩家2:** 等待中...", battle.ID[:8], modeName, battle.Player1.Username),
			Color:       0xFFCB05,
		}
		var buttons []discordgo.MessageComponent
		if !isInBattle {
			buttons = append(buttons, discordgo.Button{Label: "⚔️ 加入对战", Style: discordgo.SuccessButton, CustomID: "pkm:join"})
		}
		buttons = append(buttons, discordgo.Button{Label: "🔄 刷新", Style: discordgo.SecondaryButton, CustomID: "pkm:refresh"})
		if isHost {
			buttons = append(buttons, discordgo.Button{Label: "❌ 取消", Style: discordgo.DangerButton, CustomID: "pkm:end"})
		}
		components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: buttons}}

	case entity.BattleStateChoosing:
		embed = &discordgo.MessageEmbed{
			Title:       "⚔️ 宝可梦对战 - 选择宝可梦",
			Description: battle.GetBattleStatus(),
			Color:       0xFFCB05,
		}
		var buttons []discordgo.MessageComponent
		if isInBattle && !player.Ready {
			buttons = append(buttons, discordgo.Button{Label: "🎮 选择宝可梦", Style: discordgo.PrimaryButton, CustomID: "pkm:select"})
		}
		buttons = append(buttons, discordgo.Button{Label: "🔄 刷新", Style: discordgo.SecondaryButton, CustomID: "pkm:refresh"})
		if isHost {
			buttons = append(buttons, discordgo.Button{Label: "❌ 取消", Style: discordgo.DangerButton, CustomID: "pkm:end"})
		}
		components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: buttons}}

	case entity.BattleStateBattling:
		embed = c.buildBattleStatusEmbed(battle, userID)
		var buttons []discordgo.MessageComponent
		if isInBattle {
			if player.Action == nil {
				buttons = append(buttons, discordgo.Button{Label: "⚡ 技能", Style: discordgo.PrimaryButton, CustomID: "pkm:moves"})
				// 3v3/6v6 模式下显示换人按钮
				if battle.TeamSize > 1 && player.HasSwitchableTeamMember() {
					buttons = append(buttons, discordgo.Button{Label: "🔄 换人", Style: discordgo.SecondaryButton, CustomID: "pkm:switch"})
				}
				buttons = append(buttons, discordgo.Button{Label: "🏳️ 认输", Style: discordgo.DangerButton, CustomID: "pkm:forfeit"})
			} else {
				buttons = append(buttons, discordgo.Button{Label: "⏳ 等待对手...", Style: discordgo.SecondaryButton, CustomID: "pkm:waiting", Disabled: true})
			}
		}
		buttons = append(buttons, discordgo.Button{Label: "🔃 刷新", Style: discordgo.SecondaryButton, CustomID: "pkm:refresh"})
		components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: buttons}}

	case entity.BattleStateFinished:
		winnerName := "无"
		if battle.Winner != nil {
			winnerName = battle.Winner.Username
		}
		embed = &discordgo.MessageEmbed{
			Title:       "🏆 对战结束",
			Description: fmt.Sprintf("**%s** 获胜！", winnerName),
			Color:       0x00D166,
		}
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "⚔️ 新对战", Style: discordgo.SuccessButton, CustomID: "pkm:create"},
				},
			},
		}
	}

	return embed, components
}

// buildBattleStatusEmbed 构建对战状态Embed
func (c *PokemonCommands) buildBattleStatusEmbed(battle *entity.Battle, userID string) *discordgo.MessageEmbed {
	p1 := battle.Player1
	p2 := battle.Player2

	// 构建HP条
	p1HP := c.buildHPBar(p1.Pokemon)
	p2HP := c.buildHPBar(p2.Pokemon)

	// 获取最近的战斗日志
	logs := ""
	if len(battle.Logs) > 0 {
		start := len(battle.Logs) - 8
		if start < 0 {
			start = 0
		}
		logs = strings.Join(battle.Logs[start:], "\n")
	}

	// 判断当前状态
	status := ""
	player := battle.GetPlayer(userID)
	if player != nil {
		if player.Action == nil {
			status = "💡 请选择你的行动！"
		} else {
			status = "⏳ 等待对手行动..."
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("⚔️ 回合 %d", battle.CurrentTurn),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("🔴 %s 的 %s", p1.Username, p1.Pokemon.Pokemon.Name),
				Value:  fmt.Sprintf("Lv.%d %s\n%s", p1.Pokemon.Level, pokeapi.GetPokemonTypeString(p1.Pokemon.Pokemon.Types), p1HP),
				Inline: true,
			},
			{
				Name:   "VS",
				Value:  "⚔️",
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("🔵 %s 的 %s", p2.Username, p2.Pokemon.Pokemon.Name),
				Value:  fmt.Sprintf("Lv.%d %s\n%s", p2.Pokemon.Level, pokeapi.GetPokemonTypeString(p2.Pokemon.Pokemon.Types), p2HP),
				Inline: true,
			},
		},
		Color: 0xFFCB05,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: p1.Pokemon.Pokemon.GetSpriteURL(),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: p2.Pokemon.Pokemon.GetSpriteURL(),
		},
	}

	if logs != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "📜 战斗日志",
			Value: logs,
		})
	}

	if status != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{Text: status}
	}

	return embed
}

// buildHPBar 构建HP条
func (c *PokemonCommands) buildHPBar(battler *entity.Battler) string {
	percent := battler.GetHPPercent()
	barLength := 10
	filled := int(percent / 10)
	if filled > barLength {
		filled = barLength
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("🟩", filled) + strings.Repeat("⬜", barLength-filled)
	return fmt.Sprintf("%s\n❤️ %d/%d (%.0f%%)", bar, battler.CurrentHP, battler.MaxHP, percent)
}

// handleComponent 处理组件交互
func (c *PokemonCommands) handleComponent(i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	channelID := i.ChannelID
	userID := i.Member.User.ID
	username := i.Member.User.Username

	parts := strings.Split(customID, ":")
	if len(parts) < 2 {
		return
	}

	prefix := parts[0]
	action := parts[1]

	if prefix != "pkm" {
		return
	}

	switch action {
	case "create":
		teamSize := 1
		if len(parts) >= 3 {
			if size, err := strconv.Atoi(parts[2]); err == nil {
				teamSize = size
			}
		}
		c.handleCreate(i, channelID, userID, username, teamSize)
	case "ai":
		teamSize := 1
		if len(parts) >= 3 {
			if size, err := strconv.Atoi(parts[2]); err == nil {
				teamSize = size
			}
		}
		c.handleCreateAI(i, channelID, userID, username, teamSize)
	case "join":
		c.handleJoin(i, channelID, userID, username)
	case "select":
		c.handleSelectMenu(i)
	case "search":
		c.handleSearch(i)
	case "browse":
		pageStr := "1"
		if len(parts) >= 3 {
			pageStr = parts[2]
		}
		c.handleBrowse(i, pageStr)
	case "choose":
		if len(parts) >= 3 {
			c.handleChoosePokemon(i, channelID, userID, parts[2])
		}
	case "nature":
		if len(parts) >= 3 {
			c.handleNatureSelect(i, channelID, userID, parts[2])
		}
	case "setnature":
		if len(parts) >= 4 {
			c.handleSetNature(i, channelID, userID, parts[2], parts[3])
		}
	case "ability":
		if len(parts) >= 3 {
			c.handleAbilitySelect(i, channelID, userID, parts[2])
		}
	case "setability":
		if len(parts) >= 4 {
			c.handleSetAbility(i, channelID, userID, parts[2], parts[3])
		}
	case "confirm":
		if len(parts) >= 3 {
			c.handleConfirmPokemon(i, channelID, userID, parts[2])
		}
	case "cfgmoves":
		if len(parts) >= 3 {
			pageStr := "1"
			if len(parts) >= 4 {
				pageStr = parts[3]
			}
			c.handleConfigMoves(i, channelID, userID, parts[2], pageStr)
		}
	case "setmove":
		if len(parts) >= 5 {
			c.handleSetMove(i, channelID, userID, parts[2], parts[3], parts[4])
		}
	case "confirmmoves":
		if len(parts) >= 3 {
			c.handleConfirmMoves(i, channelID, userID, parts[2])
		}
	case "searchmove":
		if len(parts) >= 3 {
			c.handleSearchMoveModal(i, parts[2])
		}
	case "selectmove":
		if len(parts) >= 4 {
			c.handleSelectSearchedMove(i, channelID, userID, parts[2], parts[3])
		}
	case "presets":
		c.handleShowPresets(i, userID)
	case "loadpreset":
		if len(parts) >= 3 {
			c.handleLoadPreset(i, channelID, userID, parts[2])
		}
	case "savepreset":
		c.handleSavePresetModal(i, channelID, userID)
	case "delpreset":
		if len(parts) >= 3 {
			c.handleDeletePreset(i, userID, parts[2])
		}
	case "moves":
		c.handleShowMoves(i, channelID, userID)
	case "move":
		if len(parts) >= 3 {
			c.handleUseMove(i, channelID, userID, parts[2])
		}
	case "forfeit":
		c.handleForfeit(i, channelID, userID)
	case "refresh":
		c.handleRefresh(i, channelID, userID)
	case "end":
		c.handleEnd(i, channelID, userID)
	case "switch":
		c.handleShowSwitchMenu(i, channelID, userID)
	case "doswitch":
		if len(parts) >= 3 {
			c.handleDoSwitch(i, channelID, userID, parts[2])
		}
	case "forceswitch":
		if len(parts) >= 3 {
			c.handleForceSwitch(i, channelID, userID, parts[2])
		}
	}
}

// handleCreate 创建对战
func (c *PokemonCommands) handleCreate(i *discordgo.InteractionCreate, channelID, userID, username string, teamSize int) {
	// 先结束可能存在的旧对战
	c.handler.EndBattle(channelID)

	// 转换 teamSize 到 entity.TeamSize
	var ts entity.TeamSize
	switch teamSize {
	case 3:
		ts = entity.TeamSize3v3
	case 6:
		ts = entity.TeamSize6v6
	default:
		ts = entity.TeamSize1v1
	}

	battle, err := c.handler.CreateBattleWithTeamSize(channelID, userID, username, ts)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	modeName := ts.GetDisplayName()
	c.bot.RespondPublic(i.Interaction, fmt.Sprintf("⚔️ **%s** 创建了 **%s** 宝可梦对战！\n对战ID: `%s`\n使用 `/pokemon` 加入对战", username, modeName, battle.ID[:8]))
}

// handleCreateAI 创建人机对战
func (c *PokemonCommands) handleCreateAI(i *discordgo.InteractionCreate, channelID, userID, username string, teamSize int) {
	// 先结束可能存在的旧对战
	c.handler.EndBattle(channelID)

	// 转换 teamSize 到 entity.TeamSize
	var ts entity.TeamSize
	switch teamSize {
	case 3:
		ts = entity.TeamSize3v3
	case 6:
		ts = entity.TeamSize6v6
	default:
		ts = entity.TeamSize1v1
	}

	battle, err := c.handler.CreateAIBattle(channelID, userID, username, ts)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	// 获取 AI 选择的宝可梦信息
	aiPlayer := battle.GetAIPlayer()
	aiPokemonInfo := ""
	if aiPlayer != nil && len(aiPlayer.Team) > 0 {
		var names []string
		for _, battler := range aiPlayer.Team {
			names = append(names, battler.Pokemon.Name)
		}
		aiPokemonInfo = fmt.Sprintf("\n🤖 AI 已选择: **%s**", strings.Join(names, "、"))
	}

	modeName := ts.GetDisplayName()
	c.bot.RespondPublic(i.Interaction, fmt.Sprintf("🤖 **%s** 创建了 **%s** 人机对战！\n对战ID: `%s`%s\n\n请选择你的宝可梦开始对战！", username, modeName, battle.ID[:8], aiPokemonInfo))
}

// handleJoin 加入对战
func (c *PokemonCommands) handleJoin(i *discordgo.InteractionCreate, channelID, userID, username string) {
	if err := c.handler.JoinBattle(channelID, userID, username); err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	c.bot.RespondPublic(i.Interaction, fmt.Sprintf("✅ **%s** 加入了对战！\n双方请选择宝可梦", username))
}

// handleSelectMenu 显示宝可梦选择菜单（私密）
func (c *PokemonCommands) handleSelectMenu(i *discordgo.InteractionCreate) {
	channelID := i.ChannelID
	userID := i.Member.User.ID

	// 获取对战信息以显示队伍选择进度
	battle, _ := c.handler.GetBattle(channelID)
	var progressInfo string
	if battle != nil {
		player := battle.GetPlayer(userID)
		if player != nil {
			current := len(player.Team)
			total := int(battle.TeamSize)
			progressInfo = fmt.Sprintf("\n\n**📋 队伍进度: %d/%d**", current, total)
			if current > 0 {
				progressInfo += "\n已选择: "
				for idx, battler := range player.Team {
					if idx > 0 {
						progressInfo += ", "
					}
					progressInfo += battler.Pokemon.Name
				}
			}
		}
	}

	// 显示搜索提示和热门宝可梦
	embed := &discordgo.MessageEmbed{
		Title:       "🎮 选择你的宝可梦",
		Description: "**搜索方式：**\n• 输入宝可梦名称（如：皮卡丘）\n• 输入图鉴编号（如：25）\n\n点击下方按钮搜索或直接选择热门宝可梦：" + progressInfo,
		Color:       0xFFCB05,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "💡 选择过程仅你可见",
		},
	}

	// 热门宝可梦快捷按钮
	popularButtons := []discordgo.MessageComponent{
		discordgo.Button{Label: "皮卡丘", Style: discordgo.PrimaryButton, CustomID: "pkm:choose:25", Emoji: &discordgo.ComponentEmoji{Name: "⚡"}},
		discordgo.Button{Label: "喷火龙", Style: discordgo.DangerButton, CustomID: "pkm:choose:6", Emoji: &discordgo.ComponentEmoji{Name: "🔥"}},
		discordgo.Button{Label: "水箭龟", Style: discordgo.PrimaryButton, CustomID: "pkm:choose:9", Emoji: &discordgo.ComponentEmoji{Name: "💧"}},
		discordgo.Button{Label: "妙蛙花", Style: discordgo.SuccessButton, CustomID: "pkm:choose:3", Emoji: &discordgo.ComponentEmoji{Name: "🌿"}},
		discordgo.Button{Label: "超梦", Style: discordgo.SecondaryButton, CustomID: "pkm:choose:150", Emoji: &discordgo.ComponentEmoji{Name: "🔮"}},
	}

	// 搜索和浏览按钮
	actionButtons := []discordgo.MessageComponent{
		discordgo.Button{Label: "🔍 搜索宝可梦", Style: discordgo.PrimaryButton, CustomID: "pkm:search"},
		discordgo.Button{Label: "📖 浏览图鉴", Style: discordgo.SecondaryButton, CustomID: "pkm:browse:1"},
	}

	rows := []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: popularButtons},
		discordgo.ActionsRow{Components: actionButtons},
	}

	// 私密响应
	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleSearch 处理搜索请求（显示模态框）
func (c *PokemonCommands) handleSearch(i *discordgo.InteractionCreate) {
	err := c.bot.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "pkm:search_modal",
			Title:    "搜索宝可梦",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "keyword",
							Label:       "输入宝可梦名称或图鉴编号",
							Style:       discordgo.TextInputShort,
							Placeholder: "例如：皮卡丘 或 25",
							Required:    true,
							MinLength:   1,
							MaxLength:   50,
						},
					},
				},
			},
		},
	})
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无法打开搜索框")
	}
}

// handleSearchModal 处理搜索模态框提交
func (c *PokemonCommands) handleSearchModal(i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	keyword := ""
	for _, comp := range data.Components {
		if row, ok := comp.(*discordgo.ActionsRow); ok {
			for _, c := range row.Components {
				if input, ok := c.(*discordgo.TextInput); ok && input.CustomID == "keyword" {
					keyword = input.Value
				}
			}
		}
	}

	if keyword == "" {
		c.bot.RespondEphemeral(i.Interaction, "❌ 请输入搜索关键词")
		return
	}

	// 搜索宝可梦
	results := pokeapi.SearchPredefinedPokemon(keyword)
	if len(results) == 0 {
		c.bot.RespondEphemeral(i.Interaction, fmt.Sprintf("❌ 未找到匹配「%s」的宝可梦", keyword))
		return
	}

	// 限制结果数量
	if len(results) > 10 {
		results = results[:10]
	}

	// 构建搜索结果
	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("搜索「%s」找到 %d 个结果：\n\n", keyword, len(results)))

	var buttons []discordgo.MessageComponent
	for _, p := range results {
		typeStr := pokeapi.GetPokemonTypeString(p.Types)
		desc.WriteString(fmt.Sprintf("**#%03d %s** (%s)\n", p.ID, p.Name, typeStr))
		buttons = append(buttons, discordgo.Button{
			Label:    fmt.Sprintf("#%d %s", p.ID, p.Name),
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("pkm:choose:%d", p.ID),
		})
	}

	// 每行最多5个按钮
	var rows []discordgo.MessageComponent
	for i := 0; i < len(buttons); i += 5 {
		end := i + 5
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, discordgo.ActionsRow{Components: buttons[i:end]})
	}

	// 添加返回按钮
	rows = append(rows, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "🔙 返回", Style: discordgo.SecondaryButton, CustomID: "pkm:select"},
		},
	})

	embed := &discordgo.MessageEmbed{
		Title:       "🔍 搜索结果",
		Description: desc.String(),
		Color:       0xFFCB05,
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleBrowse 浏览图鉴
func (c *PokemonCommands) handleBrowse(i *discordgo.InteractionCreate, pageStr string) {
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pokemons := c.handler.GetAvailablePokemon()
	sort.Slice(pokemons, func(a, b int) bool {
		return pokemons[a].ID < pokemons[b].ID
	})

	// 每页10个
	perPage := 10
	totalPages := (len(pokemons) + perPage - 1) / perPage
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	end := start + perPage
	if end > len(pokemons) {
		end = len(pokemons)
	}

	pagePokemons := pokemons[start:end]

	// 构建列表
	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("📖 第 %d/%d 页\n\n", page, totalPages))

	var buttons []discordgo.MessageComponent
	for _, p := range pagePokemons {
		typeStr := pokeapi.GetPokemonTypeString(p.Types)
		desc.WriteString(fmt.Sprintf("**#%03d %s** (%s)\n", p.ID, p.Name, typeStr))
		buttons = append(buttons, discordgo.Button{
			Label:    fmt.Sprintf("#%d %s", p.ID, p.Name),
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("pkm:choose:%d", p.ID),
		})
	}

	// 每行最多5个按钮
	var rows []discordgo.MessageComponent
	for i := 0; i < len(buttons); i += 5 {
		end := i + 5
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, discordgo.ActionsRow{Components: buttons[i:end]})
	}

	// 分页按钮
	var navButtons []discordgo.MessageComponent
	if page > 1 {
		navButtons = append(navButtons, discordgo.Button{Label: "⬅️ 上一页", Style: discordgo.SecondaryButton, CustomID: fmt.Sprintf("pkm:browse:%d", page-1)})
	}
	navButtons = append(navButtons, discordgo.Button{Label: "🔍 搜索", Style: discordgo.PrimaryButton, CustomID: "pkm:search"})
	if page < totalPages {
		navButtons = append(navButtons, discordgo.Button{Label: "➡️ 下一页", Style: discordgo.SecondaryButton, CustomID: fmt.Sprintf("pkm:browse:%d", page+1)})
	}
	navButtons = append(navButtons, discordgo.Button{Label: "🔙 返回", Style: discordgo.SecondaryButton, CustomID: "pkm:select"})
	rows = append(rows, discordgo.ActionsRow{Components: navButtons})

	embed := &discordgo.MessageEmbed{
		Title:       "📖 宝可梦图鉴",
		Description: desc.String(),
		Color:       0xFFCB05,
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleChoosePokemon 选择宝可梦（显示配置界面）
func (c *PokemonCommands) handleChoosePokemon(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr string) {
	pokemonID, err := strconv.Atoi(pokemonIDStr)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无效的宝可梦")
		return
	}

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	// 初始化配置
	config := &pokemon_app.PokemonConfig{
		PokemonID:   pokemonID,
		Nature:      valueobject.NatureHardy,
		AbilitySlot: 0,
		MoveIndices: []int{0, 1, 2, 3},
	}
	c.handler.SetConfig(channelID, userID, config)

	// 显示配置界面
	c.showConfigPanel(i, channelID, userID, pokemon, config)
}

// showConfigPanel 显示宝可梦配置面板
func (c *PokemonCommands) showConfigPanel(i *discordgo.InteractionCreate, channelID, userID string, pokemon *entity.Pokemon, config *pokemon_app.PokemonConfig) {
	typeStr := pokeapi.GetPokemonTypeString(pokemon.Types)
	
	// 构建描述
	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("**#%03d %s** (%s)\n\n", pokemon.ID, pokemon.Name, typeStr))
	desc.WriteString(fmt.Sprintf("📊 **种族值**: HP %d / 攻 %d / 防 %d / 特攻 %d / 特防 %d / 速 %d\n\n",
		pokemon.BaseHP, pokemon.BaseAtk, pokemon.BaseDef, pokemon.BaseSpAtk, pokemon.BaseSpDef, pokemon.BaseSpeed))
	
	// 当前配置
	natureMod := valueobject.GetNatureModifier(config.Nature)
	desc.WriteString(fmt.Sprintf("🎭 **性格**: %s (%s)\n", config.Nature, formatNatureEffect(natureMod)))
	
	// 显示特性
	if config.AbilitySlot == -1 && pokemon.HiddenAbility != nil {
		// 隐藏特性
		desc.WriteString(fmt.Sprintf("✨ **特性**: %s (隐藏)\n", pokemon.HiddenAbility.Name))
	} else if config.AbilitySlot >= 0 && config.AbilitySlot < len(pokemon.Abilities) {
		// 普通特性
		desc.WriteString(fmt.Sprintf("✨ **特性**: %s\n", pokemon.Abilities[config.AbilitySlot].Name))
	} else if len(pokemon.Abilities) > 0 {
		// 默认显示第一个特性
		desc.WriteString(fmt.Sprintf("✨ **特性**: %s\n", pokemon.Abilities[0].Name))
	}
	
	desc.WriteString("\n**技能**:\n")
	for idx, moveIdx := range config.MoveIndices {
		if moveIdx < len(pokemon.LearnableMoves) {
			m := pokemon.LearnableMoves[moveIdx]
			desc.WriteString(fmt.Sprintf("%d. %s (%s, 威力%d)\n", idx+1, m.Name, m.Type, m.Power))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "⚙️ 配置你的宝可梦",
		Description: desc.String(),
		Color:       0xFFCB05,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: c.handler.GetSpriteURL(pokemon.ID)},
		Footer:      &discordgo.MessageEmbedFooter{Text: "💡 配置完成后点击「确认选择」"},
	}

	// 配置按钮
	rows := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "🎭 选择性格", Style: discordgo.PrimaryButton, CustomID: fmt.Sprintf("pkm:nature:%d", pokemon.ID)},
				discordgo.Button{Label: "✨ 选择特性", Style: discordgo.PrimaryButton, CustomID: fmt.Sprintf("pkm:ability:%d", pokemon.ID)},
				discordgo.Button{Label: "⚔️ 选择技能", Style: discordgo.PrimaryButton, CustomID: fmt.Sprintf("pkm:cfgmoves:%d:1", pokemon.ID)},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "✅ 确认选择", Style: discordgo.SuccessButton, CustomID: fmt.Sprintf("pkm:confirm:%d", pokemon.ID)},
				discordgo.Button{Label: "💾 保存预设", Style: discordgo.SecondaryButton, CustomID: fmt.Sprintf("pkm:savepreset:%d", pokemon.ID)},
				discordgo.Button{Label: "🔙 返回选择", Style: discordgo.SecondaryButton, CustomID: "pkm:select"},
			},
		},
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// formatNatureEffect 格式化性格效果
func formatNatureEffect(mod valueobject.NatureModifier) string {
	if mod.Atk == 1.0 && mod.Def == 1.0 && mod.SpAtk == 1.0 && mod.SpDef == 1.0 && mod.Speed == 1.0 {
		return "无修正"
	}
	var effects []string
	if mod.Atk > 1.0 {
		effects = append(effects, "攻击↑")
	} else if mod.Atk < 1.0 {
		effects = append(effects, "攻击↓")
	}
	if mod.Def > 1.0 {
		effects = append(effects, "防御↑")
	} else if mod.Def < 1.0 {
		effects = append(effects, "防御↓")
	}
	if mod.SpAtk > 1.0 {
		effects = append(effects, "特攻↑")
	} else if mod.SpAtk < 1.0 {
		effects = append(effects, "特攻↓")
	}
	if mod.SpDef > 1.0 {
		effects = append(effects, "特防↑")
	} else if mod.SpDef < 1.0 {
		effects = append(effects, "特防↓")
	}
	if mod.Speed > 1.0 {
		effects = append(effects, "速度↑")
	} else if mod.Speed < 1.0 {
		effects = append(effects, "速度↓")
	}
	return strings.Join(effects, " ")
}

// handleShowMoves 显示技能列表
func (c *PokemonCommands) handleShowMoves(i *discordgo.InteractionCreate, channelID, userID string) {
	battle, err := c.handler.GetBattle(channelID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 对战不存在")
		return
	}

	player := battle.GetPlayer(userID)
	if player == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 你不在对战中")
		return
	}

	if player.Action != nil {
		c.bot.RespondEphemeral(i.Interaction, "⏳ 你已选择行动，等待对手...")
		return
	}

	// 构建技能按钮
	var buttons []discordgo.MessageComponent
	for idx, move := range player.Pokemon.Moves {
		ppInfo := fmt.Sprintf("%d/%d", move.PP, move.MaxPP)
		label := fmt.Sprintf("%s (%s) %s", move.Name, move.Type, ppInfo)
		disabled := !move.CanUse()

		style := discordgo.PrimaryButton
		if move.Category == entity.CategoryPhysical {
			style = discordgo.DangerButton
		} else if move.Category == entity.CategoryStatus {
			style = discordgo.SecondaryButton
		}

		buttons = append(buttons, discordgo.Button{
			Label:    label,
			Style:    style,
			CustomID: fmt.Sprintf("pkm:move:%d", idx),
			Disabled: disabled,
		})
	}

	// 根据按钮数量动态构建行，避免数组越界
	var rows []discordgo.MessageComponent
	if len(buttons) == 0 {
		c.bot.RespondEphemeral(i.Interaction, "❌ 没有可用的技能")
		return
	} else if len(buttons) <= 2 {
		rows = []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: buttons},
		}
	} else if len(buttons) <= 4 {
		rows = []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: buttons[:2]},
			discordgo.ActionsRow{Components: buttons[2:]},
		}
	} else {
		// 超过4个技能时，只显示前4个
		rows = []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: buttons[:2]},
			discordgo.ActionsRow{Components: buttons[2:4]},
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("⚡ %s 的技能", player.Pokemon.Pokemon.Name),
		Description: "选择要使用的技能",
		Color:       0xFFCB05,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.Pokemon.Pokemon.GetSpriteURL(),
		},
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleUseMove 使用技能
func (c *PokemonCommands) handleUseMove(i *discordgo.InteractionCreate, channelID, userID, moveIndexStr string) {
	moveIndex, err := strconv.Atoi(moveIndexStr)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无效的技能")
		return
	}

	logs, err := c.handler.UseMove(channelID, userID, moveIndex)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	// 检查是否为人机对战，如果是则触发 AI 行动
	battle, _ := c.handler.GetBattle(channelID)
	if battle != nil && battle.IsAIBattle && len(logs) == 0 {
		// 玩家已行动，触发 AI 行动并执行回合
		aiLogs, _ := c.handler.ExecuteAITurn(channelID)
		logs = aiLogs
		// 重新获取对战状态
		battle, _ = c.handler.GetBattle(channelID)
	}

	if len(logs) > 0 {
		// 回合执行完毕，发送战斗日志
		logText := strings.Join(logs, "\n")

		if battle != nil && battle.State == entity.BattleStateFinished {
			c.bot.RespondPublic(i.Interaction, logText)
			c.handler.EndBattle(channelID)
		} else {
			c.bot.RespondPublic(i.Interaction, logText)
			c.sendBattlePanel(i, channelID)
		}
	} else {
		// 等待对手（普通 PVP 模式）
		c.bot.RespondEphemeral(i.Interaction, "✅ 已选择技能，等待对手行动...")
	}
}

// handleForfeit 认输
func (c *PokemonCommands) handleForfeit(i *discordgo.InteractionCreate, channelID, userID string) {
	logs, err := c.handler.Forfeit(channelID, userID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	logText := strings.Join(logs, "\n")
	c.bot.RespondPublic(i.Interaction, logText)
	c.handler.EndBattle(channelID)
}

// handleRefresh 刷新面板
func (c *PokemonCommands) handleRefresh(i *discordgo.InteractionCreate, channelID, userID string) {
	battle, err := c.handler.GetBattle(channelID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 没有进行中的对战")
		return
	}
	embed, components := c.buildBattlePanel(battle, userID)
	c.bot.UpdateWithEmbed(i.Interaction, embed, components)
}

// handleEnd 结束对战
func (c *PokemonCommands) handleEnd(i *discordgo.InteractionCreate, channelID, userID string) {
	battle, err := c.handler.GetBattle(channelID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 没有进行中的对战")
		return
	}

	// 只有房主可以结束
	if battle.Player1 == nil || battle.Player1.ID != userID {
		c.bot.RespondEphemeral(i.Interaction, "❌ 只有房主可以结束对战")
		return
	}

	c.handler.EndBattle(channelID)
	c.bot.RespondPublic(i.Interaction, "🛑 对战已结束")
}

// sendBattlePanel 发送对战面板到频道
func (c *PokemonCommands) sendBattlePanel(i *discordgo.InteractionCreate, channelID string) {
	battle, err := c.handler.GetBattle(channelID)
	if err != nil || battle.State != entity.BattleStateBattling {
		return
	}

	embed := c.buildBattleStatusEmbed(battle, "")
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "⚡ 技能", Style: discordgo.PrimaryButton, CustomID: "pkm:moves"},
				discordgo.Button{Label: "🔄 刷新", Style: discordgo.SecondaryButton, CustomID: "pkm:refresh"},
			},
		},
	}

	c.bot.SendChannelEmbed(channelID, embed, components)
}

// handleNatureSelect 显示性格选择菜单
func (c *PokemonCommands) handleNatureSelect(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)
	
	natures := []struct {
		Nature valueobject.Nature
		Desc   string
	}{
		{valueobject.NatureAdamant, "攻击↑ 特攻↓"},
		{valueobject.NatureJolly, "速度↑ 特攻↓"},
		{valueobject.NatureModest, "特攻↑ 攻击↓"},
		{valueobject.NatureTimid, "速度↑ 攻击↓"},
		{valueobject.NatureBold, "防御↑ 攻击↓"},
		{valueobject.NatureCalm, "特防↑ 攻击↓"},
		{valueobject.NatureCareful, "特防↑ 特攻↓"},
		{valueobject.NatureImpish, "防御↑ 特攻↓"},
		{valueobject.NatureHardy, "无修正"},
	}

	var buttons []discordgo.MessageComponent
	for _, n := range natures {
		buttons = append(buttons, discordgo.Button{
			Label:    fmt.Sprintf("%s (%s)", n.Nature, n.Desc),
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("pkm:setnature:%d:%s", pokemonID, n.Nature),
		})
	}

	var rows []discordgo.MessageComponent
	for j := 0; j < len(buttons); j += 3 {
		end := j + 3
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, discordgo.ActionsRow{Components: buttons[j:end]})
	}

	rows = append(rows, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "🔙 返回配置", Style: discordgo.SecondaryButton, CustomID: fmt.Sprintf("pkm:choose:%d", pokemonID)},
		},
	})

	embed := &discordgo.MessageEmbed{
		Title:       "🎭 选择性格",
		Description: "性格会影响宝可梦的能力值成长（+10%/-10%）",
		Color:       0xFFCB05,
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleSetNature 设置性格
func (c *PokemonCommands) handleSetNature(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr, natureStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)
	nature := valueobject.Nature(natureStr)

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 配置已过期，请重新选择宝可梦")
		return
	}

	config.Nature = nature
	c.handler.SetConfig(channelID, userID, config)

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	c.showConfigPanel(i, channelID, userID, pokemon, config)
}

// handleAbilitySelect 显示特性选择菜单
func (c *PokemonCommands) handleAbilitySelect(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)
	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	var buttons []discordgo.MessageComponent
	for idx, ability := range pokemon.Abilities {
		label := ability.Name
		if ability.Description != "" {
			label = fmt.Sprintf("%s: %s", ability.Name, ability.Description)
			if len(label) > 80 {
				label = label[:77] + "..."
			}
		}
		buttons = append(buttons, discordgo.Button{
			Label:    label,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("pkm:setability:%d:%d", pokemonID, idx),
		})
	}

	if pokemon.HiddenAbility != nil {
		label := fmt.Sprintf("[隐藏] %s", pokemon.HiddenAbility.Name)
		buttons = append(buttons, discordgo.Button{
			Label:    label,
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("pkm:setability:%d:hidden", pokemonID),
		})
	}

	var rows []discordgo.MessageComponent
	for _, btn := range buttons {
		rows = append(rows, discordgo.ActionsRow{Components: []discordgo.MessageComponent{btn}})
	}

	rows = append(rows, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "🔙 返回配置", Style: discordgo.SecondaryButton, CustomID: fmt.Sprintf("pkm:choose:%d", pokemonID)},
		},
	})

	embed := &discordgo.MessageEmbed{
		Title:       "✨ 选择特性",
		Description: "特性会在对战中产生特殊效果",
		Color:       0xFFCB05,
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleSetAbility 设置特性
func (c *PokemonCommands) handleSetAbility(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr, slotStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 配置已过期，请重新选择宝可梦")
		return
	}

	if slotStr == "hidden" {
		config.AbilitySlot = -1
	} else {
		slot, _ := strconv.Atoi(slotStr)
		config.AbilitySlot = slot
	}
	c.handler.SetConfig(channelID, userID, config)

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	c.showConfigPanel(i, channelID, userID, pokemon, config)
}

// handleConfirmPokemon 确认选择宝可梦
func (c *PokemonCommands) handleConfirmPokemon(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	level := 50

	if err := c.handler.SelectPokemon(channelID, userID, pokemonID, level); err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	battle, _ := c.handler.GetBattle(channelID)
	if battle == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 对战不存在")
		return
	}

	player := battle.GetPlayer(userID)
	if player == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 你不在对战中")
		return
	}

	teamSize := int(battle.TeamSize)
	currentCount := len(player.Team)

	if battle.State == entity.BattleStateBattling {
		// 双方都已准备好，对战开始
		c.bot.RespondPublic(i.Interaction, fmt.Sprintf("✅ **%s** 选择了 **%s**！\n\n⚔️ 双方准备完毕，对战开始！", i.Member.User.Username, pokemon.Name))
		c.sendBattlePanel(i, channelID)
	} else if currentCount < teamSize {
		// 队伍未满，继续选择
		c.bot.RespondEphemeral(i.Interaction, fmt.Sprintf("✅ 已添加 **%s** 到队伍！\n\n📋 队伍进度: %d/%d\n请继续选择下一只宝可梦", pokemon.Name, currentCount, teamSize))
	} else if player.Ready {
		// 队伍已满且已准备
		c.bot.RespondPublic(i.Interaction, fmt.Sprintf("✅ **%s** 的队伍已准备完毕！等待对手...", i.Member.User.Username))
	} else {
		c.bot.RespondPublic(i.Interaction, fmt.Sprintf("✅ **%s** 选择了 **%s**！等待对手选择...", i.Member.User.Username, pokemon.Name))
	}
}

// handleConfigMoves 显示技能配置界面（支持分页）
func (c *PokemonCommands) handleConfigMoves(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr, pageStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		config = &pokemon_app.PokemonConfig{PokemonID: pokemonID, MoveIndices: []int{}}
		c.handler.SetConfig(channelID, userID, config)
	}

	// 每页显示12个技能（3行x4个）
	perPage := 12
	totalMoves := len(pokemon.LearnableMoves)
	totalPages := (totalMoves + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	end := start + perPage
	if end > totalMoves {
		end = totalMoves
	}

	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("## 🎯 %s 技能配置\n\n", pokemon.Name))
	desc.WriteString(fmt.Sprintf("**已选技能 (%d/4)：**\n", len(config.MoveIndices)))
	if len(config.MoveIndices) == 0 {
		desc.WriteString("_未选择技能（将使用默认技能）_\n")
	} else {
		for slot, idx := range config.MoveIndices {
			if idx < len(pokemon.LearnableMoves) {
				move := pokemon.LearnableMoves[idx]
				powerStr := "-"
				if move.Power > 0 {
					powerStr = fmt.Sprintf("%d", move.Power)
				}
				desc.WriteString(fmt.Sprintf("%d. %s (%s) 威力:%s\n", slot+1, move.Name, move.Type, powerStr))
			}
		}
	}
	desc.WriteString(fmt.Sprintf("\n📖 可学技能: %d个 (第%d/%d页)", totalMoves, page, totalPages))

	embed := &discordgo.MessageEmbed{
		Title:       "⚡ 技能配置",
		Description: desc.String(),
		Color:       0x3498DB,
		Footer:      &discordgo.MessageEmbedFooter{Text: "点击技能选择/取消，绿色为已选"},
	}

	var rows []discordgo.MessageComponent
	var buttons []discordgo.MessageComponent

	// 显示当前页的技能
	for idx := start; idx < end; idx++ {
		move := pokemon.LearnableMoves[idx]
		style := discordgo.SecondaryButton
		for _, si := range config.MoveIndices {
			if si == idx {
				style = discordgo.SuccessButton
				break
			}
		}
		// 截断过长的技能名
		label := move.Name
		if len(label) > 12 {
			label = label[:12] + "…"
		}
		buttons = append(buttons, discordgo.Button{
			Label:    label,
			Style:    style,
			CustomID: fmt.Sprintf("pkm:setmove:%d:%d:%d", pokemonID, idx, page),
		})
		if len(buttons) == 4 {
			rows = append(rows, discordgo.ActionsRow{Components: buttons})
			buttons = nil
		}
	}
	if len(buttons) > 0 {
		rows = append(rows, discordgo.ActionsRow{Components: buttons})
	}

	// 分页和操作按钮
	var navButtons []discordgo.MessageComponent
	if page > 1 {
		navButtons = append(navButtons, discordgo.Button{
			Label:    "⬅️",
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("pkm:cfgmoves:%d:%d", pokemonID, page-1),
		})
	}
	navButtons = append(navButtons, discordgo.Button{
		Label:    "🔍 搜索",
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("pkm:searchmove:%d", pokemonID),
	})
	if page < totalPages {
		navButtons = append(navButtons, discordgo.Button{
			Label:    "➡️",
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("pkm:cfgmoves:%d:%d", pokemonID, page+1),
		})
	}
	navButtons = append(navButtons, discordgo.Button{
		Label:    "🔙 返回",
		Style:    discordgo.SecondaryButton,
		CustomID: fmt.Sprintf("pkm:choose:%d", pokemonID),
	})
	navButtons = append(navButtons, discordgo.Button{
		Label:    "✅ 确认",
		Style:    discordgo.SuccessButton,
		CustomID: fmt.Sprintf("pkm:confirmmoves:%d", pokemonID),
	})
	rows = append(rows, discordgo.ActionsRow{Components: navButtons})

	// 确保不超过 Discord 的 5 行限制
	if len(rows) > 5 {
		rows = rows[:5]
	}

	if err := c.bot.RespondWithEmbed(i.Interaction, embed, rows, true); err != nil {
		// 如果响应失败，可能是交互已过期，尝试发送新消息
		log.Printf("技能配置响应失败: %v", err)
	}
}

// handleSetMove 设置/取消技能
func (c *PokemonCommands) handleSetMove(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr, moveIdxStr, pageStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)
	moveIdx, _ := strconv.Atoi(moveIdxStr)

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		config = &pokemon_app.PokemonConfig{PokemonID: pokemonID, MoveIndices: []int{}}
	}

	found := -1
	for j, idx := range config.MoveIndices {
		if idx == moveIdx {
			found = j
			break
		}
	}

	if found >= 0 {
		config.MoveIndices = append(config.MoveIndices[:found], config.MoveIndices[found+1:]...)
	} else if len(config.MoveIndices) < 4 {
		config.MoveIndices = append(config.MoveIndices, moveIdx)
	} else {
		c.bot.RespondEphemeral(i.Interaction, "❌ 最多只能选择4个技能")
		return
	}

	c.handler.SetConfig(channelID, userID, config)
	c.handleConfigMoves(i, channelID, userID, pokemonIDStr, pageStr)
}

// handleConfirmMoves 确认技能选择，返回配置面板
func (c *PokemonCommands) handleConfirmMoves(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		config = &pokemon_app.PokemonConfig{PokemonID: pokemonID, MoveIndices: []int{}}
		c.handler.SetConfig(channelID, userID, config)
	}

	c.showConfigPanel(i, channelID, userID, pokemon, config)
}

// handleShowPresets 显示预设列表
func (c *PokemonCommands) handleShowPresets(i *discordgo.InteractionCreate, userID string) {
	presets := c.handler.GetPresets(userID)

	var desc strings.Builder
	desc.WriteString("## 📋 我的配队预设\n\n")

	if len(presets) == 0 {
		desc.WriteString("_暂无保存的预设_\n\n")
		desc.WriteString("在选择宝可梦配置时点击「保存预设」来保存当前配置。")
	} else {
		for _, p := range presets {
			desc.WriteString(fmt.Sprintf("**%s** `[%s]`\n", p.Name, p.ID))
			desc.WriteString(fmt.Sprintf("  宝可梦: %s | 性格: %s\n", p.PokemonName, p.Nature))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📋 配队预设",
		Description: desc.String(),
		Color:       0x9B59B6,
	}

	var buttons []discordgo.MessageComponent
	for _, p := range presets {
		buttons = append(buttons, discordgo.Button{
			Label:    p.Name,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("pkm:loadpreset:%s", p.ID),
		})
		if len(buttons) >= 5 {
			break
		}
	}

	var components []discordgo.MessageComponent
	if len(buttons) > 0 {
		components = append(components, discordgo.ActionsRow{Components: buttons})
	}

	c.bot.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

// handleLoadPreset 加载预设
func (c *PokemonCommands) handleLoadPreset(i *discordgo.InteractionCreate, channelID, userID, presetID string) {
	if err := c.handler.LoadPresetToConfig(channelID, userID, presetID); err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	config := c.handler.GetConfig(channelID, userID)
	pokemon := c.handler.GetPokemonByID(config.PokemonID)

	c.bot.RespondEphemeral(i.Interaction, fmt.Sprintf("✅ 已加载预设！宝可梦: %s", pokemon.Name))
}

// handleSavePresetModal 显示保存预设的模态框
func (c *PokemonCommands) handleSavePresetModal(i *discordgo.InteractionCreate, channelID, userID string) {
	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 请先选择宝可梦")
		return
	}

	pokemon := c.handler.GetPokemonByID(config.PokemonID)
	defaultName := ""
	if pokemon != nil {
		defaultName = pokemon.Name
	}

	c.bot.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "pkm:savepreset_modal",
			Title:    "保存配队预设",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "preset_name",
							Label:       "预设名称",
							Style:       discordgo.TextInputShort,
							Placeholder: "输入预设名称...",
							Value:       defaultName,
							Required:    true,
							MinLength:   1,
							MaxLength:   20,
						},
					},
				},
			},
		},
	})
}

// handleDeletePreset 删除预设
func (c *PokemonCommands) handleDeletePreset(i *discordgo.InteractionCreate, userID, presetID string) {
	if err := c.handler.DeletePreset(userID, presetID); err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	c.bot.RespondEphemeral(i.Interaction, "✅ 预设已删除")
}

// handleSavePresetSubmit 处理保存预设模态框提交
func (c *PokemonCommands) handleSavePresetSubmit(i *discordgo.InteractionCreate) {
	channelID := i.ChannelID
	userID := i.Member.User.ID
	data := i.ModalSubmitData()

	// 获取预设名称
	var presetName string
	for _, row := range data.Components {
		if actionRow, ok := row.(*discordgo.ActionsRow); ok {
			for _, comp := range actionRow.Components {
				if input, ok := comp.(*discordgo.TextInput); ok && input.CustomID == "preset_name" {
					presetName = input.Value
				}
			}
		}
	}

	if presetName == "" {
		c.bot.RespondEphemeral(i.Interaction, "❌ 预设名称不能为空")
		return
	}

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 请先选择宝可梦")
		return
	}

	preset, err := c.handler.SavePreset(userID, presetName, config)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	c.bot.RespondEphemeral(i.Interaction, fmt.Sprintf("✅ 预设 **%s** 已保存！\n宝可梦: %s", preset.Name, preset.PokemonName))
}

// handleShowSwitchMenu 显示换人菜单
func (c *PokemonCommands) handleShowSwitchMenu(i *discordgo.InteractionCreate, channelID, userID string) {
	battle, err := c.handler.GetBattle(channelID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 对战不存在")
		return
	}

	player := battle.GetPlayer(userID)
	if player == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 你不在对战中")
		return
	}

	if player.Action != nil {
		c.bot.RespondEphemeral(i.Interaction, "⏳ 你已选择行动，等待对手...")
		return
	}

	// 构建可换上场的宝可梦按钮
	var buttons []discordgo.MessageComponent
	for idx, battler := range player.Team {
		if idx == player.ActiveIndex {
			continue // 跳过当前在场的宝可梦
		}
		if !battler.IsAlive() {
			continue // 跳过已倒下的宝可梦
		}
		hpPercent := battler.GetHPPercent()
		label := fmt.Sprintf("%s (%.0f%%)", battler.Pokemon.Name, hpPercent)
		buttons = append(buttons, discordgo.Button{
			Label:    label,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("pkm:doswitch:%d", idx),
		})
	}

	if len(buttons) == 0 {
		c.bot.RespondEphemeral(i.Interaction, "❌ 没有可以换上场的宝可梦")
		return
	}

	// 添加取消按钮
	buttons = append(buttons, discordgo.Button{
		Label:    "🔙 取消",
		Style:    discordgo.SecondaryButton,
		CustomID: "pkm:refresh",
	})

	var rows []discordgo.MessageComponent
	for j := 0; j < len(buttons); j += 5 {
		end := j + 5
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, discordgo.ActionsRow{Components: buttons[j:end]})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🔄 换人",
		Description: "选择要换上场的宝可梦：",
		Color:       0x3498DB,
	}

	c.bot.RespondWithEmbed(i.Interaction, embed, rows, true)
}

// handleDoSwitch 执行换人
func (c *PokemonCommands) handleDoSwitch(i *discordgo.InteractionCreate, channelID, userID, indexStr string) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无效的选择")
		return
	}

	logs, err := c.handler.SwitchPokemon(channelID, userID, index)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}

	// 检查是否为人机对战，如果是则触发 AI 行动
	battle, _ := c.handler.GetBattle(channelID)
	if battle != nil && battle.IsAIBattle && len(logs) == 0 {
		// 玩家已行动，触发 AI 行动并执行回合
		aiLogs, _ := c.handler.ExecuteAITurn(channelID)
		logs = aiLogs
		// 重新获取对战状态
		battle, _ = c.handler.GetBattle(channelID)
	}

	if len(logs) > 0 {
		// 回合执行完毕，发送战斗日志
		logText := strings.Join(logs, "\n")

		if battle != nil && battle.State == entity.BattleStateFinished {
			c.bot.RespondPublic(i.Interaction, logText)
			c.handler.EndBattle(channelID)
		} else {
			c.bot.RespondPublic(i.Interaction, logText)
			c.sendBattlePanel(i, channelID)
		}
	} else {
		// 等待对手（普通 PVP 模式）
		c.bot.RespondEphemeral(i.Interaction, "✅ 已选择换人，等待对手行动...")
	}
}

// handleForceSwitch 强制换人（宝可梦倒下时）
func (c *PokemonCommands) handleForceSwitch(i *discordgo.InteractionCreate, channelID, userID, indexStr string) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无效的选择")
		return
	}

	battle, err := c.handler.GetBattle(channelID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 对战不存在")
		return
	}

	player := battle.GetPlayer(userID)
	if player == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 你不在对战中")
		return
	}

	if index < 0 || index >= len(player.Team) {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无效的宝可梦")
		return
	}

	if !player.Team[index].IsAlive() {
		c.bot.RespondEphemeral(i.Interaction, "❌ 该宝可梦已倒下")
		return
	}

	// 直接换人，不消耗行动
	oldName := player.Team[player.ActiveIndex].Pokemon.Name
	player.ActiveIndex = index
	newName := player.Team[index].Pokemon.Name

	c.bot.RespondPublic(i.Interaction, fmt.Sprintf("🔄 **%s** 收回了 **%s**，换上了 **%s**！", player.Username, oldName, newName))
	c.sendBattlePanel(i, channelID)
}

// handleSearchMoveModal 显示技能搜索模态框
func (c *PokemonCommands) handleSearchMoveModal(i *discordgo.InteractionCreate, pokemonIDStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)

	err := c.bot.Session().InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: fmt.Sprintf("pkm:searchmove_modal:%d", pokemonID),
			Title:    "🔍 搜索技能",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "keyword",
							Label:       "输入技能名称关键字",
							Style:       discordgo.TextInputShort,
							Placeholder: "例如：十万伏特、冲浪、地震...",
							Required:    true,
							MinLength:   1,
							MaxLength:   20,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("显示技能搜索模态框失败: %v", err)
	}
}

// handleSearchMoveModalSubmit 处理技能搜索模态框提交
func (c *PokemonCommands) handleSearchMoveModalSubmit(i *discordgo.InteractionCreate, pokemonIDStr string) {
	channelID := i.ChannelID
	userID := i.Member.User.ID
	pokemonID, _ := strconv.Atoi(pokemonIDStr)

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	// 获取搜索关键字
	data := i.ModalSubmitData()
	var keyword string
	for _, row := range data.Components {
		if ar, ok := row.(*discordgo.ActionsRow); ok {
			for _, comp := range ar.Components {
				if ti, ok := comp.(*discordgo.TextInput); ok && ti.CustomID == "keyword" {
					keyword = ti.Value
				}
			}
		}
	}

	if keyword == "" {
		c.bot.RespondEphemeral(i.Interaction, "❌ 请输入搜索关键字")
		return
	}

	// 搜索匹配的技能
	var matchedMoves []*entity.Move
	keywordLower := strings.ToLower(keyword)
	for _, move := range pokemon.LearnableMoves {
		if strings.Contains(strings.ToLower(move.Name), keywordLower) {
			matchedMoves = append(matchedMoves, move)
		}
	}

	if len(matchedMoves) == 0 {
		c.bot.RespondEphemeral(i.Interaction, fmt.Sprintf("❌ 未找到包含 \"%s\" 的技能", keyword))
		return
	}

	// 限制最多显示 20 个结果
	if len(matchedMoves) > 20 {
		matchedMoves = matchedMoves[:20]
	}

	// 获取当前配置
	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		config = &pokemon_app.PokemonConfig{PokemonID: pokemonID, MoveIndices: []int{}}
	}

	// 构建搜索结果
	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("🔍 搜索 \"%s\" 的结果 (%d 个)\n\n", keyword, len(matchedMoves)))

	// 显示已选技能
	desc.WriteString("**已选技能：**")
	if len(config.MoveIndices) == 0 {
		desc.WriteString(" 无\n\n")
	} else {
		desc.WriteString("\n")
		for _, idx := range config.MoveIndices {
			if idx < len(pokemon.LearnableMoves) {
				move := pokemon.LearnableMoves[idx]
				desc.WriteString(fmt.Sprintf("• %s (%s)\n", move.Name, string(move.Type)))
			}
		}
		desc.WriteString("\n")
	}

	desc.WriteString("**搜索结果：**\n")

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🔍 %s 的技能搜索", pokemon.Name),
		Description: desc.String(),
		Color:       0x3498DB,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/%d.gif", pokemonID),
		},
	}

	// 构建技能选择按钮 (每行最多 5 个)
	var components []discordgo.MessageComponent
	var currentRow []discordgo.MessageComponent

	for _, move := range matchedMoves {
		// 找到这个技能在原始列表中的索引
		moveIdx := -1
		for idx, m := range pokemon.LearnableMoves {
			if m.Name == move.Name {
				moveIdx = idx
				break
			}
		}
		if moveIdx == -1 {
			continue
		}

		// 检查是否已选择
		isSelected := false
		for _, idx := range config.MoveIndices {
			if idx == moveIdx {
				isSelected = true
				break
			}
		}

		label := move.Name
		if isSelected {
			label = "✓ " + label
		}

		style := discordgo.SecondaryButton
		if isSelected {
			style = discordgo.SuccessButton
		}

		currentRow = append(currentRow, discordgo.Button{
			Label:    label,
			Style:    style,
			CustomID: fmt.Sprintf("pkm:selectmove:%d:%d", pokemonID, moveIdx),
			Disabled: isSelected || len(config.MoveIndices) >= 4,
		})

		if len(currentRow) == 5 {
			components = append(components, discordgo.ActionsRow{Components: currentRow})
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		components = append(components, discordgo.ActionsRow{Components: currentRow})
	}

	// 添加返回按钮
	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "⬅️ 返回技能列表",
				Style:    discordgo.SecondaryButton,
				CustomID: fmt.Sprintf("pkm:cfgmoves:%d:1", pokemonID),
			},
			discordgo.Button{
				Label:    "✅ 确认技能",
				Style:    discordgo.SuccessButton,
				CustomID: fmt.Sprintf("pkm:confirmmoves:%d", pokemonID),
				Disabled: len(config.MoveIndices) == 0,
			},
		},
	})

	c.bot.RespondWithEmbed(i.Interaction, embed, components, true)
}

// handleSelectSearchedMove 处理从搜索结果中选择技能
func (c *PokemonCommands) handleSelectSearchedMove(i *discordgo.InteractionCreate, channelID, userID, pokemonIDStr, moveIdxStr string) {
	pokemonID, _ := strconv.Atoi(pokemonIDStr)
	moveIdx, _ := strconv.Atoi(moveIdxStr)

	pokemon := c.handler.GetPokemonByID(pokemonID)
	if pokemon == nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ 未找到宝可梦")
		return
	}

	if moveIdx < 0 || moveIdx >= len(pokemon.LearnableMoves) {
		c.bot.RespondEphemeral(i.Interaction, "❌ 无效的技能")
		return
	}

	config := c.handler.GetConfig(channelID, userID)
	if config == nil {
		config = &pokemon_app.PokemonConfig{PokemonID: pokemonID, MoveIndices: []int{}}
	}

	// 检查是否已选择 4 个技能
	if len(config.MoveIndices) >= 4 {
		c.bot.RespondEphemeral(i.Interaction, "❌ 已选择 4 个技能，请先取消一个")
		return
	}

	// 检查是否已经选择了这个技能
	for _, idx := range config.MoveIndices {
		if idx == moveIdx {
			c.bot.RespondEphemeral(i.Interaction, "❌ 该技能已被选择")
			return
		}
	}

	// 添加技能
	config.MoveIndices = append(config.MoveIndices, moveIdx)
	c.handler.SetConfig(channelID, userID, config)

	// 返回技能配置页面
	c.handleConfigMoves(i, channelID, userID, pokemonIDStr, "1")
}
