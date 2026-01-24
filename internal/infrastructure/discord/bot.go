package discord

import (
	"bytes"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot(token string) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("创建会话失败: %w", err)
	}
	s.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildVoiceStates
	return &Bot{session: s}, nil
}

func (b *Bot) Session() *discordgo.Session {
	return b.session
}

func (b *Bot) AddHandler(handler interface{}) {
	b.session.AddHandler(handler)
}

// SyncCommands 同步命令：先清理旧命令，再注册新命令
func (b *Bot) SyncCommands(guildID string, commands []*discordgo.ApplicationCommand) error {
	log.Println("正在同步斜杠命令...")

	existing, err := b.session.ApplicationCommands(b.session.State.User.ID, guildID)
	if err != nil {
		log.Printf("获取现有命令失败: %v", err)
	} else {
		for _, cmd := range existing {
			err := b.session.ApplicationCommandDelete(b.session.State.User.ID, guildID, cmd.ID)
			if err != nil {
				log.Printf("删除命令 %s 失败: %v", cmd.Name, err)
			} else {
				log.Printf("已删除旧命令: %s", cmd.Name)
			}
		}
	}

	for _, cmd := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, guildID, cmd)
		if err != nil {
			return fmt.Errorf("注册命令 %s 失败: %w", cmd.Name, err)
		}
		log.Printf("已注册命令: %s", cmd.Name)
	}

	log.Println("命令同步完成")
	return nil
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	log.Println("Bot 已上线")
	return nil
}

func (b *Bot) Stop() error {
	return b.session.Close()
}

func (b *Bot) RespondEphemeral(i *discordgo.Interaction, content string) error {
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (b *Bot) RespondPublic(i *discordgo.Interaction, content string) error {
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func (b *Bot) RespondWithComponents(i *discordgo.Interaction, content string, components []discordgo.MessageComponent, ephemeral bool) error {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: components,
			Flags:      flags,
		},
	})
}

func (b *Bot) RespondWithEmbed(i *discordgo.Interaction, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent, ephemeral bool) error {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
			Flags:      flags,
		},
	})
}

func (b *Bot) RespondWithFile(i *discordgo.Interaction, content, fileName string, data []byte, components []discordgo.MessageComponent, ephemeral bool) error {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Files: []*discordgo.File{
				{Name: fileName, Reader: bytes.NewReader(data)},
			},
			Components: components,
			Flags:      flags,
		},
	})
}

// RespondPublicWithFile 公屏发送带图片的消息
func (b *Bot) RespondPublicWithFile(i *discordgo.Interaction, content, fileName string, data []byte) error {
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Files: []*discordgo.File{
				{Name: fileName, Reader: bytes.NewReader(data)},
			},
		},
	})
}

// FollowUpEphemeralWithFile 发送仅用户可见的后续消息（带图片和按钮）
func (b *Bot) FollowUpEphemeralWithFile(i *discordgo.Interaction, content, fileName string, data []byte, components []discordgo.MessageComponent) error {
	_, err := b.session.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: content,
		Files: []*discordgo.File{
			{Name: fileName, Reader: bytes.NewReader(data)},
		},
		Components: components,
		Flags:      discordgo.MessageFlagsEphemeral,
	})
	return err
}

// FollowUpPublicWithEmbed 发送公开的后续消息（带 Embed 和按钮）
func (b *Bot) FollowUpPublicWithEmbed(i *discordgo.Interaction, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) error {
	_, err := b.session.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	return err
}

// SendChannelEmbed 直接向频道发送 Embed 消息（不依赖交互）
func (b *Bot) SendChannelEmbed(channelID string, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) error {
	_, err := b.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	return err
}

// RespondWithEmbedAndFile 发送带 Embed 和图片的消息（图片嵌入 Embed 中）
func (b *Bot) RespondWithEmbedAndFile(i *discordgo.Interaction, embed *discordgo.MessageEmbed, fileName string, data []byte, components []discordgo.MessageComponent, ephemeral bool) error {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "attachment://" + fileName,
	}
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Files:      []*discordgo.File{{Name: fileName, Reader: bytes.NewReader(data)}},
			Components: components,
			Flags:      flags,
		},
	})
}

// RespondPublicWithEmbedAndFile 公开发送带 Embed 和图片的消息
func (b *Bot) RespondPublicWithEmbedAndFile(i *discordgo.Interaction, embed *discordgo.MessageEmbed, fileName string, data []byte) error {
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "attachment://" + fileName,
	}
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Files:  []*discordgo.File{{Name: fileName, Reader: bytes.NewReader(data)}},
		},
	})
}

// SendChannelEmbedWithFile 直接向频道发送带图片的 Embed 消息
func (b *Bot) SendChannelEmbedWithFile(channelID string, embed *discordgo.MessageEmbed, fileName string, data []byte, components []discordgo.MessageComponent) error {
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "attachment://" + fileName,
	}
	_, err := b.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Files:      []*discordgo.File{{Name: fileName, Reader: bytes.NewReader(data)}},
		Components: components,
	})
	return err
}

func (b *Bot) UpdateMessage(i *discordgo.Interaction, content string, components []discordgo.MessageComponent) error {
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: components,
		},
	})
}

func (b *Bot) UpdateWithEmbed(i *discordgo.Interaction, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) error {
	return b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// SendDMWithFile 通过私信发送带图片和按钮的消息给指定用户
func (b *Bot) SendDMWithFile(userID, content, fileName string, data []byte, components []discordgo.MessageComponent) error {
	channel, err := b.session.UserChannelCreate(userID)
	if err != nil {
		return fmt.Errorf("创建私信频道失败: %w", err)
	}
	_, err = b.session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: content,
		Files: []*discordgo.File{
			{Name: fileName, Reader: bytes.NewReader(data)},
		},
		Components: components,
	})
	if err != nil {
		return fmt.Errorf("发送私信失败: %w", err)
	}
	return nil
}
