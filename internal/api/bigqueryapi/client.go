package bigqueryapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/einride/bigquery-importer-google-workspace/internal/tables"
	"go.uber.org/zap"
	workspace "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

type JobClient struct {
	Config         JobConfig
	BigQueryClient *bigquery.Client
	Logger         *zap.Logger
}

func allTableRows() []tables.Row {
	return []tables.Row{
		&tables.GroupsRow{},
		&tables.UsersRow{},
		&tables.GroupMembersRow{},
	}
}

// EnsureTables ensures that all tables exists in an empty state. If a table already exists, and error will be returned.
func (c *JobClient) EnsureTables(ctx context.Context) error {
	c.Logger.Info("ensuring tables")
	for _, tableRow := range allTableRows() {
		if err := c.createTable(ctx, tableRow); err != nil {
			return err
		}
	}
	return nil
}

func (c *JobClient) PutGroups(ctx context.Context, groups ...*workspace.Group) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put groups: %w", err)
		}
	}()
	if len(groups) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(groups))
	for _, group := range groups {
		row := tables.GroupsRow{
			Org: c.Config.Org,
		}
		row.UnmarshalGroup(group)
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting groups", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.GroupsRow{}).Put(ctx, valueSavers)
}

func (c *JobClient) PutUsers(ctx context.Context, users ...*workspace.User) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put users: %w", err)
		}
	}()
	if len(users) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(users))
	for _, user := range users {
		row := tables.UsersRow{}
		err := row.UnmarshalUser(user)
		if err != nil {
			return err
		}
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting user", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.UsersRow{}).Put(ctx, valueSavers)
}

func (c *JobClient) PutGroupMembers(ctx context.Context, group *workspace.Group, members ...*workspace.Member) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put group members: %w", err)
		}
	}()
	if len(members) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(members))
	for _, member := range members {
		row := tables.GroupMembersRow{
			GroupId:   group.Id,
			GroupName: group.Name,
		}
		row.UnmarshalMember(member)
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting group members", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.GroupMembersRow{}).Put(ctx, valueSavers)
}

func (c *JobClient) inserter(row tables.Row) *bigquery.Inserter {
	tableID := row.TableID(c.Config.Date)
	if c.Config.AppendIDSuffix {
		tableID = tableID + "_" + c.Config.ID.String()
	}
	return c.BigQueryClient.Dataset(c.Config.Dataset).Table(tableID).Inserter()
}

func (c *JobClient) createTable(ctx context.Context, row tables.Row) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("create table %s: %w", row.TableID(c.Config.Date), err)
		}
	}()
	tableID := row.TableID(c.Config.Date)
	if c.Config.AppendIDSuffix {
		tableID = tableID + "_" + c.Config.ID.String()
	}
	table := c.BigQueryClient.Dataset(c.Config.Dataset).Table(tableID)
	_, err = table.Metadata(ctx)
	if err == nil {
		return fmt.Errorf("table already exists: %s", table.FullyQualifiedName())
	}
	var errAPI *googleapi.Error
	if ok := errors.As(err, &errAPI); err != nil && (!ok || errAPI.Code != http.StatusNotFound) {
		c.Logger.Debug("error",
			zap.Error(err))
		return err
	}
	c.Logger.Info("creating table", zap.Any("fullyQualifiedName", table.FullyQualifiedName()))
	return table.Create(ctx, row.TableMetadata())
}
