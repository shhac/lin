package output

import (
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output/pretty"
)

// prettyAnnotation marks a command (or command group) as accepting
// --format pretty. lin sets and reads this itself rather than relying on
// lib-agent-cli's format validator, because lin replaces NewRoot's
// PersistentPreRunE with its own config-aware one (so libcli's validator isn't
// in the path).
const prettyAnnotation = "lin.pretty-format"

// AllowPretty opts cmd into the pretty card format. Call it in the command's
// registration alongside building the get command.
func AllowPretty(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[prettyAnnotation] = "1"
}

// prettyAllowed reports whether cmd or any ancestor opted into pretty.
func prettyAllowed(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Annotations[prettyAnnotation] == "1" {
			return true
		}
	}
	return false
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
