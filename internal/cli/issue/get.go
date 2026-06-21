package issue

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/output/pretty"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Get issue details: title, description, status, assignee, labels, relationships",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fetchDetail := func(client graphql.Client, id string) (map[string]any, error) {
				resp, err := linear.IssueGet(context.Background(), client, id)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}
				return mappers.MapIssueDetail(resp.Issue), nil
			}
			if output.WantsPretty() {
				full, _ := cmd.Flags().GetBool("full")
				getOne := func(client graphql.Client, id string) (any, error) {
					d, err := fetchDetail(client, id)
					if err != nil {
						return nil, err
					}
					card := issueCard{detail: d}
					if full {
						card.relations, card.comments = fetchFullSections(client, id)
					}
					return card, nil
				}
				return shared.GetEntitiesPretty(args, getOne, func(item any, opts pretty.Options) string {
					return renderIssueCard(item.(issueCard), opts)
				})
			}
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				return fetchDetail(client, id)
			})
		},
	}
	output.AllowPretty(cmd)

	parent.AddCommand(cmd)
}
