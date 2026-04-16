package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/projectstatuses"
	"github.com/shhac/lin/internal/ptr"
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
	registerUpdateLead(update)
	registerUpdatePriority(update)

	registerSimpleProjectUpdate(update, "description <id> <description>", "Update project description",
		func(v string) linear.ProjectUpdateInput { return linear.ProjectUpdateInput{Description: &v} })
	registerSimpleProjectUpdate(update, "content <id> <content>", "Update project content (markdown body)",
		func(v string) linear.ProjectUpdateInput { return linear.ProjectUpdateInput{Content: &v} })
	registerSimpleProjectUpdate(update, "start-date <id> <date>", "Update project start date",
		func(v string) linear.ProjectUpdateInput { return linear.ProjectUpdateInput{StartDate: &v} })
	registerSimpleProjectUpdate(update, "target-date <id> <date>", "Update project target date",
		func(v string) linear.ProjectUpdateInput { return linear.ProjectUpdateInput{TargetDate: &v} })
	registerSimpleProjectUpdate(update, "icon <id> <icon>", "Update project icon",
		func(v string) linear.ProjectUpdateInput { return linear.ProjectUpdateInput{Icon: &v} })
	registerSimpleProjectUpdate(update, "color <id> <color>", "Update project color",
		func(v string) linear.ProjectUpdateInput { return linear.ProjectUpdateInput{Color: &v} })

	output.HandleUnknownCommand(update, "Run `lin project usage` for available update subcommands")
}

func registerSimpleProjectUpdate(parent *cobra.Command, use, short string, buildInput func(string) linear.ProjectUpdateInput) {
	parent.AddCommand(&cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, buildInput(args[1]))
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	})
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
				output.HandleGraphQLError(err)
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
			normalized, err := projectstatuses.Validate(args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				State: &normalized,
			})
			if err != nil {
				output.HandleGraphQLError(err)
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
				output.HandleGraphQLError(err)
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
				Priority: ptr.To(p),
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
