package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Discord  DiscordConfig  `yaml:"discord"`
	Uno      UnoConfig      `yaml:"uno"`
	LLM      LLMConfig      `yaml:"llm"`
	Activity ActivityConfig `yaml:"activity"`
}

type DiscordConfig struct {
	Token   string `yaml:"token"`
	GuildID string `yaml:"guild_id"`
}

type UnoConfig struct {
	AssetsPath string `yaml:"assets_path"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	BaseURL  string `yaml:"base_url"`
	Model    string `yaml:"model"`
}

type ActivityConfig struct {
	Enabled      bool   `yaml:"enabled"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Port         int    `yaml:"port"`
	PublicURL    string `yaml:"public_url"`
	GamePath     string `yaml:"game_path"`
	DevMode      bool   `yaml:"dev_mode"`  // 是否使用 Vite 开发服务器
	VitePort     int    `yaml:"vite_port"` // Vite 开发服务器端口
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if cfg.Uno.AssetsPath == "" {
		cfg.Uno.AssetsPath = "./assets/uno"
	}
	if cfg.LLM.Provider == "" {
		cfg.LLM.Provider = "openai"
	}
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = "gpt-4"
	}
	if cfg.Activity.Port == 0 {
		cfg.Activity.Port = 8080
	}
	if cfg.Activity.GamePath == "" {
		cfg.Activity.GamePath = "./noname"
	}
	if cfg.Activity.VitePort == 0 {
		cfg.Activity.VitePort = 5173
	}

	return &cfg, nil
}