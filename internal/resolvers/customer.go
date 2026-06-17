package resolvers

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
)

type ResolvedCustomer struct {
	ID     string
	Name   string
	SlugId string
}

// ResolveCustomer resolves a customer from a UUID, slug, or exact name.
func ResolveCustomer(client graphql.Client, input string) (ResolvedCustomer, error) {
	resp, err := linear.CustomerGet(ctx(), client, input)
	if err == nil {
		return ResolvedCustomer{
			ID: resp.Customer.Id, Name: resp.Customer.Name, SlugId: resp.Customer.SlugId,
		}, nil
	}
	filter := filters.BuildCustomerNameFilter(input)
	listResp, err := linear.CustomerList(ctx(), client, filter, 1, nil)
	if err != nil {
		return ResolvedCustomer{}, err
	}
	if len(listResp.Customers.Nodes) == 0 {
		return ResolvedCustomer{}, fmt.Errorf("customer not found: %q, provide a UUID, slug ID, or exact name", input)
	}
	c := listResp.Customers.Nodes[0]
	return ResolvedCustomer{ID: c.Id, Name: c.Name, SlugId: c.SlugId}, nil
}
