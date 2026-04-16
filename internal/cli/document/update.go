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

	registerUpdateProject(update)

	registerSimpleDocumentUpdate(update, "title <id> <new-title>", "Update document title",
		func(v string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{Title: &v} })
	registerSimpleDocumentUpdate(update, "content <id> <content>", "Update document content",
		func(v string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{Content: &v} })
	registerSimpleDocumentUpdate(update, "icon <id> <icon>", "Update document icon",
		func(v string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{Icon: &v} })
	registerSimpleDocumentUpdate(update, "color <id> <color>", "Update document color",
		func(v string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{Color: &v} })

	output.HandleUnknownCommand(update, "Run `lin document usage` for available update subcommands")
}

func registerSimpleDocumentUpdate(parent *cobra.Command, use, short string, buildInput func(string) linear.DocumentUpdateInput) {
	parent.AddCommand(&cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, buildInput(args[1]))
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	})
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
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
