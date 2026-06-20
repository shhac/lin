package output

import (
	"encoding/json"
	"fmt"
	"sync"

	_ "github.com/shhac/lib-agent-cli/yaml" // registers the YAML encoder (yaml.v3) for out.FormatYAML
	out "github.com/shhac/lib-agent-output"
	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/truncation"
)

// Format and its values come from the shared contract; ParseFormat is therefore
// the family's lenient parser (accepts "ndjson"/"yml", case-insensitive).
type Format = out.Format

const (
	FormatJSON   = out.FormatJSON
	FormatYAML   = out.FormatYAML
	FormatNDJSON = out.FormatNDJSON
)

// ParseFormat accepts the family's lenient set (json/yaml/jsonl plus the
// yml/ndjson aliases, case-insensitive) but preserves lin's original message
// and actionable hint on rejection, so the stderr error stays byte-identical.
func ParseFormat(s string) (Format, error) {
	f, err := out.ParseFormat(s)
	if err != nil {
		return "", out.New(fmt.Sprintf("unknown format %q, expected: json, yaml, jsonl", s), out.FixableByAgent).
			WithHint("use --format json, --format yaml, or --format jsonl")
	}
	return f, nil
}

var (
	formatMu   sync.RWMutex
	flagFormat string
)

func ConfigureFormat(format string) error {
	if format != "" {
		if _, err := ParseFormat(format); err != nil {
			return err
		}
	}
	formatMu.Lock()
	flagFormat = format
	formatMu.Unlock()
	return nil
}

// ResolveFormat returns the effective format: the --format flag, else the
// configured default, else defaultFormat. (lin keeps its one-arg, config-aware
// contract — the shared two-arg ResolveFormat doesn't read config.)
func ResolveFormat(defaultFormat Format) Format {
	formatMu.RLock()
	f := flagFormat
	formatMu.RUnlock()
	if f != "" {
		parsed, err := ParseFormat(f)
		if err == nil {
			return parsed
		}
		return defaultFormat
	}
	cfg := config.Read()
	if cfg.Settings != nil && cfg.Settings.Output != nil && cfg.Settings.Output.DefaultFormat != "" {
		parsed, err := ParseFormat(cfg.Settings.Output.DefaultFormat)
		if err == nil {
			return parsed
		}
	}
	return defaultFormat
}

// Print cleans (prune + truncate) then encodes data in the given format via the
// shared encoder. Pruning is opt-in; truncation rides along with it.
func Print(data any, format Format, prune bool) {
	cleaned, ok := toCleanAny(data, prune)
	if !ok {
		return
	}
	// Data is already cleaned, so pass a nil pruner — out.Print just encodes.
	_ = out.Print(Stdout(), cleaned, format, nil)
}

func PrintJSON(data any) {
	Print(data, ResolveFormat(FormatJSON), true)
}

func printNDJSON(data any) {
	w := NewNDJSONWriter(Stdout())
	if arr, ok := data.([]any); ok {
		for _, item := range arr {
			_ = w.WriteItem(item)
		}
		return
	}
	_ = w.WriteItem(data)
}

func toCleanAny(data any, prune bool) (any, bool) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}
	var decoded any
	if err := json.Unmarshal(b, &decoded); err != nil {
		return nil, false
	}
	if prune {
		decoded = truncation.Apply(pruneEmpty(decoded))
	}
	return decoded, true
}
