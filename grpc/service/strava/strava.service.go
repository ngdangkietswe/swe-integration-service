package strava

import (
	"context"
	"github.com/google/uuid"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	stravarepo "github.com/ngdangkietswe/swe-integration-service/data/repository/strava"
	"github.com/ngdangkietswe/swe-integration-service/grpc/mapper"
	stravavalidator "github.com/ngdangkietswe/swe-integration-service/grpc/validator/strava"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"go.uber.org/zap"
)

type stravaService struct {
	logger          *logger.Logger
	stravaValidator stravavalidator.IStravaValidator
	stravaRepo      stravarepo.IStravaRepository
}

const GetAthleteActivitiesEndpoint = "https://www.strava.com/api/v3/athlete/activities"

// SyncStravaActivities is a function that syncs Strava activities with the user's account.
func (s stravaService) SyncStravaActivities(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error) {
	//TODO implement me
	panic("implement me")
}

// GetStravaActivities is a function that gets Strava activities of the user.
func (s stravaService) GetStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq) (*integration.GetStravaActivitiesResp, error) {
	//TODO implement me
	panic("implement me")
}

// IntegrateStravaAccount is a function that integrates a Strava account with the user's account.
func (s stravaService) IntegrateStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*common.EmptyResp, error) {
	// Validate request
	if err := s.stravaValidator.ValidateIntegrateStravaAccountRequest(ctx, req); err != nil {
		s.logger.Error("Failed to validate Strava account", zap.String("error", err.Error()))
		return nil, err
	}

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

func NewStravaService(
	logger *logger.Logger,
	stravavalidator stravavalidator.IStravaValidator,
	stravaRepo stravarepo.IStravaRepository) IStravaService {
	return &stravaService{
		logger:          logger,
		stravaValidator: stravavalidator,
		stravaRepo:      stravaRepo,
	}
}
