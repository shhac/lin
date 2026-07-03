package project

import (
	"context"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/projectstatuses"
	"github.com/shhac/lin/internal/resolvers"
)

// newProjectOpts collects the flag-bound values for `lin project new`.
type newProjectOpts struct {
	Name        string
	Team        string
	Description string
	Lead        string
	StartDate   string
	TargetDate  string
	Status      string
	Content     string
}

func registerNew(parent *cobra.Command) {
	var opts newProjectOpts

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			opts.Name = args[0]

			client := linear.GetClient()
			input, err := buildProjectCreateInput(client, opts)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectCreate(context.Background(), client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			result := map[string]any{
				"created": resp.ProjectCreate.Success,
			}
			if resp.ProjectCreate.Project != nil {
				p := resp.ProjectCreate.Project
				result["id"] = p.Id
				result["slugId"] = p.SlugId
				result["url"] = p.Url
				result["name"] = p.Name
			}

			output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&opts.Team, "team", "", "Team ID(s) or key(s), comma-separated")
	_ = cmd.MarkFlagRequired("team")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Project description")
	cmd.Flags().StringVar(&opts.Lead, "lead", "", "Project lead: name, email, or user ID")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.TargetDate, "target-date", "", "Target date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Status: backlog|planned|started|paused|completed|canceled")
	cmd.Flags().StringVar(&opts.Content, "content", "", "Project content body (markdown)")
	parent.AddCommand(cmd)
}

// buildProjectCreateInput resolves human-friendly flag values to Linear IDs and
// returns the GraphQL input. Any user-input error short-circuits with a
// descriptive message; resolver errors propagate up unchanged.
func buildProjectCreateInput(client graphql.Client, opts newProjectOpts) (linear.ProjectCreateInput, error) {
	teamIds, err := resolveTeamIDs(client, opts.Team)
	if err != nil {
		return linear.ProjectCreateInput{}, err
	}

	input := linear.ProjectCreateInput{
		Name:    opts.Name,
		TeamIds: teamIds,
	}
	if opts.Description != "" {
		input.Description = &opts.Description
	}
	if opts.Lead != "" {
		u, err := resolvers.ResolveUser(client, opts.Lead)
		if err != nil {
			return linear.ProjectCreateInput{}, err
		}
		input.LeadId = &u.ID
	}
	if opts.StartDate != "" {
		input.StartDate = &opts.StartDate
	}
	if opts.TargetDate != "" {
		input.TargetDate = &opts.TargetDate
	}
	if opts.Status != "" {
		normalized, err := projectstatuses.Validate(opts.Status)
		if err != nil {
			return linear.ProjectCreateInput{}, err
		}
		input.State = &normalized
	}
	if opts.Content != "" {
		input.Content = &opts.Content
	}
	return input, nil
}

// resolveTeamIDs resolves a comma-separated list of team keys/IDs/names to their
// Linear team IDs, skipping blank entries. Any resolver error short-circuits.
func resolveTeamIDs(client graphql.Client, csv string) ([]string, error) {
	parts := strings.Split(csv, ",")
	teamIds := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		t, err := resolvers.ResolveTeam(client, p)
		if err != nil {
			return nil, err
		}
		teamIds = append(teamIds, t.ID)
	}
	return teamIds, nil
}
