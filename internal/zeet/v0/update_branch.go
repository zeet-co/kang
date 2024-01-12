package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) UpdateProjectBranch(ctx context.Context, projectID uuid.UUID, branch string) error {

	_ = `# @genqlient
mutation updateProjectBranch($id: ID!, $branch: String) {
	updateProject(input: {
		id: $id,
		productionBranch: $branch
	}) {
		id
	}
}
`

	_, err := updateProjectBranch(ctx, c.gql, projectID, branch)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
