package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Repo struct {
	ID   uuid.UUID
	Name string
}

func (c *Client) GetRepo(ctx context.Context, id uuid.UUID) (*Repo, error) {
	out := &Repo{}

	_ = `# @genqlient
query getRepo($id: UUID) {
  repo(id: $id) {
    id
		name
  }
}
`
	res, err := getRepo(ctx, c.gql, id)
	if err := copier.Copy(out, res.Repo); err != nil {
		return nil, err
	}

	return out, err
}
