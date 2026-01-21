package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	unoapp "github.com/user/dcminigames/internal/application/uno"
	"github.com/user/dcminigames/internal/domain/uno/entity"
	"github.com/user/dcminigames/internal/domain/uno/valueobject"
	"github.com/user/dcminigames/internal/infrastructure/discord"
)

type UnoCommands struct {
	bot     *discord.Bot
	handler *unoapp.Handler
}

func NewUnoCommands(bot *discord.Bot, handler *unoapp.Handler) *UnoCommands {
	return &UnoCommands{bot: bot, handler: handler}
}

func (c *UnoCommands) Commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "uno",
			Description: "打开 UNO 游戏面板",
		},
	}
}

func (c *UnoCommands) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		if data.Name == "uno" {
			c.showPanel(i)
		}
	} else if i.Type == discordgo.InteractionMessageComponent {
		c.handleComponent(i)
	}
}

func (c *UnoCommands) showPanel(i *discordgo.InteractionCreate) {
	channelID := i.ChannelID
	userID := i.Member.User.ID
	game, err := c.handler.GetGame(channelID)
	var embed *discordgo.MessageEmbed
	var components []discordgo.MessageComponent
	if err != nil {
		embed = &discordgo.MessageEmbed{
			Title:       "🎴 UNO 游戏",
			Description: "当前没有进行中的游戏\n点击下方按钮创建新游戏",
			Color:       0x00D166,
		}
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "🎮 创建游戏", Style: discordgo.SuccessButton, CustomID: "uno:create"},
				},
			},
		}
	} else {
		embed, components = c.buildGamePanel(game, userID)
	}
	c.bot.RespondWithEmbed(i.Interaction, embed, components, true)
}

func (c *UnoCommands) buildGamePanel(game *entity.Game, userID string) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	var embed *discordgo.MessageEmbed
	var components []discordgo.MessageComponent
	player := game.GetPlayer(userID)
	isInGame := player != nil
	isHost := len(game.Players) > 0 && game.Players[0].ID == userID
	switch game.State {
	case entity.GameStateWaiting:
		var playerList []string
		for _, p := range game.Players {
			playerList = append(playerList, p.Username)
		}
		embed = &discordgo.MessageEmbed{
			Title:       "🎴 UNO - 等待玩家",
			Description: fmt.Sprintf("游戏ID: `%s`\n\n**已加入玩家 (%d/10):**\n%s", game.ID[:8], len(game.Players), strings.Join(playerList, "\n")),
			Color:       0xFEE75C,
		}
		var buttons []discordgo.MessageComponent
		if !isInGame {
			buttons = append(buttons, discordgo.Button{Label: "✋ 加入游戏", Style: discordgo.SuccessButton, CustomID: "uno:join"})
		}
		if isHost && len(game.Players) >= 2 {
			buttons = append(buttons, discordgo.Button{Label: "🚀 开始游戏", Style: discordgo.PrimaryButton, CustomID: "uno:start"})
		}
		buttons = append(buttons, discordgo.Button{Label: "🔄 刷新", Style: discordgo.SecondaryButton, CustomID: "uno:refresh"})
		if isHost {
			buttons = append(buttons, discordgo.Button{Label: "❌ 解散", Style: discordgo.DangerButton, CustomID: "uno:end"})
		}
		components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: buttons}}
	case entity.GameStatePlaying:
		current := game.GetCurrentPlayer()
		topCard := game.GetTopCard()
		var handInfo string
		for _, p := range game.Players {
			marker := ""
			if p.ID == current.ID {
				marker = " 👈"
			}
			handInfo += fmt.Sprintf("%s: %d张%s\n", p.Username, p.HandSize(), marker)
		}
		embed = &discordgo.MessageEmbed{
			Title: "🎴 UNO - 游戏中",
			Fields: []*discordgo.MessageEmbedField{
				{Name: "当前牌", Value: topCard.String(), Inline: true},
				{Name: "当前颜色", Value: string(game.CurrentColor), Inline: true},
				{Name: "当前玩家", Value: current.Username, Inline: true},
				{Name: "玩家手牌", Value: handInfo, Inline: false},
			},
			Color: c.getColorCode(game.CurrentColor),
		}
		var buttons []discordgo.MessageComponent
		if isInGame {
			buttons = append(buttons, discordgo.Button{Label: "🃏 查看手牌", Style: discordgo.PrimaryButton, CustomID: "uno:hand"})
		}
		buttons = append(buttons, discordgo.Button{Label: "🔄 刷新", Style: discordgo.SecondaryButton, CustomID: "uno:refresh"})
		if isHost {
			buttons = append(buttons, discordgo.Button{Label: "❌ 结束", Style: discordgo.DangerButton, CustomID: "uno:end"})
		}
		components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: buttons}}
	case entity.GameStateFinished:
		embed = &discordgo.MessageEmbed{
			Title:       "🎉 游戏结束",
			Description: fmt.Sprintf("**%s** 获胜！", game.Winner.Username),
			Color:       0x00D166,
		}
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "🎮 新游戏", Style: discordgo.SuccessButton, CustomID: "uno:create"},
				},
			},
		}
	}
	return embed, components
}

func (c *UnoCommands) getColorCode(color valueobject.Color) int {
	switch color {
	case valueobject.ColorRed:
		return 0xED4245
	case valueobject.ColorBlue:
		return 0x5865F2
	case valueobject.ColorGreen:
		return 0x57F287
	case valueobject.ColorYellow:
		return 0xFEE75C
	default:
		return 0x99AAB5
	}
}

func (c *UnoCommands) handleComponent(i *discordgo.InteractionCreate) {
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
	switch prefix {
	case "uno":
		c.handleUnoAction(i, action, channelID, userID, username)
	case "play":
		c.handlePlayCard(i, action, channelID, userID)
	case "color":
		c.handleColorSelect(i, action, channelID, userID, parts)
	case "draw":
		c.handleDraw(i, channelID, userID)
	case "pass":
		c.handlePass(i, channelID, userID)
	}
}

func (c *UnoCommands) handleUnoAction(i *discordgo.InteractionCreate, action, channelID, userID, username string) {
	switch action {
	case "create":
		_, err := c.handler.CreateGame(channelID)
		if err != nil {
			c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
			return
		}
		c.handler.JoinGame(channelID, userID, username)
		c.bot.RespondPublic(i.Interaction, fmt.Sprintf("🎴 **%s** 创建了 UNO 游戏！\n使用 `/uno` 打开面板加入游戏", username))
	case "join":
		if err := c.handler.JoinGame(channelID, userID, username); err != nil {
			c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
			return
		}
		game, _ := c.handler.GetGame(channelID)
		c.bot.RespondPublic(i.Interaction, fmt.Sprintf("✅ **%s** 加入了游戏！当前 %d 人", username, len(game.Players)))
	case "start":
		if err := c.handler.StartGame(channelID, userID); err != nil {
			c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
			return
		}
		game, _ := c.handler.GetGame(channelID)
		topCard := game.GetTopCard()
		cardImg, err := c.handler.RenderSingleCard(topCard)
		if err != nil {
			c.bot.RespondPublic(i.Interaction, c.formatGameStart(game))
		} else {
			embed := &discordgo.MessageEmbed{
				Title:       "🎮 游戏开始！",
				Description: c.formatGameStart(game),
				Color:       c.getColorCode(game.CurrentColor),
			}
			c.bot.RespondPublicWithEmbedAndFile(i.Interaction, embed, "card.jpg", cardImg)
		}
		c.sendGamePanel(i, channelID)
	case "hand":
		imgData, err := c.handler.RenderPlayerHand(channelID, userID)
		if err != nil {
			c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
			return
		}
		game, _ := c.handler.GetGame(channelID)
		player := game.GetPlayer(userID)
		components := c.buildHandComponents(player, game)
		embed := &discordgo.MessageEmbed{
			Title: "🃏 你的手牌",
			Color: c.getColorCode(game.CurrentColor),
		}
		c.bot.RespondWithEmbedAndFile(i.Interaction, embed, "hand.jpg", imgData, components, true)
	case "refresh":
		game, err := c.handler.GetGame(channelID)
		if err != nil {
			c.bot.RespondEphemeral(i.Interaction, "❌ 没有进行中的游戏")
			return
		}
		embed, components := c.buildGamePanel(game, userID)
		c.bot.UpdateWithEmbed(i.Interaction, embed, components)
	case "end":
		if err := c.handler.EndGame(channelID); err != nil {
			c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
			return
		}
		c.bot.RespondPublic(i.Interaction, "🛑 游戏已结束")
	}
}

func (c *UnoCommands) handlePlayCard(i *discordgo.InteractionCreate, action, channelID, userID string) {
	index, _ := strconv.Atoi(action)
	cards, _ := c.handler.GetPlayerHand(channelID, userID)
	if index >= 0 && index < len(cards) && cards[index].Type.IsWildCard() {
		c.bot.RespondWithComponents(i.Interaction, "选择颜色:", c.buildColorPicker(index), true)
		return
	}
	playedCard, err := c.handler.PlayCardAndGetCard(channelID, userID, index, "")
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	c.announcePlayWithCard(i, channelID, playedCard, i.Member.User.Username)
}

func (c *UnoCommands) handleColorSelect(i *discordgo.InteractionCreate, action, channelID, userID string, parts []string) {
	index, _ := strconv.Atoi(action)
	color := valueobject.Color(parts[2])
	playedCard, err := c.handler.PlayCardAndGetCard(channelID, userID, index, color)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	c.announcePlayWithCard(i, channelID, playedCard, i.Member.User.Username)
}

func (c *UnoCommands) handleDraw(i *discordgo.InteractionCreate, channelID, userID string) {
	card, err := c.handler.DrawCard(channelID, userID)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	cardImg, err := c.handler.RenderSingleCard(card)
	if err != nil {
		c.bot.RespondEphemeral(i.Interaction, fmt.Sprintf("📥 你摸了一张: %s", card.String()))
		return
	}
	embed := &discordgo.MessageEmbed{
		Title: "📥 你摸了一张牌",
		Color: c.getColorCode(card.Color),
	}
	c.bot.RespondWithEmbedAndFile(i.Interaction, embed, "card.jpg", cardImg, nil, true)
}

func (c *UnoCommands) handlePass(i *discordgo.InteractionCreate, channelID, userID string) {
	if err := c.handler.PassTurn(channelID, userID); err != nil {
		c.bot.RespondEphemeral(i.Interaction, "❌ "+err.Error())
		return
	}
	game, _ := c.handler.GetGame(channelID)
	nextPlayer := game.GetCurrentPlayer()
	c.bot.RespondPublic(i.Interaction, fmt.Sprintf("⏭️ **%s** 跳过回合，轮到 <@%s>", i.Member.User.Username, nextPlayer.ID))
	c.sendGamePanel(i, channelID)
}

func (c *UnoCommands) announcePlayWithCard(i *discordgo.InteractionCreate, channelID string, playedCard *entity.Card, username string) {
	game, _ := c.handler.GetGame(channelID)
	if game.State == entity.GameStateFinished {
		c.bot.RespondPublic(i.Interaction, fmt.Sprintf("🎉 **%s** 打出 **%s** 获胜！游戏结束！", game.Winner.Username, playedCard.String()))
		c.handler.EndGame(channelID)
		return
	}
	nextPlayer := game.GetCurrentPlayer()
	c.bot.RespondPublic(i.Interaction, fmt.Sprintf("🎴 **%s** 打出了 **%s**\n当前颜色: %s\n轮到 <@%s>",
		username, playedCard.String(), game.CurrentColor, nextPlayer.ID))
	c.sendGamePanel(i, channelID)
}

func (c *UnoCommands) sendGamePanel(i *discordgo.InteractionCreate, channelID string) {
	game, err := c.handler.GetGame(channelID)
	if err != nil || game.State != entity.GameStatePlaying {
		return
	}
	currentPlayer := game.GetCurrentPlayer()
	topCard := game.GetTopCard()
	
	// 构建玩家手牌信息
	var handInfo string
	for _, p := range game.Players {
		marker := ""
		if p.ID == currentPlayer.ID {
			marker = " 👈"
		}
		handInfo += fmt.Sprintf("%s: %d张%s\n", p.Username, p.HandSize(), marker)
	}
	
	embed := &discordgo.MessageEmbed{
		Title: "🎴 UNO - 游戏中",
		Description: fmt.Sprintf("轮到 <@%s> 出牌！", currentPlayer.ID),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "当前牌", Value: topCard.String(), Inline: true},
			{Name: "当前颜色", Value: string(game.CurrentColor), Inline: true},
			{Name: "当前玩家", Value: currentPlayer.Username, Inline: true},
			{Name: "玩家手牌", Value: handInfo, Inline: false},
		},
		Color: c.getColorCode(game.CurrentColor),
	}
	
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "🃏 查看手牌", Style: discordgo.PrimaryButton, CustomID: "uno:hand"},
				discordgo.Button{Label: "🔄 刷新", Style: discordgo.SecondaryButton, CustomID: "uno:refresh"},
			},
		},
	}
	
	// 渲染当前牌图片并嵌入 Embed
	cardImg, imgErr := c.handler.RenderSingleCard(topCard)
	if imgErr != nil {
		err = c.bot.SendChannelEmbed(channelID, embed, components)
	} else {
		err = c.bot.SendChannelEmbedWithFile(channelID, embed, "card.jpg", cardImg, components)
	}
	if err != nil {
		log.Printf("发送游戏面板失败: %v", err)
	}
}

func (c *UnoCommands) buildHandComponents(player *entity.Player, game *entity.Game) []discordgo.MessageComponent {
	var buttons []discordgo.MessageComponent
	isMyTurn := game.GetCurrentPlayer().ID == player.ID
	for idx, card := range player.Hand {
		if idx >= 20 {
			break
		}
		label := card.String()
		if len(label) > 10 {
			label = label[:10]
		}
		canPlay := isMyTurn && card.CanPlayOn(game.GetTopCard(), game.CurrentColor)
		style := discordgo.SecondaryButton
		if canPlay {
			style = discordgo.PrimaryButton
		}
		buttons = append(buttons, discordgo.Button{
			Label:    label,
			Style:    style,
			CustomID: fmt.Sprintf("play:%d", idx),
			Disabled: !canPlay,
		})
	}
	rows := make([]discordgo.MessageComponent, 0)
	for i := 0; i < len(buttons); i += 5 {
		end := i + 5
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, discordgo.ActionsRow{Components: buttons[i:end]})
	}
	if isMyTurn {
		rows = append(rows, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "📥 摸牌", Style: discordgo.SuccessButton, CustomID: "draw:"},
				discordgo.Button{Label: "⏭️ 跳过", Style: discordgo.DangerButton, CustomID: "pass:"},
			},
		})
	}
	return rows
}

func (c *UnoCommands) buildColorPicker(cardIndex int) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "🔴 红", Style: discordgo.DangerButton, CustomID: fmt.Sprintf("color:%d:Red", cardIndex)},
				discordgo.Button{Label: "🔵 蓝", Style: discordgo.PrimaryButton, CustomID: fmt.Sprintf("color:%d:Blue", cardIndex)},
				discordgo.Button{Label: "🟢 绿", Style: discordgo.SuccessButton, CustomID: fmt.Sprintf("color:%d:Green", cardIndex)},
				discordgo.Button{Label: "🟡 黄", Style: discordgo.SecondaryButton, CustomID: fmt.Sprintf("color:%d:Yellow", cardIndex)},
			},
		},
	}
}

func (c *UnoCommands) formatGameStart(game *entity.Game) string {
	var players []string
	for _, p := range game.Players {
		players = append(players, p.Username)
	}
	topCard := game.GetTopCard()
	currentPlayer := game.GetCurrentPlayer()
	return fmt.Sprintf("🎮 **游戏开始！**\n\n玩家: %s\n起始牌: **%s**\n当前颜色: %s\n\n轮到 <@%s> 出牌",
		strings.Join(players, ", "), topCard.String(), game.CurrentColor, currentPlayer.ID)
}