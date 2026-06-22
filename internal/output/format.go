package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	_ "github.com/shhac/lib-agent-cli/yaml" // registers the YAML encoder (yaml.v3) for out.FormatYAML
	out "github.com/shhac/lib-agent-output"
	"github.com/shhac/lin/internal/truncation"
)

// Format and its values come from the shared contract; ParseFormat is therefore
// the family's lenient parser (accepts "ndjson"/"yml", case-insensitive).
type Format = out.Format

const (
	FormatJSON   = out.FormatJSON
	FormatYAML   = out.FormatYAML
	FormatNDJSON = out.FormatNDJSON
	// FormatPretty is lin's human-readable card renderer. It is NOT part of the
	// shared/universal enum (ParseFormat rejects it): it's opt-in per command via
	// AllowPretty and handled by the command's own render branch, never reaching
	// the universal encoder. Mirrors the agent-slack "transcript" pattern.
	FormatPretty Format = "pretty"
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

// ConfigureFormat records the resolved --format value (the flag, or the config
// default applied in the root's ConfigDefaults) so ResolveFormat can read it at
// render time. Validation is libcli.NewRoot's job now — the universal set plus
// any per-command AllowFormats("pretty") — so this no longer validates.
func ConfigureFormat(format string) {
	formatMu.Lock()
	flagFormat = format
	formatMu.Unlock()
}

// ResolveFormat returns the effective format: the recorded --format value (the
// flag, or the config default applied in ConfigDefaults), else defaultFormat.
// "pretty" is mapped here; a command only reaches a pretty value if it opted in
// via AllowPretty, so NewRoot's validator has already accepted it.
func ResolveFormat(defaultFormat Format) Format {
	formatMu.RLock()
	f := flagFormat
	formatMu.RUnlock()
	if f == "" {
		return defaultFormat
	}
	if strings.EqualFold(f, string(FormatPretty)) {
		return FormatPretty
	}
	if parsed, err := ParseFormat(f); err == nil {
		return parsed
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
