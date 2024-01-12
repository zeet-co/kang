package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type GetRepoResponse struct {
	ID uuid.UUID
}

func (c *Client) GetRepo(ctx context.Context, id uuid.UUID) (*GetRepoResponse, error) {
	out := &GetRepoResponse{}

	_ = `# @genqlient
query getRepo($id: UUID) {
  repo(id: $id) {
    id
  }
}
`
	res, err := getRepo(ctx, c.gql, id)
	if err := copier.Copy(out, res.Repo); err != nil {
		return nil, err
	}

	return out, err
}
