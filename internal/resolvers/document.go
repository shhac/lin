package resolvers

import (
	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

type ResolvedDocument struct {
	ID     string
	SlugId string
	Title  string
}

func ResolveDocument(client graphql.Client, input string) (ResolvedDocument, error) {
	resp, err := linear.DocumentGet(ctx(), client, input)
	if err == nil {
		return ResolvedDocument{
			ID: resp.Document.Id, SlugId: resp.Document.SlugId, Title: resp.Document.Title,
		}, nil
	}
	filter := &linear.DocumentFilter{
		SlugId: &linear.StringComparator{Eq: ptr.To(input)},
	}
	listResp, err := linear.DocumentList(ctx(), client, filter, 1, nil, nil)
	if err != nil {
		return ResolvedDocument{}, err
	}
	if len(listResp.Documents.Nodes) == 0 {
		return ResolvedDocument{}, apierrors.Newf(apierrors.FixableByAgent, "document not found: %q, provide a UUID or slug ID", input)
	}
	d := listResp.Documents.Nodes[0]
	return ResolvedDocument{ID: d.Id, SlugId: d.SlugId, Title: d.Title}, nil
}
