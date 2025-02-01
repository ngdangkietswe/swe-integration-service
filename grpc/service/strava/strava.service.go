package strava

import (
	"context"
	"github.com/google/uuid"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	stravarepo "github.com/ngdangkietswe/swe-integration-service/data/repository/strava"
	"github.com/ngdangkietswe/swe-integration-service/grpc/mapper"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"go.uber.org/zap"
)

type stravaService struct {
	logger     *logger.Logger
	stravaRepo stravarepo.IStravaRepository
}

// IntegrateStravaAccount is a function that integrates a Strava account with the user's account.
func (s stravaService) IntegrateStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*common.EmptyResp, error) {
	_, err := s.stravaRepo.SaveStravaAccount(ctx, req)
	if err != nil {
		s.logger.Error("Failed to save Strava account", zap.String("error", err.Error()))
		return nil, err
	}

	return &common.EmptyResp{
		Success: true,
	}, nil
}

// GetStravaAccount is a function that gets the Strava account of the user.
func (s stravaService) GetStravaAccount(ctx context.Context, req *integration.GetStravaAccountReq) (*integration.GetStravaAccountResp, error) {
	var userId uuid.UUID
	if req.UserId != "" {
		userId = uuid.MustParse(req.UserId)
	} else {
		userId = uuid.MustParse(grpcutil.GetGrpcPrincipal(ctx).UserId)
	}

	stravaAccount, err := s.stravaRepo.GetStravaAccountByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("Failed to get Strava account", zap.String("error", err.Error()))
		return nil, err
	}

	return &integration.GetStravaAccountResp{
		Success: true,
		Resp: &integration.GetStravaAccountResp_StravaAccount{
			StravaAccount: mapper.AsMonoStravaAccount(stravaAccount),
		},
	}, nil
}

func NewStravaService(logger *logger.Logger, stravaRepo stravarepo.IStravaRepository) IStravaService {
	return &stravaService{
		logger:     logger,
		stravaRepo: stravaRepo,
	}
}
