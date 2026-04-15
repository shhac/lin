package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Workspace struct {
	APIKey string `json:"api_key"`
	Name   string `json:"name,omitempty"`
	URLKey string `json:"urlKey,omitempty"`
}

type TruncationSettings struct {
	MaxLength *int `json:"maxLength,omitempty"`
}

type PaginationSettings struct {
	DefaultPageSize *int `json:"defaultPageSize,omitempty"`
}

type Settings struct {
	Truncation *TruncationSettings `json:"truncation,omitempty"`
	Pagination *PaginationSettings `json:"pagination,omitempty"`
}

type Config struct {
	LegacyAPIKey     string               `json:"api_key,omitempty"`
	DefaultWorkspace string               `json:"default_workspace,omitempty"`
	Workspaces       map[string]Workspace `json:"workspaces,omitempty"`
	Settings         *Settings            `json:"settings,omitempty"`
}

var (
	cache       *Config
	cacheMu     sync.Mutex
	overrideDir string
)

func SetConfigDir(dir string) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	overrideDir = dir
	cache = nil
}

func ConfigDir() string {
	if overrideDir != "" {
		return overrideDir
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "lin")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "lin")
}

func configPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func Read() *Config {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if cache != nil {
		return cache
	}
	data, err := os.ReadFile(configPath())
	if err != nil {
		return defaultConfig()
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}
	if cfg.Workspaces == nil {
		cfg.Workspaces = make(map[string]Workspace)
	}
	cache = &cfg
	return cache
}

func Write(cfg *Config) error {
	cacheMu.Lock()
	cache = nil
	cacheMu.Unlock()

	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), append(data, '\n'), 0o644)
}

func ClearCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache = nil
}

func defaultConfig() *Config {
	cfg := &Config{
		Workspaces: make(map[string]Workspace),
	}
	cache = cfg
	return cfg
}
