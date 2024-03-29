package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Deployment struct {
	ID        uuid.UUID `json:"id"`
	Endpoints []string  `json:"endpoints"`
}

type Repo struct {
	ID                   uuid.UUID  `json:"id"`
	Name                 string     `json:"name"`
	Owner                string     `json:"owner"`
	GroupName            string     `json:"groupName"`
	SubGroupName         string     `json:"subGroupName"`
	ProductionDeployment Deployment `json:"deployment"`
}

func (c *Client) GetRepoByID(ctx context.Context, id uuid.UUID) (*Repo, error) {
	out := &Repo{}

	_ = `# @genqlient
query getRepo($id: UUID) {
  repo(id: $id) {
    id
		name
		owner {
			login
		}
		project{
			name
		}
		projectEnvironment {
			name
		}
		productionDeployment {
			id
			endpoints
		}
  }
}
`
	res, err := getRepo(ctx, c.gql, &id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out = &Repo{
		ID:           res.Repo.Id,
		Name:         res.Repo.Name,
		Owner:        res.Repo.Owner.Login,
		GroupName:    res.Repo.Project.Name,
		SubGroupName: res.Repo.ProjectEnvironment.Name,
		ProductionDeployment: Deployment{
			ID:        res.Repo.ProductionDeployment.Id,
			Endpoints: res.Repo.ProductionDeployment.Endpoints,
		},
	}

	return out, err
}

func (c *Client) GetRepoByName(ctx context.Context, name string) (uuid.UUID, error) {

	_ = `# @genqlient
query getRepoByName($name: String) {
  repo(path: $name) {
    id
		name
  }
}
`
	res, err := getRepoByName(ctx, c.gql, &name)

	if err != nil {
		return uuid.Nil, err
	}

	return res.Repo.Id, nil
}
