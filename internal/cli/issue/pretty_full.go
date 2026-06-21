package issue

import (
	"context"
	"sync"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

// fetchFullSections fetches the --full data for an issue card — relations and
// comments — concurrently. A secondary fetch that fails yields an empty section
// rather than failing the whole card (the base detail already succeeded).
func fetchFullSections(client graphql.Client, id string) (relations []relationRow, comments []commentRow) {
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		resp, err := linear.IssueRelations(ctx, client, id)
		if err != nil {
			return
		}
		for _, r := range resp.Issue.Relations.Nodes {
			relations = append(relations, relationRow{Type: r.Type, Identifier: r.RelatedIssue.Identifier})
		}
		for _, r := range resp.Issue.InverseRelations.Nodes {
			relations = append(relations, relationRow{Type: inverseRelationLabel(r.Type), Identifier: r.Issue.Identifier})
		}
	}()

	go func() {
		defer wg.Done()
		resp, err := linear.IssueComments(ctx, client, id, 250, nil)
		if err != nil {
			return
		}
		for _, cm := range resp.Issue.Comments.Nodes {
			author := "Unknown"
			if cm.User != nil {
				author = cm.User.Name
			}
			comments = append(comments, commentRow{Author: author, CreatedAt: cm.CreatedAt, Body: cm.Body})
		}
	}()

	wg.Wait()
	return relations, comments
}
