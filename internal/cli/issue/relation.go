package issue

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

const validRelationTypes = "blocks | duplicate | related"

var validRelationTypeSet = map[string]linear.IssueRelationType{
	"blocks":    linear.IssueRelationTypeBlocks,
	"duplicate": linear.IssueRelationTypeDuplicate,
	"related":   linear.IssueRelationTypeRelated,
}

// resolveRelationType maps a user-provided relation type (any case) to the
// Linear enum value. Returns ok=false for unknown types.
func resolveRelationType(input string) (linear.IssueRelationType, bool) {
	v, ok := validRelationTypeSet[strings.ToLower(input)]
	return v, ok
}

// inverseRelationLabel rewrites the inverse-direction "blocks" relation to
// "blocked_by"; other types pass through.
func inverseRelationLabel(forward string) string {
	if forward == "blocks" {
		return "blocked_by"
	}
	return forward
}

func registerRelation(parent *cobra.Command) {
	relation := &cobra.Command{
		Use:   "relation",
		Short: "Issue relation operations",
	}
	parent.AddCommand(relation)

	registerRelationList(relation)
	registerRelationAdd(relation)
	registerRelationRemove(relation)

	output.HandleUnknownCommand(relation, "Run 'lin issue usage' for available relation subcommands")
}

func registerRelationList(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "list <issue-id>",
		Short: "List all relations on an issue (both directions)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueRelations(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			var mapped []any
			for _, r := range resp.Issue.Relations.Nodes {
				mapped = append(mapped, map[string]any{
					"id":           r.Id,
					"type":         r.Type,
					"relatedIssue": r.RelatedIssue.Identifier,
				})
			}
			for _, r := range resp.Issue.InverseRelations.Nodes {
				mapped = append(mapped, map[string]any{
					"id":           r.Id,
					"type":         inverseRelationLabel(r.Type),
					"relatedIssue": r.Issue.Identifier,
				})
			}

			output.PrintJSON(mapped)
		},
	})
}

func registerRelationAdd(parent *cobra.Command) {
	var (
		relType   string
		relatedID string
	)

	cmd := &cobra.Command{
		Use:   "add <issue-id>",
		Short: "Add a relation between two issues",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			enumVal, ok := resolveRelationType(relType)
			if !ok {
				output.PrintErrorf("Invalid relation type: %q. Valid values: %s", relType, validRelationTypes)
			}

			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueRelationCreate(ctx, client, linear.IssueRelationCreateInput{
				IssueId:        args[0],
				RelatedIssueId: relatedID,
				Type:           enumVal,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			rel := resp.IssueRelationCreate.IssueRelation
			output.PrintJSON(map[string]any{
				"id":                      rel.Id,
				"type":                    rel.Type,
				"issueIdentifier":         rel.Issue.Identifier,
				"relatedIssueIdentifier":  rel.RelatedIssue.Identifier,
				"created":                 resp.IssueRelationCreate.Success,
			})
		},
	}

	cmd.Flags().StringVar(&relType, "type", "", "Relation type: blocks|duplicate|related")
	_ = cmd.MarkFlagRequired("type")
	cmd.Flags().StringVar(&relatedID, "related", "", "Target issue ID or key")
	_ = cmd.MarkFlagRequired("related")
	parent.AddCommand(cmd)
}

func registerRelationRemove(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "remove <relation-id>",
		Short: "Remove a relation",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueRelationDelete(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.IssueRelationDelete.Success})
		},
	})
}
