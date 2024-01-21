package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type GetSubGroupsForGroupResp struct {
	ID        uuid.UUID
	Name      string
	SubGroups []SubGroup
}

type SubGroup struct {
	ID       uuid.UUID
	Name     string
	Projects []Project
}

type Project struct {
	ID      uuid.UUID
	Name    string
	Enabled bool
}

func (c *Client) GetGroup(ctx context.Context, group string) (*GetSubGroupsForGroupResp, error) {
	out := &GetSubGroupsForGroupResp{}

	_ = `# @genqlient
query getGroup($path: String) {
	project(path: $path) {
		id
		name
		environments {
			id
			name
			repos {
				id
				enabled
				name
			}
		}
	}
}
`

	res, err := getGroup(ctx, c.gql, group)

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

	out.ID = res.Project.Id
	out.Name = res.Project.Name
	out.SubGroups = make([]SubGroup, len(res.Project.Environments))
	for i, e := range res.Project.Environments {

		projects := make([]Project, len(e.Repos))
		for i, r := range e.Repos {
			projects[i] = Project{
				ID:      r.Id,
				Name:    r.Name,
				Enabled: r.Enabled,
			}
		}
		out.SubGroups[i] = SubGroup{
			ID:       e.Id,
			Name:     e.Name,
			Projects: projects,
		}
	}

	return out, nil
}
