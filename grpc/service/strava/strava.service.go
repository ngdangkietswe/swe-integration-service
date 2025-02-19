package strava

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	stravarepo "github.com/ngdangkietswe/swe-integration-service/data/repository/strava"
	"github.com/ngdangkietswe/swe-integration-service/grpc/mapper"
	"github.com/ngdangkietswe/swe-integration-service/grpc/utils"
	stravavalidator "github.com/ngdangkietswe/swe-integration-service/grpc/validator/strava"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type stravaService struct {
	logger          *logger.Logger
	stravaValidator stravavalidator.IStravaValidator
	stravaRepo      stravarepo.IStravaRepository
}

const (
	GetAthleteActivitiesEndpoint = "https://www.strava.com/api/v3/athlete/activities"
	MaxPerPage                   = 100
)

// SyncStravaActivities is a function that syncs Strava activities with the user's account.
func (s stravaService) SyncStravaActivities(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error) {
	userId := uuid.MustParse(grpcutil.GetGrpcPrincipal(ctx).UserId)

	stravaAccount, err := s.stravaRepo.GetStravaAccountByUserId(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			s.logger.Error("Strava account not found", zap.String("error", err.Error()))
			return nil, err
		} else {
			s.logger.Error("Failed to get Strava account", zap.String("error", err.Error()))
			return nil, err
		}
	}

	var stravaActivities, newStravaActivities []map[string]interface{}
	page := 1
	client := resty.New()

	s.logger.Info("Getting Strava activities from Strava API endpoint", zap.String("endpoint", GetAthleteActivitiesEndpoint))

	for {
		var pageStravaActivities []map[string]interface{}
		_, err = client.R().
			SetAuthToken(stravaAccount.AccessToken).
			SetResult(&pageStravaActivities).
			Get(fmt.Sprintf("%s?per_page=%d&page=%d", GetAthleteActivitiesEndpoint, MaxPerPage, page))
		if err != nil {
			s.logger.Error("Failed to get Strava activities from Strava API endpoint", zap.String("error", err.Error()))
			return nil, err
		}

		if len(pageStravaActivities) == 0 {
			break
		}

		stravaActivities = append(stravaActivities, pageStravaActivities...)
		page++
	}

	entStravaActivities, err := s.stravaRepo.GetListStravaActivitiesByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("Failed to get Strava activities", zap.String("error", err.Error()))
		return nil, err
	}

	existsIntegrationIds := lo.Map(entStravaActivities, func(sa *ent.StravaActivity, _ int) int64 {
		return sa.StravaActivityID
	})

	lo.ForEach(stravaActivities, func(item map[string]interface{}, _ int) {
		if !lo.Contains(existsIntegrationIds, int64(item["id"].(float64))) {
			newStravaActivities = append(newStravaActivities, item)
		}
	})

	if len(newStravaActivities) > 0 {
		s.logger.Info("Syncing Strava activities...", zap.Int("count", len(newStravaActivities)))
		err = s.stravaRepo.SyncStravaActivities(ctx, userId, stravaAccount.AthleteID, newStravaActivities)
		if err != nil {
			s.logger.Error("Failed to sync Strava activities", zap.String("error", err.Error()))
			return nil, err
		}
	}

	return &common.EmptyResp{
		Success: true,
	}, nil
}

// GetStravaActivities is a function that gets Strava activities of the user.
func (s stravaService) GetStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq) (*integration.GetStravaActivitiesResp, error) {
	normalizePageable := utils.NormalizePageable(req.Pageable)

	userId := uuid.MustParse(grpcutil.GetGrpcPrincipal(ctx).UserId)
	stravaActivities, count, err := s.stravaRepo.GetListStravaActivities(ctx, req, userId, normalizePageable)
	if err != nil {
		s.logger.Error("Failed to get Strava activities", zap.String("error", err.Error()))
		return nil, err
	}

	return &integration.GetStravaActivitiesResp{
		Success: true,
		Resp: &integration.GetStravaActivitiesResp_Data_{
			Data: &integration.GetStravaActivitiesResp_Data{
				PageMetaData: utils.AsPageMetaData(normalizePageable, count),
				Activities:   mapper.AsListStravaActivity(stravaActivities),
			},
		},
	}, nil
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
