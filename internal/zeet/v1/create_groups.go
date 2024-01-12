package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) CreateGroup(ctx context.Context, group string, teamID uuid.UUID) (uuid.UUID, error) {

	_ = `# @genqlient

mutation createGroup(
	$input: CreateGroupInput!
) {
	createGroup(input: $input) {
		id
	}
}
`

	groupRes, err := createGroup(ctx, c.gql, CreateGroupInput{
		Name:   group,
		TeamId: teamID,
	})

	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	return groupRes.CreateGroup.Id, nil
}

func (c *Client) CreateSubGroup(ctx context.Context, subgroup string, groupID, teamID uuid.UUID) (uuid.UUID, error) {

	_ = `# @genqlient

mutation createSubGroup(
	$input: CreateSubGroupInput!
) {
	createSubGroup(input: $input) {
		id
	}
}
`

	subGroupRes, err := createSubGroup(ctx, c.gql, CreateSubGroupInput{
		GroupId: groupID,
		Name:    subgroup,
	})

	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	return subGroupRes.CreateSubGroup.Id, nil
}
