package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/user/dcminigames/internal/infrastructure/activity"
	"github.com/user/dcminigames/internal/infrastructure/discord"
)

// Entry Point Command 类型常量
const (
	ApplicationCommandTypePrimaryEntryPoint = 4 // PRIMARY_ENTRY_POINT
)

// Interaction Response 类型常量
const (
	InteractionResponseTypeLaunchActivity = 12 // LAUNCH_ACTIVITY
)

type ActivityCommands struct {
	bot    *discord.Bot
	server *activity.Server
}

func NewActivityCommands(bot *discord.Bot, server *activity.Server) *ActivityCommands {
	return &ActivityCommands{bot: bot, server: server}
}

func (c *ActivityCommands) Commands() []*discordgo.ApplicationCommand {
	// 注册 Entry Point 命令 (type: 4)
	entryPointType := discordgo.ApplicationCommandType(ApplicationCommandTypePrimaryEntryPoint)
	return []*discordgo.ApplicationCommand{
		{
			Name:        "sanguosha",
			Description: "启动三国杀 (无名杀) Activity",
			Type:        entryPointType,
		},
	}
}

func (c *ActivityCommands) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "sanguosha", "noname":
			c.launchActivity(s, i)
		}
	}
}

func (c *ActivityCommands) launchActivity(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 对于 Entry Point 命令，响应类型必须是 LAUNCH_ACTIVITY (12)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(InteractionResponseTypeLaunchActivity),
	})

	if err != nil {
		log.Printf("启动 Activity 失败: %v", err)
		// 如果 LAUNCH_ACTIVITY 失败，尝试发送错误消息
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ 启动 Activity 失败，请确保应用已正确配置。",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	log.Printf("Activity 已启动，用户: %s", i.Member.User.Username)
}