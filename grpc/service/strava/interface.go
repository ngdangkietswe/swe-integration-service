package strava

import (
	"context"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type IStravaService interface {
	IntegrateStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*common.EmptyResp, error)
	GetStravaAccount(ctx context.Context, req *integration.GetStravaAccountReq) (*integration.GetStravaAccountResp, error)
}
