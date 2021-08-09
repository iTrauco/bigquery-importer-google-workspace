package app

import (
	"context"
	"fmt"

	"github.com/einride/bigquery-importer-google-workspace/internal/api/bigqueryapi"
	"github.com/einride/bigquery-importer-google-workspace/internal/api/workspaceapi"
	"go.uber.org/zap"
	workspace "google.golang.org/api/admin/directory/v1"
)

type App struct {
	BigQueryJobClient *bigqueryapi.JobClient
	WorkspaceClient   *workspaceapi.WorkspaceClient
	Logger            *zap.Logger
}

// Run creates new tables and scrapes Google Workspace on Groups, Group Members and Users.
func (a *App) Run(ctx context.Context) error {
	a.Logger.Info("running")
	defer a.Logger.Info("stopped")
	if err := a.BigQueryJobClient.EnsureTables(ctx); err != nil {
		return err
	}
	if err := a.exportUsers(ctx); err != nil {
		return err
	}
	if err := a.exportGroups(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) exportGroups(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export groups: %w", err)
		}
	}()
	a.Logger.Info("exporting groups")
	return a.WorkspaceClient.ListDomainGroups(ctx, func(ctx context.Context, groups ...*workspace.Group) error {
		err := a.BigQueryJobClient.PutGroups(ctx, groups...)
		if err != nil {
			return nil
		}
		for _, group := range groups {
			err := a.exportGroupMembers(ctx, group)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *App) exportUsers(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export users: %w", err)
		}
	}()
	a.Logger.Info("exporting users")
	return a.WorkspaceClient.ListDomainUsers(ctx, a.BigQueryJobClient.PutUsers)
}

func (a *App) exportGroupMembers(ctx context.Context, group *workspace.Group) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export group members: %w", err)
		}
	}()
	a.Logger.Info("exporting group members")
	err = a.WorkspaceClient.ListGroupMembers(ctx, group, a.BigQueryJobClient.PutGroupMembers)
	if err != nil {
		return err
	}
	return nil
}
