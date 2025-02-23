package strava

import (
	"context"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type IStravaService interface {
	IntegrateStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*common.EmptyResp, error)
	GetStravaAccount(ctx context.Context, req *integration.GetStravaAccountReq) (*integration.GetStravaAccountResp, error)
	SyncStravaActivities(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error)
	GetStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq) (*integration.GetStravaActivitiesResp, error)
	RemoveStravaAccount(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error)
	RemoveStravaActivity(ctx context.Context, req *common.IdReq) (*common.EmptyResp, error)
	BulkRemoveStravaActivities(ctx context.Context, req *common.IdsReq) (*common.EmptyResp, error)
}
