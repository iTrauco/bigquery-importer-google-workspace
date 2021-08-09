package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	workspace "google.golang.org/api/admin/directory/v1"
)

// GroupMembersRow follows the structure of the WebAPI. See the official documentation for field descriptions:
// https://developers.google.com/admin-sdk/directory/reference/rest/v1/members
type GroupMembersRow struct {
	Id               string `bigquery:"id"`
	GroupId          string `bigquery:"group_id"`
	GroupName        string `bigquery:"group_name"`
	DeliverySettings string `bigquery:"delivery_settings"`
	Email            string `bigquery:"email"`
	Etag             string `bigquery:"etag"`
	Kind             string `bigquery:"kind"`
	Role             string `bigquery:"role"`
	Status           string `bigquery:"status"`
	Type             string `bigquery:"type"`
}

func (g *GroupMembersRow) TableID(date civil.Date) string {
	return "group_members_" + strings.ReplaceAll(date.String(), "-", "")
}

func (g *GroupMembersRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   g.Schema(),
		InsertID: g.InsertID(jobID),
		Struct:   g,
	}
}

func (g *GroupMembersRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(g)
	return schema
}

func (g *GroupMembersRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "Group Members follows the structure of the WebAPI. See the official documentation for field " +
			"descriptions: https://developers.google.com/admin-sdk/directory/reference/rest/v1/members",
		Schema: g.Schema(),
	}
}

func (g *GroupMembersRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		g.Id,
		g.GroupId,
	}, "-")
}

func (g *GroupMembersRow) UnmarshalMember(wg *workspace.Member) {
	g.Id = wg.Id
	g.DeliverySettings = wg.DeliverySettings
	g.Email = wg.Email
	g.Etag = wg.Etag
	g.Kind = wg.Kind
	g.Role = wg.Role
	g.Status = wg.Status
	g.Type = wg.Type
}
