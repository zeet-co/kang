package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Cluster struct {
	ID     uuid.UUID
	Region string
}

func (c *Client) GetClusterByID(ctx context.Context, clusterID, teamID uuid.UUID) (*Cluster, error) {
	var out *Cluster

	_ = `# @genqlient
query getCluster($userID: ID!, $clusterID: UUID!) {
	user(id: $userID) {
		cluster(id: $clusterID) {
			id
			region
		}
	}
}
`
	res, err := getCluster(ctx, c.gql, teamID, clusterID)

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
	out = &Cluster{
		ID:     res.User.Cluster.Id,
		Region: *res.User.Cluster.Region,
	}
	return out, nil
}
