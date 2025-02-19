package strava

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

// SyncStravaActivities is a function that syncs Strava activities with the user's account.
func (s stravaRepository) SyncStravaActivities(ctx context.Context, userId uuid.UUID, athleteId int64, stravaActivities []map[string]interface{}) error {
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

	return s.entClient.StravaActivity.CreateBulk(stravaActivityCreates...).Exec(ctx)
}

// GetListStravaActivitiesByUserId is a function that gets a list of Strava activities by user ID.
func (s stravaRepository) GetListStravaActivitiesByUserId(ctx context.Context, userId uuid.UUID) ([]*ent.StravaActivity, error) {
	return s.entClient.StravaActivity.Query().Where(stravaactivity.UserID(userId)).All(ctx)
}

// GetListStravaActivities is a function that gets a list of Strava activities by user ID.
func (s stravaRepository) GetListStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq, userId uuid.UUID, pageable *common.Pageable) ([]*ent.StravaActivity, int64, error) {
	entSAs := s.entClient.StravaActivity.Query().Where(stravaactivity.UserID(userId))

	if req.Type != nil {
		entSAs.Where(stravaactivity.ActivityType(int(req.GetType())))
	}

	count, err := entSAs.Count(ctx)
	data, err := entSAs.Limit(int(pageable.Size)).Offset(int(utils.AsOffset(pageable.Page, pageable.Size))).All(ctx)
	return data, int64(count), err
}

// ExistsByUserIdAndAthleteId is a function that checks if a Strava account exists by user ID and athlete ID
func (s stravaRepository) ExistsByUserIdAndAthleteId(ctx context.Context, userId uuid.UUID, athleteId int64) (bool, error) {
	return s.entClient.StravaAccount.Query().Where(stravaaccount.UserID(userId), stravaaccount.AthleteID(athleteId)).Exist(ctx)
}

// SaveStravaAccount is a function that saves a Strava account
func (s stravaRepository) SaveStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*ent.StravaAccount, error) {
	data, err := s.entClient.StravaAccount.Create().
		SetUserID(uuid.MustParse(req.UserId)).
		SetAthleteID(req.Strava.AthleteId).
		SetAccessToken(req.Strava.AccessToken).
		SetRefreshToken(req.Strava.RefreshToken).
		SetExpiresAt(time.UnixMilli(req.Strava.ExpiresAt)).
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
