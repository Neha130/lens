// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/devtron-labs/common-lib/pubsub-lib"
	"github.com/devtron-labs/lens/api"
	"github.com/devtron-labs/lens/client"
	"github.com/devtron-labs/lens/client/gitSensor"
	"github.com/devtron-labs/lens/internal/logger"
	"github.com/devtron-labs/lens/internal/sql"
	"github.com/devtron-labs/lens/pkg"
)

// Injectors from Wire.go:

func InitializeApp() (*App, error) {
	sugaredLogger := logger.NewSugardLogger()
	config, err := sql.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.NewDbConnection(config, sugaredLogger)
	if err != nil {
		return nil, err
	}
	leadTimeRepositoryImpl := sql.NewLeadTimeRepositoryImpl(db, sugaredLogger)
	pipelineMaterialRepositoryImpl := sql.NewPipelineMaterialRepositoryImpl(db, sugaredLogger)
	appReleaseRepositoryImpl := sql.NewAppReleaseRepositoryImpl(db, sugaredLogger, leadTimeRepositoryImpl, pipelineMaterialRepositoryImpl)
	deploymentMetricServiceImpl := pkg.NewDeploymentMetricServiceImpl(sugaredLogger, appReleaseRepositoryImpl, pipelineMaterialRepositoryImpl, leadTimeRepositoryImpl)
	gitSensorConfig, err := gitSensor.GetGitSensorConfig()
	if err != nil {
		return nil, err
	}
	gitSensorClientImpl, err := gitSensor.NewGitSensorSession(gitSensorConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	gitSensorGrpcClientConfig, err := gitSensor.GetConfig()
	if err != nil {
		return nil, err
	}
	gitSensorGrpcClientImpl := gitSensor.NewGitSensorGrpcClientImpl(sugaredLogger, gitSensorGrpcClientConfig)
	ingestionServiceImpl := pkg.NewIngestionServiceImpl(sugaredLogger, appReleaseRepositoryImpl, pipelineMaterialRepositoryImpl, leadTimeRepositoryImpl, gitSensorClientImpl, gitSensorGrpcClientImpl)
	restHandlerImpl := api.NewRestHandlerImpl(sugaredLogger, deploymentMetricServiceImpl, ingestionServiceImpl)
	muxRouter := api.NewMuxRouter(sugaredLogger, restHandlerImpl)
	pubSubClientServiceImpl := pubsub_lib.NewPubSubClientServiceImpl(sugaredLogger)
	natsSubscriptionImpl, err := client.NewNatsSubscription(pubSubClientServiceImpl, sugaredLogger, ingestionServiceImpl)
	if err != nil {
		return nil, err
	}
	app := NewApp(muxRouter, sugaredLogger, db, ingestionServiceImpl, natsSubscriptionImpl, pubSubClientServiceImpl)
	return app, nil
}
