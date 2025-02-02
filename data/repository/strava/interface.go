package strava

import (
	"context"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type IStravaRepository interface {
	SaveStravaAccount(ctx context.Context, req *integration.IntegrateStravaAccountReq) (*ent.StravaAccount, error)
	GetStravaAccountByUserId(ctx context.Context, userId uuid.UUID) (*ent.StravaAccount, error)
	ExistsByUserIdAndAthleteId(ctx context.Context, userId uuid.UUID, athleteId int64) (bool, error)
}
