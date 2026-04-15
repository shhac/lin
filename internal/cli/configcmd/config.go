package configcmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

type settingDef struct {
	section     string
	field       string
	parse       func(string) (any, error)
	description string
}

var settingDefs = map[string]settingDef{
	"truncation.maxLength": {
		section: "truncation",
		field:   "maxLength",
		parse: func(v string) (any, error) {
			n, err := strconv.Atoi(v)
			if err != nil || n < 0 {
				return nil, fmt.Errorf("Invalid value: %s. Must be a non-negative integer.", v)
			}
			return n, nil
		},
		description: "Max characters before truncating description/body/content fields (default: 200)",
	},
	"pagination.defaultPageSize": {
		section: "pagination",
		field:   "defaultPageSize",
		parse: func(v string) (any, error) {
			n, err := strconv.Atoi(v)
			if err != nil || n < 1 || n > 250 {
				return nil, fmt.Errorf("Invalid value: %s. Must be an integer between 1 and 250.", v)
			}
			return n, nil
		},
		description: "Default number of results for list/search commands (default: 50)",
	},
}

var validKeys = func() []string {
	keys := make([]string, 0, len(settingDefs))
	for k := range settingDefs {
		keys = append(keys, k)
	}
	return keys
}()

func validKeysStr() string {
	return strings.Join(validKeys, ", ")
}

func getNestedValue(settings *config.Settings, key string) any {
	def, ok := settingDefs[key]
	if !ok {
		return nil
	}
	switch def.section {
	case "truncation":
		if settings.Truncation == nil {
			return nil
		}
		if def.field == "maxLength" {
			if settings.Truncation.MaxLength == nil {
				return nil
			}
			return *settings.Truncation.MaxLength
		}
	case "pagination":
		if settings.Pagination == nil {
			return nil
		}
		if def.field == "defaultPageSize" {
			if settings.Pagination.DefaultPageSize == nil {
				return nil
			}
			return *settings.Pagination.DefaultPageSize
		}
	}
	return nil
}

// Register adds the config command group to the parent command.
func Register(parent *cobra.Command) {
	cfg := &cobra.Command{
		Use:   "config",
		Short: "View and update CLI settings",
	}
	output.HandleUnknownCommand(cfg, "Run 'lin config usage' for help")

	registerGet(cfg)
	registerSet(cfg)
	registerReset(cfg)
	registerListKeys(cfg)
	registerUsage(cfg)

	parent.AddCommand(cfg)
}

func registerGet(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Show current settings",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			settings := config.GetSettings()

			if len(args) == 0 {
				output.PrintJSON(settings)
				return
			}

			key := args[0]
			if _, ok := settingDefs[key]; !ok {
				output.PrintError(fmt.Sprintf("Unknown setting: %s. Valid keys: %s", key, validKeysStr()))
				return
			}

			value := getNestedValue(settings, key)
			output.PrintJSON(map[string]any{key: value})
		},
	}
	cfg.AddCommand(cmd)
}

func registerSet(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update a setting",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			key, value := args[0], args[1]
			def, ok := settingDefs[key]
			if !ok {
				output.PrintError(fmt.Sprintf("Unknown setting: %s. Valid keys: %s", key, validKeysStr()))
				return
			}

			parsed, err := def.parse(value)
			if err != nil {
				output.PrintError(err.Error())
				return
			}

			intVal := parsed.(int)
			settings := config.GetSettings()
			switch def.section {
			case "truncation":
				if settings.Truncation == nil {
					settings.Truncation = &config.TruncationSettings{}
				}
				settings.Truncation.MaxLength = &intVal
			case "pagination":
				if settings.Pagination == nil {
					settings.Pagination = &config.PaginationSettings{}
				}
				settings.Pagination.DefaultPageSize = &intVal
			}

			if err := config.UpdateSettings(settings); err != nil {
				output.PrintError(err.Error())
				return
			}
			output.PrintJSON(map[string]any{key: parsed})
		},
	}
	cfg.AddCommand(cmd)
}

func registerReset(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "reset [key]",
		Short: "Reset settings to defaults",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := config.ResetSettings(); err != nil {
					output.PrintError(err.Error())
					return
				}
				output.PrintJSON(map[string]any{"reset": "all"})
				return
			}

			key := args[0]
			def, ok := settingDefs[key]
			if !ok {
				output.PrintError(fmt.Sprintf("Unknown setting: %s. Valid keys: %s", key, validKeysStr()))
				return
			}

			settings := config.GetSettings()
			switch def.section {
			case "truncation":
				if settings.Truncation != nil && def.field == "maxLength" {
					settings.Truncation.MaxLength = nil
				}
			case "pagination":
				if settings.Pagination != nil && def.field == "defaultPageSize" {
					settings.Pagination.DefaultPageSize = nil
				}
			}

			if err := config.UpdateSettings(settings); err != nil {
				output.PrintError(err.Error())
				return
			}
			output.PrintJSON(map[string]any{"reset": key})
		},
	}
	cfg.AddCommand(cmd)
}

func registerListKeys(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list-keys",
		Short: "List all available setting keys",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			keys := make([]map[string]string, 0, len(settingDefs))
			for key, def := range settingDefs {
				keys = append(keys, map[string]string{
					"key":         key,
					"description": def.description,
				})
			}
			output.PrintJSON(map[string]any{"keys": keys})
		},
	}
	cfg.AddCommand(cmd)
}

func registerUsage(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed config command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(strings.TrimSpace(configUsageText))
		},
	}
	cfg.AddCommand(cmd)
}

const configUsageText = `lin config — View and update persistent CLI settings

SUBCOMMANDS:
  config get [key]            Show current settings (omit key for all)
  config set <key> <value>    Update a setting
  config reset [key]          Reset to defaults (omit key to reset all)
  config list-keys            List all available setting keys with descriptions

SETTING KEYS:
  truncation.maxLength         Max chars before truncating description/body/content fields
                               Default: 200. Must be a non-negative integer (0 = no truncation).
  pagination.defaultPageSize   Default number of results for list/search commands
                               Default: 50. Must be an integer between 1 and 250.

EXAMPLES:
  config set truncation.maxLength 500       Show more content before truncating
  config set pagination.defaultPageSize 20  Fetch fewer results per page
  config get truncation.maxLength           Check current truncation setting
  config reset truncation.maxLength         Reset truncation to default (200)
  config reset                              Reset all settings to defaults

STORAGE: Settings persisted in ~/.config/lin/config.json alongside auth credentials.

OUTPUT: JSON to stdout. Unknown keys return error with valid key list.`
