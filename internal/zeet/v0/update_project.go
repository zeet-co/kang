package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) UpdateProject(ctx context.Context, projectID uuid.UUID, input *UpdateProjectInput) error {

	_ = `# @genqlient
mutation updateProject($input: UpdateProjectInput!) {
	updateProject(input: $input) {
		id
	}
}
`

	_, err := updateProject(ctx, c.gql, *input)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
