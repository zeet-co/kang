package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (c *Client) GetTeamName(ctx context.Context, teamID uuid.UUID) (*string, error) {
	_ = `# @genqlient
query getTeam($id: ID!) {
	user(id: $id) {
		id
		login
	}
}
`

	res, err := getTeam(ctx, c.gql, teamID)

	var errList gqlerror.List
	if errors.As(err, &errList) {
		for _, err := range errList {
			if err.Message == "not found" {
				return nil, NotFoundError
			}
		}
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &res.User.Login, nil
}
