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

func (c *Client) EnsureGroupsExist(group, subgroup string, teamID uuid.UUID) (uuid.UUID, uuid.UUID, error) {

	ctx := context.Background()

	resp, err := c.v0Client.GetSubGroupsForGroup(ctx, group)
	if err != nil {
		//if not found, create
		return uuid.Nil, uuid.Nil, err
	}

	groupID := resp.ID
	subGroupID := uuid.Nil

	for _, sg := range resp.SubGroups {
		if sg.Name == subgroup {
			subGroupID = sg.ID
			break
		}
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
