package zeet

import (
	"context"

	"github.com/google/uuid"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
	v1 "github.com/zeet-co/kang/internal/zeet/v1"
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

func (c *Client) GetRepo(ctx context.Context, id uuid.UUID) (*v0.GetRepoResponse, error) {
	return c.v0Client.GetRepo(ctx, id)
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

func (c *Client) UpdateProjectBranch(ctx context.Context, projectID uuid.UUID, branch string) error {
	return c.v0Client.UpdateProjectBranch(ctx, projectID, branch)
}

func (c *Client) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return c.v0Client.DeleteProject(ctx, id)
}
