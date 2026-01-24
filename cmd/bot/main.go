package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	pokemonapp "github.com/user/dcminigames/internal/application/pokemon"
	unoapp "github.com/user/dcminigames/internal/application/uno"
	"github.com/user/dcminigames/internal/infrastructure/activity"
	"github.com/user/dcminigames/internal/infrastructure/discord"
	"github.com/user/dcminigames/internal/infrastructure/imaging"
	"github.com/user/dcminigames/internal/infrastructure/persistence/memory"
	"github.com/user/dcminigames/internal/infrastructure/pokeapi"
	"github.com/user/dcminigames/internal/interfaces/discord/commands"
	"github.com/user/dcminigames/pkg/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if cfg.Discord.Token == "" {
		log.Fatal("Discord Token 未配置")
	}

	bot, err := discord.NewBot(cfg.Discord.Token)
	if err != nil {
		log.Fatalf("创建 Bot 失败: %v", err)
	}

	// 初始化依赖
	gameRepo := memory.NewGameRepository()
	cardRenderer := imaging.NewCardRenderer(cfg.Uno.AssetsPath)
	unoHandler := unoapp.NewHandler(gameRepo, cardRenderer)
	unoCommands := commands.NewUnoCommands(bot, unoHandler)

	// 预加载宝可梦数据（避免首次使用时超时）
	log.Println("正在预加载宝可梦数据...")
	if err := pokeapi.EnsureDataLoaded(); err != nil {
		log.Printf("预加载宝可梦数据失败: %v", err)
	} else {
		log.Printf("宝可梦数据加载完成，共 %d 只宝可梦", pokeapi.GetTotalPokemonCount())
	}

	// 初始化宝可梦对战
	battleRepo := memory.NewBattleRepository()
	pokemonHandler := pokemonapp.NewHandler(battleRepo)
	pokemonCommands := commands.NewPokemonCommands(bot, pokemonHandler)

	// 初始化 Activity 服务（无名杀/三国杀）
	var activityServer *activity.Server
	var activityCommands *commands.ActivityCommands
	if cfg.Activity.Enabled {
		log.Println("正在启动 Activity 服务...")
		activityServer = activity.NewServer(cfg.Activity)
		if err := activityServer.Start(); err != nil {
			log.Printf("Activity 服务启动失败: %v", err)
		} else {
			activityCommands = commands.NewActivityCommands(bot, activityServer)
			log.Printf("Activity 服务已启动: %s", activityServer.GetPublicURL())
		}
	}

	// 注册事件处理器
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot 已登录: %s#%s", s.State.User.Username, s.State.User.Discriminator)
	})
	bot.AddHandler(unoCommands.HandleInteraction)
	bot.AddHandler(pokemonCommands.HandleInteraction)
	if activityCommands != nil {
		bot.AddHandler(activityCommands.HandleInteraction)
	}

	// 启动 Bot
	if err := bot.Start(); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
	defer bot.Stop()

	// 同步斜杠命令（启动时重置并重新注册）
	allCommands := append(unoCommands.Commands(), pokemonCommands.Commands()...)
	if activityCommands != nil {
		allCommands = append(allCommands, activityCommands.Commands()...)
	}
	if err := bot.SyncCommands(cfg.Discord.GuildID, allCommands); err != nil {
		log.Printf("同步命令失败: %v", err)
	}

	log.Println("Bot 运行中，按 Ctrl+C 退出")

	// 等待退出信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	log.Println("正在关闭...")

	// 关闭 Activity 服务
	if activityServer != nil {
		if err := activityServer.Stop(); err != nil {
			log.Printf("关闭 Activity 服务失败: %v", err)
		}
	}
}