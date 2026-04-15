package document

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
		Short: "Update document fields",
	}
	parent.AddCommand(update)

	registerUpdateTitle(update)
	registerUpdateContent(update)
	registerUpdateProject(update)
	registerUpdateIcon(update)
	registerUpdateColor(update)
	output.HandleUnknownCommand(update, "Run `lin document usage` for available update subcommands")
}

func registerUpdateTitle(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "title <id> <new-title>",
		Short: "Update document title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, linear.DocumentUpdateInput{
				Title: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateContent(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "content <id> <markdown>",
		Short: "Update document content",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, linear.DocumentUpdateInput{
				Content: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateProject(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "project <id> <project>",
		Short: "Move document to project",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resolved, err := resolvers.ResolveProject(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, linear.DocumentUpdateInput{
				ProjectId: &resolved.ID,
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateIcon(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "icon <id> <icon>",
		Short: "Update document icon",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, linear.DocumentUpdateInput{
				Icon: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUpdateColor(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "color <id> <color>",
		Short: "Update document color",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, linear.DocumentUpdateInput{
				Color: &args[1],
			})
			if err != nil {
				output.PrintError(err.Error())
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
