package strava

import (
	"context"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/stravaaccount"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"time"
)

type stravaRepository struct {
	entClient *ent.Client
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
