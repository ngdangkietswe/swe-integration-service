package validator

import (
	"github.com/ngdangkietswe/swe-integration-service/grpc/validator/strava"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	strava.NewStravaValidator,
)
