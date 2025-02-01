package mapper

import (
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
)

// AsMonoStravaAccount is a function that maps an ent.StravaAccount to an integration.StravaAccount.
func AsMonoStravaAccount(entSA *ent.StravaAccount) *integration.StravaAccount {
	return &integration.StravaAccount{
		Id:           entSA.ID.String(),
		UserId:       entSA.UserID.String(),
		AthleteId:    entSA.AthleteID,
		AccessToken:  entSA.AccessToken,
		RefreshToken: entSA.RefreshToken,
		ExpiresAt:    util.Format(&entSA.ExpiresAt, util.LayoutISOWithTime),
		CreatedAt:    util.Format(&entSA.CreatedAt, util.LayoutISOWithTime),
		UpdatedAt:    util.Format(&entSA.UpdatedAt, util.LayoutISOWithTime),
		TokenType:    "Bearer",
		Username:     entSA.Username,
		FirstName:    entSA.FirstName,
		LastName:     entSA.LastName,
	}
}
