package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Discord DiscordConfig `yaml:"discord"`
	Uno     UnoConfig     `yaml:"uno"`
	LLM     LLMConfig     `yaml:"llm"`
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

	return &cfg, nil
}