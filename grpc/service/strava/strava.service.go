package strava

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-integration-service/data/repository"
	stravarepo "github.com/ngdangkietswe/swe-integration-service/data/repository/strava"
	"github.com/ngdangkietswe/swe-integration-service/grpc/mapper"
	stravavalidator "github.com/ngdangkietswe/swe-integration-service/grpc/validator/strava"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"time"
)

type stravaService struct {
	client          *ent.Client
	logger          *logger.Logger
	stravaValidator stravavalidator.IStravaValidator
	stravaRepo      stravarepo.IStravaRepository
}

const (
	GetAthleteActivitiesEndpoint      = "https://www.strava.com/api/v3/athlete/activities"
	RefreshExpiredAccessTokenEndpoint = "https://www.strava.com/api/v3/oauth/token"
	MaxPerPage                        = 100
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
	accessToken := stravaAccount.AccessToken

	// Check if access token has expired. If so, refresh the access token.
	if stravaAccount.ExpiresAt.Before(time.Now()) {
		s.logger.Info("Access token of Strava account has expired, refreshing access token...")

		accessToken = s.refreshExpiredAccessTokenOfStravaAccount(ctx, stravaAccount)
		if accessToken == "" {
			return nil, errors.New("failed to refresh access token")
		}
	}

	s.logger.Info("Getting Strava activities from Strava API endpoint", zap.String("endpoint", GetAthleteActivitiesEndpoint))

	for {
		var pageStravaActivities []map[string]interface{}
		_, err = client.R().
			SetAuthToken(accessToken).
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
		if err = repository.WithTx(ctx, s.client, s.logger, func(tx *ent.Tx) error {
			return s.stravaRepo.SyncStravaActivities(ctx, tx, userId, stravaAccount.AthleteID, newStravaActivities)
		}); err != nil {
			s.logger.Error("Failed to sync Strava activities", zap.String("error", err.Error()))
			return nil, err
		}
	}

	return &common.EmptyResp{
		Success: true,
	}, nil
}

// refreshExpiredAccessTokenOfStravaAccount is a function that refreshes the expired access token of a Strava account.
func (s stravaService) refreshExpiredAccessTokenOfStravaAccount(ctx context.Context, stravaAccount *ent.StravaAccount) string /*new_access_token*/ {
	var tokenResp map[string]interface{}
	client := resty.New()

	_, err := client.R().
		SetFormData(map[string]string{
			"client_id":     config.GetString("STRAVA_CLIENT_ID", ""),
			"client_secret": config.GetString("STRAVA_CLIENT_SECRET", ""),
			"grant_type":    "refresh_token",
			"refresh_token": stravaAccount.RefreshToken,
		}).
		SetResult(&tokenResp).
		Post(RefreshExpiredAccessTokenEndpoint)

	if err != nil {
		s.logger.Error("Failed to refresh expired access token of Strava account", zap.String("error", err.Error()))
		return ""
	}

	newAccessToken := tokenResp["access_token"].(string)

	s.logger.Info("Successfully refreshed expired access token of Strava account", zap.String("new_access_token", newAccessToken))

	if err = repository.WithTx(ctx, s.client, s.logger, func(tx *ent.Tx) error {
		return s.stravaRepo.UpdateTokenStravaAccount(ctx, tx, stravaAccount.ID, newAccessToken, tokenResp["refresh_token"].(string), int64(tokenResp["expires_at"].(float64)))
	}); err != nil {
		s.logger.Error("Failed to update Strava account with new access token", zap.String("error", err.Error()))
		return ""
	}

	return newAccessToken
}

// GetStravaActivities is a function that gets Strava activities of the user.
func (s stravaService) GetStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq) (*integration.GetStravaActivitiesResp, error) {
	normalizePageable := util.NormalizePageable(req.Pageable)

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
				PageMetaData: util.AsPageMetaData(normalizePageable, count),
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

	if _, err := repository.WithTxResult(ctx, s.client, s.logger, func(tx *ent.Tx) (*ent.StravaAccount, error) {
		return s.stravaRepo.SaveStravaAccount(ctx, tx, req)
	}); err != nil {
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

// RemoveStravaAccount is a function that removes the Strava account of the user.
func (s stravaService) RemoveStravaAccount(ctx context.Context, req *common.EmptyReq) (*common.EmptyResp, error) {
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

	if err = repository.WithTx(ctx, s.client, s.logger, func(tx *ent.Tx) error {
		return s.stravaRepo.RemoveStravaAccountByUserId(ctx, tx, stravaAccount.ID)
	}); err != nil {
		s.logger.Error("Failed to remove Strava account", zap.String("error", err.Error()))
		return nil, err
	}

	return &common.EmptyResp{}, nil
}

// RemoveStravaActivity is a function that removes a Strava activity.
func (s stravaService) RemoveStravaActivity(ctx context.Context, req *common.IdReq) (*common.EmptyResp, error) {
	exists, err := s.stravaRepo.ExistsStravaActivityById(ctx, uuid.MustParse(req.Id))
	if err != nil {
		s.logger.Error("Failed to check if Strava activity exists", zap.String("error", err.Error()))
		return nil, err
	} else if !exists {
		return nil, errors.New("strava activity not found")
	}

	if err = repository.WithTx(ctx, s.client, s.logger, func(tx *ent.Tx) error {
		return s.stravaRepo.DeleteStravaActivityById(ctx, tx, uuid.MustParse(req.Id))
	}); err != nil {
		s.logger.Error("Failed to delete Strava activity", zap.String("error", err.Error()))
		return nil, err
	}

	return &common.EmptyResp{
		Success: true,
	}, nil
}

// BulkRemoveStravaActivities is a function that bulk removes Strava activities.
func (s stravaService) BulkRemoveStravaActivities(ctx context.Context, req *common.IdsReq) (*common.EmptyResp, error) {
	existsAll, err := s.stravaRepo.ExistsAllStravaActivitiesByIdIn(ctx, util.Convert2UUID(req.Ids))
	if err != nil {
		s.logger.Error("Failed to check if all Strava activities exist", zap.String("error", err.Error()))
		return nil, err
	} else if !existsAll {
		return nil, errors.New("strava activities not found")
	}

	if err = repository.WithTx(ctx, s.client, s.logger, func(tx *ent.Tx) error {
		return s.stravaRepo.DeleteStravaActivitiesByIdIn(ctx, tx, util.Convert2UUID(req.Ids))
	}); err != nil {
		s.logger.Error("Failed to delete Strava activities", zap.String("error", err.Error()))
		return nil, err
	}

	return &common.EmptyResp{
		Success: true,
	}, nil
}

func NewStravaService(
	client *ent.Client,
	logger *logger.Logger,
	stravavalidator stravavalidator.IStravaValidator,
	stravaRepo stravarepo.IStravaRepository) IStravaService {
	return &stravaService{
		client:          client,
		logger:          logger,
		stravaValidator: stravavalidator,
		stravaRepo:      stravaRepo,
	}
}
