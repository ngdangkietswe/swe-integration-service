package mapper

import (
	"fmt"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"github.com/samber/lo"
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

// AsMonoStravaActivity is a function that maps an ent.StravaActivity to an integration.StravaActivity.
func AsMonoStravaActivity(entSA *ent.StravaActivity) *integration.StravaActivity {
	return &integration.StravaActivity{
		Id:                 entSA.ID.String(),
		StravaActivityId:   entSA.StravaActivityID,
		AthleteId:          entSA.AthleteID,
		UserId:             entSA.UserID.String(),
		ActivityName:       entSA.ActivityName,
		ActivityType:       int32(entSA.ActivityType),
		ActivityUrl:        entSA.ActivityURL,
		Distance:           fmt.Sprintf("%.2f km", entSA.Distance/1000),
		StartDate:          util.Format(&entSA.StartDate, util.LayoutISOWithTime),
		MovingTime:         secondsToTime(entSA.MovingTime),
		ElapsedTime:        secondsToTime(entSA.ElapsedTime),
		TotalElevationGain: entSA.TotalElevationGain,
		AverageSpeed:       entSA.AverageSpeed,
		MaxSpeed:           entSA.MaxSpeed,
		CreatedAt:          util.Format(&entSA.CreatedAt, util.LayoutISOWithTime),
		Pace:               fmt.Sprintf("%.2f /km", float64(entSA.MovingTime/60)/(entSA.Distance/1000)),
	}
}

// secondsToTime is a function that converts seconds to time (HH:mm:ss).
func secondsToTime(seconds int32) string {
	return fmt.Sprintf("%02d:%02d:%02d", seconds/3600, (seconds%3600)/60, seconds%60)
}

// AsListStravaActivity is a function that maps a list of ent.StravaActivity to a list of integration.StravaActivity.
func AsListStravaActivity(entSAs []*ent.StravaActivity) []*integration.StravaActivity {
	var res []*integration.StravaActivity
	lo.ForEach(entSAs, func(entSA *ent.StravaActivity, _ int) {
		res = append(res, AsMonoStravaActivity(entSA))
	})
	return res
}
