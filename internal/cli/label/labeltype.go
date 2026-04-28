package label

import (
	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
)

const (
	labelTypeIssue   = "issue"
	labelTypeProject = "project"
)

// addTypeFlag attaches `--type issue|project` (default issue) to the command.
func addTypeFlag(cmd *cobra.Command, dest *string) {
	cmd.Flags().StringVar(dest, "type", labelTypeIssue, "Label type: issue (default) or project")
}

// validateType returns an error if the value isn't a recognised label type.
func validateType(t string) error {
	if t != labelTypeIssue && t != labelTypeProject {
		return apierrors.Newf(apierrors.FixableByAgent,
			"--type must be %q or %q (got %q)", labelTypeIssue, labelTypeProject, t)
	}
	return nil
}

// rejectTeamForProject returns an error if --team was supplied alongside
// --type=project. Project labels are workspace-scoped and have no team.
func rejectTeamForProject(typeFlag, teamFlag string) error {
	if typeFlag == labelTypeProject && teamFlag != "" {
		return apierrors.Newf(apierrors.FixableByAgent,
			"--team is not valid with --type=project (project labels are workspace-scoped)").
			WithHint("re-run without --team, or drop --type=project to search issue labels")
	}
	return nil
}

// validateLabelTypeFlags runs the type/team guards that every `lin label`
// subcommand shares. Returns the first error or nil.
func validateLabelTypeFlags(typeFlag, teamFlag string) error {
	if err := validateType(typeFlag); err != nil {
		return err
	}
	return rejectTeamForProject(typeFlag, teamFlag)
}

// labelNotFoundErr produces the "no matches" error used by `label get` for
// both issue and project labels, keeping wording in lockstep.
func labelNotFoundErr(kind, input, searchHint string) error {
	return apierrors.Newf(apierrors.FixableByAgent, "%s not found: %q", kind, input).
		WithHint(searchHint)
}

// ambiguousLabelErr produces the "too many matches" error, kept in lockstep
// across issue and project label paths.
func ambiguousLabelErr(kind, input string, n int, disambigHint, listHint string) error {
	return apierrors.Newf(apierrors.FixableByAgent,
		"%d %ss match %q — %s", n, kind, input, disambigHint).
		WithHint(listHint)
}
