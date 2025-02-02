package strava

import (
	"context"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type IStravaValidator interface {
	ValidateIntegrateStravaAccountRequest(ctx context.Context, request *integration.IntegrateStravaAccountReq) error
}
