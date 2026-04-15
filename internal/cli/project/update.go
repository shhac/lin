package project

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/projectstatuses"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdate(parent *cobra.Command) {
	update := &cobra.Command{
		Use:   "update",
		Short: "Update project fields",
	}
	parent.AddCommand(update)

	registerUpdateTitle(update)
	registerUpdateStatus(update)
	registerUpdateDescription(update)
	registerUpdateContent(update)
	registerUpdateLead(update)
	registerUpdateStartDate(update)
	registerUpdateTargetDate(update)
	registerUpdatePriority(update)
	registerUpdateIcon(update)
	registerUpdateColor(update)
	output.HandleUnknownCommand(update, "Run `lin project usage` for available update subcommands")
}

func registerUpdateTitle(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "title <id> <new-title>",
		Short: "Update project title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Name: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{
				"id":      resolved.ID,
				"name":    args[1],
				"updated": resp.ProjectUpdate.Success,
			})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateStatus(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status <id> <new-status>",
		Short: "Update project status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			newStatus := args[1]
			valid := false
			for _, s := range projectstatuses.List {
				if strings.EqualFold(s, newStatus) {
					valid = true
					break
				}
			}
			if !valid {
				output.PrintErrorf("Invalid project status: %q. Valid values: %s", newStatus, projectstatuses.Values)
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			lower := strings.ToLower(newStatus)
			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				State: &lower,
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateDescription(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "description <id> <description>",
		Short: "Update project description",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Description: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateContent(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "content <id> <content>",
		Short: "Update project content (markdown body)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Content: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateLead(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "lead <id> <user>",
		Short: "Update project lead",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			user, err := resolvers.ResolveUser(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				LeadId: &user.ID,
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateStartDate(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "start-date <id> <date>",
		Short: "Update project start date",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				StartDate: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateTargetDate(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "target-date <id> <date>",
		Short: "Update project target date",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				TargetDate: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdatePriority(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "priority <id> <priority>",
		Short: "Update project priority",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := priorities.Resolve(args[1])
			if !ok {
				output.PrintErrorf("Invalid priority: %q. Valid values: %s", args[1], priorities.Values)
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Priority: intPtr(p),
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateIcon(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "icon <id> <icon>",
		Short: "Update project icon",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Icon: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateColor(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "color <id> <color>",
		Short: "Update project color",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Color: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
