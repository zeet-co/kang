package v0

import (
	"context"

	"github.com/google/uuid"
)

func (c *Client) DeleteProject(ctx context.Context, id uuid.UUID) error {

	_ = `# @genqlient
mutation deleteRepo($id: ID!) {
  deleteRepo(id: $id)
}
`

	_, err := deleteRepo(ctx, c.gql, id)

	return err

}
