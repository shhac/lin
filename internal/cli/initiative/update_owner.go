package initiative

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/clearable"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/resolvers"
)

var ownerArg = clearable.Arg{Parent: "lin initiative update", Name: "owner", Placeholder: "user", Clearable: true}

func registerUpdateOwner(parent *cobra.Command) {
	parent.AddCommand(ownerArg.Command("Update or remove initiative owner", setOwner, clearOwner))
}

func setOwner(ctx context.Context, client graphql.Client, initiativeInput, userInput string) (map[string]any, error) {
	initiative, err := resolvers.ResolveInitiative(client, initiativeInput)
	if err != nil {
		return nil, err
	}
	user, err := resolvers.ResolveUser(client, userInput)
	if err != nil {
		return nil, err
	}
	resp, err := linear.InitiativeUpdate(ctx, client, initiative.ID, linear.InitiativeUpdateInput{OwnerId: &user.ID})
	if err != nil {
		return nil, err
	}
	result := map[string]any{"updated": resp.InitiativeUpdate.Success, "id": initiative.ID}
	if o := resp.InitiativeUpdate.Initiative.Owner; o != nil {
		result["owner"] = map[string]any{"id": o.Id, "name": o.Name}
	}
	return result, nil
}

func clearOwner(ctx context.Context, client graphql.Client, initiativeInput string) (map[string]any, error) {
	initiative, err := resolvers.ResolveInitiative(client, initiativeInput)
	if err != nil {
		return nil, err
	}
	resp, err := linear.InitiativeClearOwner(ctx, client, initiative.ID)
	if err != nil {
		return nil, err
	}
	if resp.InitiativeUpdate.Initiative.Owner != nil {
		return nil, ownerArg.StillSetErr()
	}
	return map[string]any{"updated": resp.InitiativeUpdate.Success, "id": initiative.ID, "cleared": "owner"}, nil
}
