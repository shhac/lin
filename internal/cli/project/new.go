package project

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/projectstatuses"
	"github.com/shhac/lin/internal/resolvers"
)

func registerNew(parent *cobra.Command) {
	var (
		team        string
		description string
		lead        string
		startDate   string
		targetDate  string
		status      string
		content     string
	)

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			name := args[0]

			if status != "" {
				normalized, err := projectstatuses.Validate(status)
				if err != nil {
					output.PrintError(err.Error())
				}
				status = normalized
			}

			var leadId *string
			if lead != "" {
				u, err := resolvers.ResolveUser(client, lead)
				if err != nil {
					output.PrintError(err.Error())
				}
				leadId = &u.ID
			}

			parts := strings.Split(team, ",")
			teamIds := make([]string, 0, len(parts))
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				t, err := resolvers.ResolveTeam(client, p)
				if err != nil {
					output.PrintError(err.Error())
				}
				teamIds = append(teamIds, t.ID)
			}

			input := linear.ProjectCreateInput{
				Name:    name,
				TeamIds: teamIds,
			}
			if description != "" {
				input.Description = &description
			}
			if leadId != nil {
				input.LeadId = leadId
			}
			if startDate != "" {
				input.StartDate = &startDate
			}
			if targetDate != "" {
				input.TargetDate = &targetDate
			}
			if status != "" {
				input.State = &status
			}
			if content != "" {
				input.Content = &content
			}

			resp, err := linear.ProjectCreate(ctx, client, input)
			if err != nil {
				output.PrintError(err.Error())
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

	cmd.Flags().StringVar(&team, "team", "", "Team ID(s) or key(s), comma-separated")
	_ = cmd.MarkFlagRequired("team")
	cmd.Flags().StringVar(&description, "description", "", "Project description")
	cmd.Flags().StringVar(&lead, "lead", "", "Project lead: name, email, or user ID")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&targetDate, "target-date", "", "Target date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&status, "status", "", "Status: backlog|planned|started|paused|completed|canceled")
	cmd.Flags().StringVar(&content, "content", "", "Project content body (markdown)")
	parent.AddCommand(cmd)
}
