package config

import (
	"path/filepath"
	"sync"

	"github.com/shhac/lib-agent-cli/creds"
	"github.com/shhac/lib-agent-cli/xdg"
)

type Workspace struct {
	APIKey string `json:"api_key"`
	Name   string `json:"name,omitempty"`
	URLKey string `json:"urlKey,omitempty"`
	// OrgID is the Linear organization id this workspace's key resolves to. It
	// is the strongest identity anchor for MCP enrollment convergence (urlKey
	// can be renamed); recorded when known, empty for legacy logins.
	OrgID string `json:"orgId,omitempty"`
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

// store is lin's config file: 0600 in a 0700 parent, atomic replacement, and
// the lock that update() serializes on. This used to be hand-rolled with
// os.ReadFile/os.WriteFile, which carried a lost-update race — two concurrent
// processes (e.g. `lin auth login` racing `lin config set`) could each build
// their write from a snapshot taken before the other landed, and the loser's
// change was silently erased. Because this file is also where a workspace's
// API key lives when the keychain is unavailable, a lost workspace entry can
// strand a live secret in the keychain that nothing references anymore.
func store() creds.Store {
	return creds.Store{Path: configPath()}
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
	var cfg Config
	if err := store().Load(&cfg); err != nil {
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

	return store().Save(cfg)
}

// update performs a read-modify-write against the config file under ONE
// exclusive lock spanning load, mutate, and save — so two concurrent
// processes serialize instead of each building its write from a snapshot
// taken before the other landed (the lost-update race Write alone cannot
// avoid).
//
// It loads FRESH state from disk rather than the package-level cache, which
// may be stale relative to another process's write; using the cache here
// would silently reintroduce the same lost update, just scoped in-process
// instead of cross-process. The cache is invalidated afterward so the next
// Read() in this process picks up what was actually written.
func update(mutate func(cfg *Config) error) error {
	var cfg Config
	err := store().Update(&cfg, func() error {
		if cfg.Workspaces == nil {
			cfg.Workspaces = make(map[string]Workspace)
		}
		return mutate(&cfg)
	})

	cacheMu.Lock()
	cache = nil
	cacheMu.Unlock()

	return err
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
