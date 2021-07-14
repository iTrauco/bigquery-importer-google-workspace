//+build wireinject

package app

import (
	"context"

	"github.com/einride/bigquery-importer-google-workspace/internal/api/bigqueryapi"
	"github.com/einride/bigquery-importer-google-workspace/internal/api/workspaceapi"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitApp(ctx context.Context, logger *zap.Logger, config *Config) (*App, func(), error) {
	panic(
		wire.Build(
			wire.Struct(new(App), "*"),
			InitBigQueryClient,
			InitWorkspaceClient,
			InitSecretManagerClient,
			wire.Struct(new(workspaceapi.WorkspaceClient), "*"), wire.FieldsOf(&config, "Work"),
			wire.Struct(new(bigqueryapi.JobClient), "*"), wire.FieldsOf(&config, "Job"),
		),
	)
}
