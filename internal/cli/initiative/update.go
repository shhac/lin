package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdate(parent *cobra.Command) {
	update := &cobra.Command{
		Use:   "update",
		Short: "Update initiative fields",
	}
	parent.AddCommand(update)

	registerUpdateName(update)
	registerUpdateStatus(update)
	registerUpdateOwner(update)

	registerSimpleInitiativeUpdate(update, "description <id> <description>", "Update initiative description",
		func(v string) linear.InitiativeUpdateInput { return linear.InitiativeUpdateInput{Description: &v} })
	registerSimpleInitiativeUpdate(update, "content <id> <content>", "Update initiative content (markdown body)",
		func(v string) linear.InitiativeUpdateInput { return linear.InitiativeUpdateInput{Content: &v} })
	registerSimpleInitiativeUpdate(update, "target-date <id> <date>", "Update initiative target date",
		func(v string) linear.InitiativeUpdateInput { return linear.InitiativeUpdateInput{TargetDate: &v} })
	registerSimpleInitiativeUpdate(update, "color <id> <color>", "Update initiative color",
		func(v string) linear.InitiativeUpdateInput { return linear.InitiativeUpdateInput{Color: &v} })
	registerSimpleInitiativeUpdate(update, "icon <id> <icon>", "Update initiative icon",
		func(v string) linear.InitiativeUpdateInput { return linear.InitiativeUpdateInput{Icon: &v} })

	output.HandleUnknownCommand(update, "Run `lin initiative usage` for available update subcommands")
}

func registerSimpleInitiativeUpdate(parent *cobra.Command, use, short string, buildInput func(string) linear.InitiativeUpdateInput) {
	parent.AddCommand(&cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, buildInput(args[1]))
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.InitiativeUpdate.Success})
		},
	})
}

func registerUpdateName(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "name <id> <new-name>",
		Short: "Update initiative name",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, linear.InitiativeUpdateInput{
				Name: &args[1],
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{
				"id":      resolved.ID,
				"name":    args[1],
				"updated": resp.InitiativeUpdate.Success,
			})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateStatus(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status <id> <new-status>",
		Short: "Update initiative status",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			normalized, err := validateInitiativeStatus(args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			s := linear.InitiativeStatus(normalized)
			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, linear.InitiativeUpdateInput{
				Status: &s,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.InitiativeUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateOwner(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "owner <id> <user>",
		Short: "Update initiative owner",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			user, err := resolvers.ResolveUser(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, linear.InitiativeUpdateInput{
				OwnerId: &user.ID,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.InitiativeUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
