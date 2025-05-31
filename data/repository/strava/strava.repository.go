package strava

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	grpcutil "github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/stravaaccount"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/stravaactivity"
	"github.com/ngdangkietswe/swe-integration-service/grpc/utils"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"github.com/samber/lo"
	"time"
)

type stravaRepository struct {
	entClient *ent.Client
}

// UpdateTokenStravaAccount is a function that updates the token of a Strava account.
func (s stravaRepository) UpdateTokenStravaAccount(ctx context.Context, tx *ent.Tx, id uuid.UUID, accessToken string, refreshToken string, expiresAt int64) error {
	return tx.StravaAccount.UpdateOneID(id).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		SetExpiresAt(time.Unix(expiresAt, 0).UTC()).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
}

// DeleteStravaActivityById is a function that deletes a Strava activity by ID.
func (s stravaRepository) DeleteStravaActivityById(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	userId := uuid.MustParse(grpcutil.GetGrpcPrincipal(ctx).UserId)
	return tx.StravaActivity.DeleteOneID(id).Where(stravaactivity.UserID(userId)).Exec(ctx)
}

// DeleteStravaActivitiesByIdIn is a function that deletes Strava activities by IDs.
func (s stravaRepository) DeleteStravaActivitiesByIdIn(ctx context.Context, tx *ent.Tx, ids []uuid.UUID) error {
	userId := uuid.MustParse(grpcutil.GetGrpcPrincipal(ctx).UserId)
	_, err := tx.StravaActivity.Delete().
		Where(
			stravaactivity.IDIn(ids...),
			stravaactivity.UserID(userId)).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ExistsAllStravaActivitiesByIdIn is a function that checks if all Strava activities exist by IDs.
func (s stravaRepository) ExistsAllStravaActivitiesByIdIn(ctx context.Context, ids []uuid.UUID) (bool, error) {
	count, err := s.entClient.StravaActivity.Query().Where(stravaactivity.IDIn(ids...)).Count(ctx)
	if err != nil {
		return false, err
	}
	return count == len(ids), nil
}

// ExistsStravaActivityById is a function that checks if a Strava activity exists by ID.
func (s stravaRepository) ExistsStravaActivityById(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.entClient.StravaActivity.Query().Where(stravaactivity.ID(id)).Exist(ctx)
}

// GetStravaActivityById is a function that gets a Strava activity by ID.
func (s stravaRepository) GetStravaActivityById(ctx context.Context, id uuid.UUID) (*ent.StravaActivity, error) {
	return s.entClient.StravaActivity.Get(ctx, id)
}

// GetListStravaActivitiesByIdIn is a function that gets a list of Strava activities by IDs.
func (s stravaRepository) GetListStravaActivitiesByIdIn(ctx context.Context, ids []uuid.UUID) ([]*ent.StravaActivity, error) {
	return s.entClient.StravaActivity.Query().Where(stravaactivity.IDIn(ids...)).All(ctx)
}

// RemoveStravaAccountByUserId is a function that removes a Strava account by user ID.
func (s stravaRepository) RemoveStravaAccountByUserId(ctx context.Context, tx *ent.Tx, stravaAccountId uuid.UUID) error {
	return tx.StravaAccount.DeleteOneID(stravaAccountId).Exec(ctx)
}

// SyncStravaActivities is a function that syncs Strava activities with the user's account.
func (s stravaRepository) SyncStravaActivities(ctx context.Context, tx *ent.Tx, userId uuid.UUID, athleteId int64, stravaActivities []map[string]interface{}) error {
	var stravaActivityCreates []*ent.StravaActivityCreate

	lo.ForEach(stravaActivities, func(item map[string]interface{}, _ int) {
		startDate, err := time.Parse(time.RFC3339, item["start_date"].(string))
		if err != nil {
			return
		}

		stravaActivityId := int64(item["id"].(float64))
		stravaActivityUrl := fmt.Sprintf("https://www.strava.com/activities/%d", stravaActivityId)

		stravaActivityCreates = append(stravaActivityCreates, s.entClient.StravaActivity.Create().
			SetUserID(userId).
			SetAthleteID(athleteId).
			SetStravaActivityID(stravaActivityId).
			SetActivityName(item["name"].(string)).
			SetActivityType(0). // TODO: Set activity type in the future
			SetActivityURL(stravaActivityUrl).
			SetDistance(item["distance"].(float64)).
			SetStartDate(startDate).
			SetMovingTime(int32(item["moving_time"].(float64))).
			SetElapsedTime(int32(item["elapsed_time"].(float64))).
			SetTotalElevationGain(item["total_elevation_gain"].(float64)).
			SetAverageSpeed(item["average_speed"].(float64)).
			SetMaxSpeed(item["max_speed"].(float64)))
	})

	return tx.StravaActivity.CreateBulk(stravaActivityCreates...).Exec(ctx)
}

// GetListStravaActivitiesByUserId is a function that gets a list of Strava activities by user ID.
func (s stravaRepository) GetListStravaActivitiesByUserId(ctx context.Context, userId uuid.UUID) ([]*ent.StravaActivity, error) {
	return s.entClient.StravaActivity.Query().Where(stravaactivity.UserID(userId)).All(ctx)
}

// GetListStravaActivities is a function that gets a list of Strava activities by user ID.
func (s stravaRepository) GetListStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq, userId uuid.UUID, pageable *common.Pageable) ([]*ent.StravaActivity, int64, error) {
	query := s.entClient.StravaActivity.Query().Where(stravaactivity.UserID(userId))

	if req.Type != nil {
		query.Where(stravaactivity.ActivityType(int(req.GetType())))
	}

	orderSpecifier := utils.AsOrderSpecifier(pageable.Sort, pageable.Direction)

	if pageable.UnPaged {
		data, err := query.Order(orderSpecifier).All(ctx)
		return data, int64(len(data)), err
	}

	count, err := query.Count(ctx)
	data, err := query.
		Order(orderSpecifier).
		Limit(int(pageable.Size)).
		Offset(int(util.AsOffset(pageable.Page, pageable.Size))).
		All(ctx)
	return data, int64(count), err
}

// ExistsByUserIdAndAthleteId is a function that checks if a Strava account exists by user ID and athlete ID
func (s stravaRepository) ExistsByUserIdAndAthleteId(ctx context.Context, userId uuid.UUID, athleteId int64) (bool, error) {
	return s.entClient.StravaAccount.Query().Where(stravaaccount.UserID(userId), stravaaccount.AthleteID(athleteId)).Exist(ctx)
}

// SaveStravaAccount is a function that saves a Strava account
func (s stravaRepository) SaveStravaAccount(ctx context.Context, tx *ent.Tx, req *integration.IntegrateStravaAccountReq) (*ent.StravaAccount, error) {
	data, err := tx.StravaAccount.Create().
		SetUserID(uuid.MustParse(req.UserId)).
		SetAthleteID(req.Strava.AthleteId).
		SetAccessToken(req.Strava.AccessToken).
		SetRefreshToken(req.Strava.RefreshToken).
		SetExpiresAt(time.Unix(req.Strava.ExpiresAt, 0).UTC()).
		SetProfile(req.Strava.Profile).
		SetUsername(req.Strava.Username).
		SetFirstName(req.Strava.FirstName).
		SetLastName(req.Strava.LastName).
		Save(ctx)

	return data, err
}

// GetStravaAccountByUserId is a function that gets a Strava account by user ID
func (s stravaRepository) GetStravaAccountByUserId(ctx context.Context, userId uuid.UUID) (*ent.StravaAccount, error) {
	return s.entClient.StravaAccount.Query().Where(stravaaccount.UserID(userId)).First(ctx)
}

func NewStravaRepository(entClient *ent.Client) IStravaRepository {
	return &stravaRepository{
		entClient: entClient,
	}
}
