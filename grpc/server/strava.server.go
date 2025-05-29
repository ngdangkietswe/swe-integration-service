package server

import (
	"context"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"github.com/ngdangkietswe/swe-integration-service/grpc/service/strava"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type StravaGrpcServer struct {
	integration.UnimplementedStravaServiceServer
	stravaSvc strava.IStravaService
}

func NewStravaGrpcServer(stravaSvc strava.IStravaService) *StravaGrpcServer {
	return &StravaGrpcServer{
		stravaSvc: stravaSvc,
	}
}

// IntegrateStravaAccount is a function that implements the IntegrateStravaAccount method of the StravaServiceServer interface
func (s *StravaGrpcServer) IntegrateStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*common.EmptyResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.IntegrateStravaAccount)
}

// GetStravaAccount is a function that implements the GetStravaAccount method of the StravaServiceServer interface
func (s *StravaGrpcServer) GetStravaAccount(ctx context.Context, req *integration.GetStravaAccountReq) (*integration.GetStravaAccountResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.GetStravaAccount)
}

// SyncStravaActivities is a function that implements the SyncStravaActivities method of the StravaServiceServer interface
func (s *StravaGrpcServer) SyncStravaActivities(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.SyncStravaActivities)
}

// GetStravaActivities is a function that implements the GetStravaActivities method of the StravaServiceServer interface
func (s *StravaGrpcServer) GetStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq) (*integration.GetStravaActivitiesResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.GetStravaActivities)
}

// RemoveStravaAccount is a function that implements the RemoveStravaAccount method of the StravaServiceServer interface
func (s *StravaGrpcServer) RemoveStravaAccount(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.RemoveStravaAccount)
}

// RemoveStravaActivity is a function that implements the RemoveStravaActivity method of the StravaServiceServer interface
func (s *StravaGrpcServer) RemoveStravaActivity(ctx context.Context, req *common.IdReq) (*common.EmptyResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.RemoveStravaActivity)
}

// BulkRemoveStravaActivities is a function that implements the BulkRemoveStravaActivities method of the StravaServiceServer interface
func (s *StravaGrpcServer) BulkRemoveStravaActivities(ctx context.Context, req *common.IdsReq) (*common.EmptyResp, error) {
	return util.HandleGrpc(ctx, req, s.stravaSvc.BulkRemoveStravaActivities)
}
