package output

import (
	"sync"
	"time"

	libcli "github.com/shhac/lib-agent-cli/cli"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output/pretty"
)

// AllowPretty opts cmd into the pretty card format via the family's shared
// format allow-list, so libcli.NewRoot's validator accepts `--format pretty` on
// this command (and rejects it elsewhere with the standard error). lin still owns
// the rendering — this only registers the format choice with lib-agent-cli, the
// same way agent-deepweb opts into raw/text.
func AllowPretty(cmd *cobra.Command) {
	libcli.AllowFormats(cmd, string(FormatPretty))
}

// WantsPretty reports whether the resolved format is the pretty card renderer.
func WantsPretty() bool {
	return ResolveFormat(FormatNDJSON) == FormatPretty
}

var (
	widthMu   sync.RWMutex
	flagWidth int
)

// ConfigureWidth records the --width flag (0 = auto-detect).
func ConfigureWidth(w int) {
	widthMu.Lock()
	flagWidth = w
	widthMu.Unlock()
}

// PrettyOptions resolves the render options once for a pretty pass: effective
// width (flag or terminal), color capability, and the wall-clock reference for
// relative times.
func PrettyOptions() pretty.Options {
	widthMu.RLock()
	w := flagWidth
	widthMu.RUnlock()
	return pretty.Options{
		Width: pretty.TerminalWidth(w),
		Color: pretty.ColorCapable(),
		Now:   time.Now(),
	}
}
