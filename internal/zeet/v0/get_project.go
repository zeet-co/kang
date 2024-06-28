package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Deployment struct {
	ID        uuid.UUID        `json:"id"`
	Status    DeploymentStatus `json:"status"`
	Endpoints []string         `json:"endpoints"`
}

type Repo struct {
	ID                   uuid.UUID         `json:"id"`
	Name                 string            `json:"name"`
	Owner                string            `json:"owner"`
	GroupName            string            `json:"groupName"`
	SubGroupName         string            `json:"subGroupName"`
	ProductionDeployment Deployment        `json:"deployment"`
	DatabaseEnvs         map[string]string `json:"databaseEnvs"`
	Envs                 map[string]string `json:"envs"`
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
			status
		}
    databaseEnvs {
			name
			value
		}
		envs {
			name
			value
		}
  }
}
`
	res, err := getRepo(ctx, c.gql, &id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dbEnvs := make(map[string]string, len(res.Repo.DatabaseEnvs))
	for _, e := range res.Repo.DatabaseEnvs {
		dbEnvs[e.Name] = e.Value
	}

	envs := make(map[string]string, len(res.Repo.Envs))
	for _, e := range res.Repo.Envs {
		envs[e.Name] = e.Value
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
			Status:    res.Repo.ProductionDeployment.Status,
		},
		DatabaseEnvs: dbEnvs,
		Envs:         envs,
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
