package resolvers

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
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
	if err := requireInput("customer", input); err != nil {
		return ResolvedCustomer{}, err
	}
	resp, err := linear.CustomerGet(ctx(), client, input)
	if err == nil {
		return ResolvedCustomer{
			ID: resp.Customer.Id, Name: resp.Customer.Name, SlugId: resp.Customer.SlugId,
		}, nil
	}
	filter := filters.BuildCustomerNameFilter(input)
	listResp, err := linear.CustomerList(ctx(), client, filter, nameMatchLimit, nil)
	if err != nil {
		return ResolvedCustomer{}, err
	}
	matches := matchByNameOrID(input, listResp.Customers.Nodes,
		func(c linear.CustomerListCustomersCustomerConnectionNodesCustomer) string { return c.SlugId },
		func(c linear.CustomerListCustomersCustomerConnectionNodesCustomer) string { return c.Name })
	switch len(matches) {
	case 1:
		match := matches[0]
		return ResolvedCustomer{ID: match.Id, Name: match.Name, SlugId: match.SlugId}, nil
	case 0:
		return ResolvedCustomer{}, apierrors.Newf(apierrors.FixableByAgent, "customer not found: %q, provide a UUID, slug ID, or exact name", input)
	default:
		return ResolvedCustomer{}, ambiguousErr("customer", "customers", input, matches, describeCustomer, "use the slug ID or UUID")
	}
}

func describeCustomer(c linear.CustomerListCustomersCustomerConnectionNodesCustomer) string {
	desc := fmt.Sprintf("%s (slug: %s", c.Name, c.SlugId)
	if c.Status.Name != "" {
		desc += ", status: " + c.Status.Name
	}
	if c.Owner != nil {
		desc += ", owner: " + c.Owner.Name
	}
	return desc + ")"
}
