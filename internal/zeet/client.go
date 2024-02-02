package zeet

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
	v1 "github.com/zeet-co/kang/internal/zeet/v1"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	v0Client *v0.Client
	v1Client *v1.Client
}

func New(token string) *Client {
	v0Client := v0.New(token)
	v1Client := v1.New(token)

	return &Client{
		v0Client,
		v1Client,
	}
}

func (c *Client) GetRepoByID(ctx context.Context, id uuid.UUID) (*v0.Repo, error) {
	return c.v0Client.GetRepoByID(ctx, id)
}

func (c *Client) GetRepoByName(ctx context.Context, name string) (uuid.UUID, error) {
	return c.v0Client.GetRepoByName(ctx, name)
}

func (c *Client) GetGroup(ctx context.Context, group string) (*v0.GetSubGroupsForGroupResp, error) {
	return c.v0Client.GetGroup(ctx, group)

}

func (c *Client) EnsureGroupsExist(group, subgroup string, teamID uuid.UUID) (uuid.UUID, uuid.UUID, error) {

	ctx := context.Background()
	groupID := uuid.Nil
	subGroupID := uuid.Nil

	resp, err := c.GetGroup(ctx, group)

	if err == nil {
		groupID = resp.ID

		for _, sg := range resp.SubGroups {
			if sg.Name == subgroup {
				subGroupID = sg.ID
				break
			}
		}

	} else if err == v0.NotFoundError {
		// create group
		groupID, err = c.v1Client.CreateGroup(context.Background(), group, teamID)
		if err != nil {
			return uuid.Nil, uuid.Nil, err
		}
	} else {
		return uuid.Nil, uuid.Nil, err
	}

	if subGroupID == uuid.Nil {
		subGroupID, err = c.v1Client.CreateSubGroup(ctx, subgroup, groupID, teamID)
		if err != nil {
			return uuid.Nil, uuid.Nil, err
		}
	}

	return groupID, subGroupID, nil
}

func (c *Client) DuplicateProject(ctx context.Context, projectID, groupID, subGroupID uuid.UUID, newName string) (uuid.UUID, error) {
	return c.v0Client.DuplicateProject(ctx, projectID, groupID, subGroupID, newName)
}

func (c *Client) UpdateProject(ctx context.Context, projectID uuid.UUID, input *v0.UpdateProjectInput) error {
	return c.v0Client.UpdateProject(ctx, projectID, input)
}

func (c *Client) UpdateEnvs(ctx context.Context, projectID uuid.UUID, input map[string]string) error {
	return c.v0Client.UpdateEnvs(ctx, projectID, input)
}

func (c *Client) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return c.v0Client.DeleteProject(ctx, id)
}

func (c *Client) GetProjectsByID(ctx context.Context, projectIDs []uuid.UUID) ([]*v0.Repo, error) {

	projects := make([]*v0.Repo, len(projectIDs))

	var wg sync.WaitGroup
	// var repos []*Repo
	eg := new(errgroup.Group)

	// Assuming you have a slice of inputs for GetRepo
	for i, id := range projectIDs {
		wg.Add(1)
		id := id // capture range variable
		i := i
		eg.Go(func() error {
			defer wg.Done()
			repo, err := c.GetRepoByID(ctx, id)
			if err != nil {
				return err
			}
			// Use a mutex or other synchronization method if needed
			projects[i] = repo
			return nil
		})
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Check if any goroutines returned an error
	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "could not fetch project information"))
	}

	return projects, nil
}
