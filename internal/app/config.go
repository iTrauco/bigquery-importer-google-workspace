package app

import (
	"github.com/einride/bigquery-importer-google-workspace/internal/api/bigqueryapi"
	"github.com/einride/bigquery-importer-google-workspace/internal/api/workspaceapi"
)

type Config struct {
	Logger struct {
		ServiceName string `required:"true"`
		Level       string `required:"true"`
		Development bool   `required:"true"`
	}

	BigQueryClient struct {
		ProjectID string `required:"true"`
	}

	WorkspaceClient struct {
		APIKeySecret string `required:"true"`
		JWTConfig    struct {
			Subject string `required:"true"`
		}
	}

	Job bigqueryapi.JobConfig

	Work workspaceapi.Config
}
