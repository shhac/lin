package document

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdate(parent *cobra.Command) {
	update := &cobra.Command{
		Use:   "update",
		Short: "Update document fields",
	}
	parent.AddCommand(update)

	for _, m := range parentMoves() {
		registerParentMove(update, m)
	}

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
				output.HandleGraphQLError(err)
			}

			resp, err := linear.DocumentUpdate(ctx, client, doc.ID, buildInput(args[1]))
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.DocumentUpdate.Success})
		},
	})
}

// parentMove is a `document update <parent>` subcommand. A document has
// exactly one parent, so moving it replaces the old one, and there is no
// clearing: Linear acknowledges a null parent but leaves it in place.
type parentMove struct {
	kind        string // parent kind, matching a MapDocumentParents key
	placeholder string
	short       string
	resolve     func(client graphql.Client, input string) (string, error)
	input       func(id string) linear.DocumentUpdateInput
}

func parentMoves() []parentMove {
	return []parentMove{
		{
			kind: "project", placeholder: "project", short: "Move document to a project",
			resolve: func(c graphql.Client, in string) (string, error) {
				p, err := resolvers.ResolveProject(c, in)
				return p.ID, err
			},
			input: func(id string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{ProjectId: &id} },
		},
		{
			kind: "team", placeholder: "team", short: "Move document to a team",
			resolve: func(c graphql.Client, in string) (string, error) {
				t, err := resolvers.ResolveTeam(c, in)
				return t.ID, err
			},
			input: func(id string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{TeamId: &id} },
		},
		{
			kind: "initiative", placeholder: "initiative", short: "Move document to an initiative",
			resolve: func(c graphql.Client, in string) (string, error) {
				i, err := resolvers.ResolveInitiative(c, in)
				return i.ID, err
			},
			input: func(id string) linear.DocumentUpdateInput { return linear.DocumentUpdateInput{InitiativeId: &id} },
		},
	}
}

func registerParentMove(parent *cobra.Command, m parentMove) {
	parent.AddCommand(&cobra.Command{
		Use:   fmt.Sprintf("%s <id> <%s>", m.kind, m.placeholder),
		Short: m.short,
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			if strings.TrimSpace(args[1]) == "" {
				output.WriteError(noParentlessDocumentErr(m.kind))
			}
			result, err := moveDocument(context.Background(), linear.GetClient(), m, args[0], args[1])
			if err != nil {
				output.HandleGraphQLError(err)
			}
			output.PrintJSON(result)
		},
	})
}

func noParentlessDocumentErr(kind string) error {
	return apierrors.Newf(apierrors.FixableByAgent, "missing %s: a document always has one parent, so it can't be removed from its %s, only moved", kind, kind).
		WithHint("lin document update project|team|initiative <id> <target>")
}

// moveDocument re-parents the document and confirms from Linear's response
// that the requested parent is now its only one.
func moveDocument(ctx context.Context, client graphql.Client, m parentMove, docInput, targetInput string) (map[string]any, error) {
	doc, err := resolvers.ResolveDocument(client, docInput)
	if err != nil {
		return nil, err
	}
	targetID, err := m.resolve(client, targetInput)
	if err != nil {
		return nil, err
	}
	resp, err := linear.DocumentUpdate(ctx, client, doc.ID, m.input(targetID))
	if err != nil {
		return nil, err
	}
	parents := mappers.MapDocumentParents(resp.DocumentUpdate.Document)
	moved, ok := parents[m.kind]
	if !ok || moved["id"] != targetID || len(parents) != 1 {
		return nil, apierrors.Newf(apierrors.FixableByHuman, "Linear reported success but the document's parent is not the requested %s", m.kind).
			WithHint("check with: lin document get " + doc.ID)
	}
	return map[string]any{"updated": resp.DocumentUpdate.Success, "id": doc.ID, m.kind: moved}, nil
}
