package app

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2/google"
	workspace "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func InitSecretManagerClient(
	ctx context.Context,
	logger *zap.Logger,
) (*secretmanager.Client, func(), error) {
	logger.Info("init Secret Manager client")
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("init Secret Manager client: %w", err)
	}
	cleanup := func() {
		logger.Info("closing Secret Manager client")
		if err := client.Close(); err != nil {
			logger.Error("close Secret Manager client", zap.Error(err))
		}
	}
	return client, cleanup, nil
}

func InitWorkspaceClient(
	ctx context.Context,
	config *Config,
	secretManagerClient *secretmanager.Client,
	logger *zap.Logger,
) (_ *workspace.Service, err error) {
	logger.Info("init Google Workspace Directory service", zap.Any("config", config.WorkspaceClient))
	secretVersion, err := secretManagerClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: config.WorkspaceClient.APIKeySecret,
	})
	if err != nil {
		return nil, fmt.Errorf("init Workspace Directory service: %w", err)
	}
	jwtConfig, err := google.JWTConfigFromJSON(secretVersion.Payload.Data, workspace.AdminDirectoryUserScope, workspace.AdminDirectoryGroupScope)
	if err != nil {
		return nil, fmt.Errorf("init Workspace Directory service: %w", err)
	}
	jwtConfig.Subject = config.WorkspaceClient.JWTConfig.Subject
	service, err := workspace.NewService(ctx, option.WithTokenSource(jwtConfig.TokenSource(ctx)))
	if err != nil {
		return nil, fmt.Errorf("init G Suite Directory service: %w", err)
	}
	return service, nil
}

func InitBigQueryClient(
	ctx context.Context,
	config *Config,
	logger *zap.Logger,
) (_ *bigquery.Client, _ func(), err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("init BigQuery client: %w", err)
		}
	}()
	logger.Info("init BigQuery client", zap.Any("config", config.BigQueryClient))
	client, err := bigquery.NewClient(ctx, config.BigQueryClient.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		logger.Info("closing BigQuery client")
		if err := client.Close(); err != nil {
			logger.Error("close BigQuery client", zap.Error(err))
		}
	}
	return client, cleanup, nil
}

func InitLogger(
	config *Config,
) (_ *zap.Logger, _ func(), err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("init logger: %w", err)
		}
	}()
	var zapConfig zap.Config
	var zapOptions []zap.Option
	if config.Logger.Development {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig = zapdriver.NewProductionEncoderConfig()
		zapOptions = append(
			zapOptions,
			zapdriver.WrapCore(
				zapdriver.ServiceName(config.Logger.ServiceName),
				zapdriver.ReportAllErrors(true),
			),
		)
	}
	if err := zapConfig.Level.UnmarshalText([]byte(config.Logger.Level)); err != nil {
		return nil, nil, err
	}
	logger, err := zapConfig.Build(zapOptions...)
	if err != nil {
		return nil, nil, err
	}
	logger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	logger.Info("logger initialized")
	cleanup := func() {
		logger.Info("closing logger, goodbye")
		_ = logger.Sync()
	}
	return logger, cleanup, nil
}
