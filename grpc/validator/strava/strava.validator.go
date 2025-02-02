package strava

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-integration-service/data/repository/strava"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

type stravaValidator struct {
	stravaRepo strava.IStravaRepository
}

// ValidateIntegrateStravaAccountRequest is a function that validates the request to integrate a Strava account
func (s stravaValidator) ValidateIntegrateStravaAccountRequest(ctx context.Context, request *integration.IntegrateStravaAccountReq) error {
	// TODO: Implement the validation logic for fields in the request (required)
	exists, err := s.stravaRepo.ExistsByUserIdAndAthleteId(ctx, uuid.MustParse(request.UserId), request.Strava.AthleteId)
	if err != nil {
		return err
	} else if exists {
		return errors.New("strava account already exists")
	}

	return nil
}

func NewStravaValidator(stravaRepo strava.IStravaRepository) IStravaValidator {
	return &stravaValidator{
		stravaRepo: stravaRepo,
	}
}
