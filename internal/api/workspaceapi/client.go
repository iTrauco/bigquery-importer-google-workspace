package workspaceapi

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	workspace "google.golang.org/api/admin/directory/v1"
)

type WorkspaceClient struct {
	DirectoryService *workspace.Service
	Logger           *zap.Logger
	Config           *Config
}

// ListDomainGroups fetches all admin.Group in a domain and executes the provided function on them.
func (c *WorkspaceClient) ListDomainGroups(
	ctx context.Context,
	put func(context.Context, ...*workspace.Group) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list domain %s groups: %w", c.Config.Domain, err)
		}
	}()
	var pageToken string
	for {
		groups, err := c.DirectoryService.Groups.List().Domain(c.Config.Domain).PageToken(pageToken).Do()
		if err != nil {
			return fmt.Errorf("groups: %w", err)
		}
		err = put(ctx, groups.Groups...)
		if err != nil {
			return err
		}
		pageToken = groups.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return nil
}

// ListDomainUsers fetches all admin.User in a domain and executes the provided function on them.
func (c *WorkspaceClient) ListDomainUsers(
	ctx context.Context,
	put func(context.Context, ...*workspace.User) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list domain %s users: %w", c.Config.Domain, err)
		}
	}()
	var pageToken string
	const projection = "full"
	for {
		users, err := c.DirectoryService.Users.List().Projection(projection).Domain(c.Config.Domain).PageToken(pageToken).Do()
		if err != nil {
			return fmt.Errorf("users: %w", err)
		}
		err = put(ctx, users.Users...)
		if err != nil {
			return err
		}
		pageToken = users.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return nil
}

// ListGroupMembers fetches all admin.Member for the provided group and execute the provided function on them.
func (c *WorkspaceClient) ListGroupMembers(
	ctx context.Context,
	group *workspace.Group,
	put func(context.Context, *workspace.Group, ...*workspace.Member) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list member in group %s: %v", group.Id, err)
		}
	}()
	var pageToken string
	for {
		members, err := c.DirectoryService.Members.List(group.Id).PageToken(pageToken).Do()
		if err != nil {
			return fmt.Errorf("group members: %w", err)
		}
		err = put(ctx, group, members.Members...)
		if err != nil {
			return fmt.Errorf("group members: %w", err)
		}
		pageToken = members.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return nil
}
