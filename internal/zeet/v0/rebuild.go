package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) RebuildProject(ctx context.Context, projectID uuid.UUID) error {
	_ = `# @genqlient
mutation buildRepo($id: ID!) {
	buildRepo(id: $id, noCache: false) {
		id
	}
}
`

	_, err := buildRepo(ctx, c.gql, projectID)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
