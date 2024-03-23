package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) DeleteSubGroup(ctx context.Context, subGroupID uuid.UUID) error {

	fmt.Printf("Deleting sub-group %s\n", subGroupID)

	_ = `# @genqlient

mutation deleteSubGroup(
	$id: UUID!
) {
	deleteSubGroup(id: $id)
}
`

	_, err := deleteSubGroup(ctx, c.gql, subGroupID)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
