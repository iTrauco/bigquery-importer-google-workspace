package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	workspace "google.golang.org/api/admin/directory/v1"
)

// GroupsRow follows the structure of the WebAPI. See the official documentation for field descriptions:
// https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups
type GroupsRow struct {
	Org                string   `bigquery:"org"`
	Id                 string   `bigquery:"id"`
	Email              string   `bigquery:"email"`
	Name               string   `bigquery:"name"`
	Description        string   `bigquery:"description"`
	AdminCreated       bool     `bigquery:"admin_created"`
	DirectMembersCount int64    `bigquery:"direct_members_count"`
	Kind               string   `bigquery:"kind"`
	Etag               string   `bigquery:"etag"`
	Aliases            []string `bigquery:"aliases"`
	NonEditableAliases []string `bigquery:"non_editable_aliases"`
}

var _ Row = &GroupsRow{}

func (g *GroupsRow) TableID(date civil.Date) string {
	return "groups_" + strings.ReplaceAll(date.String(), "-", "")
}

func (g *GroupsRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   g.Schema(),
		InsertID: g.InsertID(jobID),
		Struct:   g,
	}
}

func (g *GroupsRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(g)
	return schema
}

func (g *GroupsRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "Groups follows the structure of the WebAPI. See the official documentation for field " +
			"descriptions: https://developers.google.com/admin-sdk/directory/reference/rest/v1/groups",
		Schema: g.Schema(),
	}
}

func (g *GroupsRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		g.Id,
	}, "-")
}

func (g *GroupsRow) UnmarshalGroup(wg *workspace.Group) {
	g.Id = wg.Id
	g.Email = wg.Email
	g.Name = wg.Name
	g.Description = wg.Description
	g.AdminCreated = wg.AdminCreated
	g.DirectMembersCount = wg.DirectMembersCount
	g.Kind = wg.Kind
	g.Etag = wg.Etag
	g.Aliases = wg.Aliases
	g.NonEditableAliases = wg.NonEditableAliases
}
