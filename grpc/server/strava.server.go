package server

import (
	"context"
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
	return s.stravaSvc.IntegrateStravaAccount(ctx, req)
}

// GetStravaAccount is a function that implements the GetStravaAccount method of the StravaServiceServer interface
func (s *StravaGrpcServer) GetStravaAccount(ctx context.Context, req *integration.GetStravaAccountReq) (*integration.GetStravaAccountResp, error) {
	return s.stravaSvc.GetStravaAccount(ctx, req)
}
