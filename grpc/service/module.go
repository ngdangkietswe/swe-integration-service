package service

import (
	"github.com/ngdangkietswe/swe-integration-service/grpc/service/strava"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	strava.NewStravaService,
)
