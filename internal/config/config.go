package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/shhac/lib-agent-cli/xdg"
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

type OutputSettings struct {
	DefaultFormat string `json:"defaultFormat,omitempty"`
}

type RequestSettings struct {
	TimeoutMS *int `json:"timeoutMS,omitempty"`
}

type Settings struct {
	Truncation *TruncationSettings `json:"truncation,omitempty"`
	Pagination *PaginationSettings `json:"pagination,omitempty"`
	Output     *OutputSettings     `json:"output,omitempty"`
	Request    *RequestSettings    `json:"request,omitempty"`
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
	return xdg.ConfigDir("lin")
}

func configPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

// CacheDir is lin's cache directory — regenerable data such as downloaded
// files, kept separate from ConfigDir (which holds credentials/config). It is
// exposed to MCP clients as the read-only "cache" file root.
func CacheDir() string {
	return xdg.CacheDir("lin")
}

// DownloadsDir is where `file download` writes by default (a subdir of
// CacheDir), so a downloaded file is fetchable over MCP via the fs tool's
// "cache" root as downloads/<name>.
func DownloadsDir() string {
	return filepath.Join(CacheDir(), "downloads")
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
