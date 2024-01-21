package v0

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (c *Client) DuplicateProject(ctx context.Context, projectID, groupID, subGroupID uuid.UUID, newName string) (uuid.UUID, error) {

	fmt.Printf("duplicating project %s to group %s subgroup %s as %s \n", projectID, groupID, subGroupID, newName)

	_ = `# @genqlient
	mutation duplicateProject($id: UUID!, $groupID: UUID, $subGroupID: UUID, $name: String!) {
		duplicateProject(input: {
			enabled: true,
			environmentID: $subGroupID,
			id: $id,
			projectID: $groupID,
			name: $name
		}) {
			id
		}
	}
`

	res, err := duplicateProject(ctx, c.gql, projectID, groupID, subGroupID, newName)
	if err != nil {
		return uuid.Nil, err
	}
	return res.DuplicateProject.Id, nil
}
