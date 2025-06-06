package repository

import (
	"github.com/ngdangkietswe/swe-integration-service/data/datasource"
	"github.com/ngdangkietswe/swe-integration-service/data/repository/cdcauthusers"
	"github.com/ngdangkietswe/swe-integration-service/data/repository/strava"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	datasource.NewEntClient,
	cdcauthusers.NewCdcAuthUsersRepository,
	strava.NewStravaRepository,
)
