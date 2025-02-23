package strava

import (
	"context"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type IStravaRepository interface {
	SaveStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*ent.StravaAccount, error)
	GetStravaAccountByUserId(ctx context.Context, userId uuid.UUID) (*ent.StravaAccount, error)
	ExistsByUserIdAndAthleteId(ctx context.Context, userId uuid.UUID, athleteId int64) (bool, error)
	RemoveStravaAccountByUserId(ctx context.Context, stravaAccountId uuid.UUID) error

	ExistsStravaActivityById(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsAllStravaActivitiesByIdIn(ctx context.Context, ids []uuid.UUID) (bool, error)
	GetStravaActivityById(ctx context.Context, id uuid.UUID) (*ent.StravaActivity, error)
	GetListStravaActivitiesByIdIn(ctx context.Context, ids []uuid.UUID) ([]*ent.StravaActivity, error)
	GetListStravaActivities(ctx context.Context, req *integration.GetStravaActivitiesReq, userId uuid.UUID, pageable *common.Pageable) ([]*ent.StravaActivity, int64, error)
	GetListStravaActivitiesByUserId(ctx context.Context, userId uuid.UUID) ([]*ent.StravaActivity, error)
	DeleteStravaActivityById(ctx context.Context, id uuid.UUID) error
	DeleteStravaActivitiesByIdIn(ctx context.Context, ids []uuid.UUID) error

	SyncStravaActivities(ctx context.Context, userId uuid.UUID, athleteId int64, stravaActivities []map[string]interface{}) error
	UpdateTokenStravaAccount(ctx context.Context, id uuid.UUID, accessToken string, refreshToken string, expiresAt int64) error
}
