package initiative

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

// newInitiativeOpts collects the flag-bound values for `lin initiative new`.
type newInitiativeOpts struct {
	Name        string
	Description string
	Owner       string
	Status      string
	TargetDate  string
	Content     string
	Color       string
	Icon        string
}

func registerNew(parent *cobra.Command) {
	var opts newInitiativeOpts

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new initiative",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			opts.Name = args[0]

			client := linear.GetClient()
			input, err := buildInitiativeCreateInput(client, opts)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeCreate(context.Background(), client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			i := resp.InitiativeCreate.Initiative
			output.PrintJSON(map[string]any{
				"id":      i.Id,
				"slugId":  i.SlugId,
				"url":     i.Url,
				"name":    i.Name,
				"created": resp.InitiativeCreate.Success,
			})
		},
	}

	cmd.Flags().StringVar(&opts.Description, "description", "", "Initiative description")
	cmd.Flags().StringVar(&opts.Owner, "owner", "", "Initiative owner: name, email, or user ID")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Status: planned|active|completed")
	cmd.Flags().StringVar(&opts.TargetDate, "target-date", "", "Target date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.Content, "content", "", "Initiative content body (markdown)")
	cmd.Flags().StringVar(&opts.Color, "color", "", "Initiative color")
	cmd.Flags().StringVar(&opts.Icon, "icon", "", "Initiative icon")
	parent.AddCommand(cmd)
}

// buildInitiativeCreateInput resolves human-friendly flag values to Linear IDs
// and returns the GraphQL input. Any user-input error short-circuits with a
// descriptive message; resolver errors propagate up unchanged.
func buildInitiativeCreateInput(client graphql.Client, opts newInitiativeOpts) (linear.InitiativeCreateInput, error) {
	input := linear.InitiativeCreateInput{
		Name: opts.Name,
	}
	if opts.Description != "" {
		input.Description = &opts.Description
	}
	if opts.Content != "" {
		input.Content = &opts.Content
	}
	if opts.Color != "" {
		input.Color = &opts.Color
	}
	if opts.Icon != "" {
		input.Icon = &opts.Icon
	}
	if opts.TargetDate != "" {
		input.TargetDate = &opts.TargetDate
	}
	if opts.Status != "" {
		normalized, err := validateInitiativeStatus(opts.Status)
		if err != nil {
			return linear.InitiativeCreateInput{}, err
		}
		s := linear.InitiativeStatus(normalized)
		input.Status = &s
	}
	if opts.Owner != "" {
		u, err := resolvers.ResolveUser(client, opts.Owner)
		if err != nil {
			return linear.InitiativeCreateInput{}, err
		}
		input.OwnerId = &u.ID
	}
	return input, nil
}
