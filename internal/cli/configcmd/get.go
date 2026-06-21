package configcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	libcli "github.com/shhac/lib-agent-cli/cli"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

func registerGet(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get [key...]",
		Short: "Show current settings",
		Args:  cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			// No args → list all settings (all valid keys in sorted order).
			keys := args
			if len(keys) == 0 {
				keys = validKeys
			}

			settings := config.GetSettings()
			format := string(output.ResolveFormat(output.FormatNDJSON))
			return libcli.EntityGet(output.Stdout(), format, keys, func(key string) (any, error) {
				def, ok := settingDefs[key]
				if !ok {
					return nil, apierrors.New(
						fmt.Sprintf("no setting %q", key),
						apierrors.FixableByAgent,
					).WithHint(fmt.Sprintf("valid keys: %s", validKeysStr()))
				}
				return map[string]any{"key": key, "value": def.get(settings)}, nil
			})
		},
	}
	cfg.AddCommand(cmd)
}
