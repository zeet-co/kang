package v0

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) UpdateEnvs(ctx context.Context, projectID uuid.UUID, vars map[string]string) error {
	// input *SetRepoEnvsInput) error {

	envs := make([]EnvVarInput, 0, len(vars))

	for k, v := range vars {
		envs = append(envs, EnvVarInput{
			Name:  k,
			Value: v,
		})
	}

	input := SetRepoEnvsInput{
		Id:   projectID,
		Envs: envs,
	}

	_ = `# @genqlient
mutation updateEnvs($input: SetRepoEnvsInput!) {
	setRepoEnvs(input: $input) {
		id
	}
}
`

	_, err := updateEnvs(ctx, c.gql, input)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
