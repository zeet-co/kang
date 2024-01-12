package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type GetSubGroupsForGroupResp struct {
	ID        uuid.UUID
	Name      string
	SubGroups []SubGroup
}

type SubGroup struct {
	ID   uuid.UUID
	Name string
}

func (c *Client) GetSubGroupsForGroup(ctx context.Context, group string) (*GetSubGroupsForGroupResp, error) {
	out := &GetSubGroupsForGroupResp{}

	_ = `# @genqlient
query getGroup($path: String) {
	project(path: $path) {
		id
		name
		environments {
			id
			name
		}
	}
}
`

	res, err := getGroup(ctx, c.gql, group)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out.ID = res.Project.Id
	out.Name = res.Project.Name
	out.SubGroups = make([]SubGroup, len(res.Project.Environments))
	for i, e := range res.Project.Environments {
		out.SubGroups[i] = SubGroup{
			ID:   e.Id,
			Name: e.Name,
		}
	}

	return out, nil
}
