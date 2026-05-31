package output

import (
	"encoding/json"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/shhac/lin/internal/config"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/truncation"
)

type Format string

const (
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
	FormatNDJSON Format = "jsonl"
)

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

func ParseFormat(s string) (Format, error) {
	switch s {
	case "json":
		return FormatJSON, nil
	case "yaml":
		return FormatYAML, nil
	case "jsonl", "ndjson":
		return FormatNDJSON, nil
	default:
		return "", apierrors.Newf(apierrors.FixableByAgent, "unknown format %q, expected: json, yaml, jsonl", s).
			WithHint("use --format json, --format yaml, or --format jsonl")
	}
}

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

func Print(data any, format Format, prune bool) {
	cleaned, ok := toCleanAny(data, prune)
	if !ok {
		return
	}
	switch format {
	case FormatYAML:
		enc := yaml.NewEncoder(Stdout())
		_ = enc.Encode(cleaned)
	case FormatNDJSON:
		printNDJSON(cleaned)
	default:
		enc := json.NewEncoder(Stdout())
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		_ = enc.Encode(cleaned)
	}
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
