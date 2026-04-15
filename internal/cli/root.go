package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "lin",
		Short:         "Linear CLI for humans and LLMs",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	return root
}

func Execute(version string) error {
	err := newRootCmd(version).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return err
}
