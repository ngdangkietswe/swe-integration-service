package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"time"
)

// StravaActivity holds the schema definition for the StravaActivity entity.
type StravaActivity struct {
	ent.Schema
}

// Fields of the StravaActivity.
func (StravaActivity) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Int64("strava_activity_id"),
		field.Int64("athlete_id"),
		field.UUID("user_id", uuid.UUID{}),
		field.String("activity_name").NotEmpty(),
		field.Int("activity_type").Default(0),
		field.String("activity_url").NotEmpty(),
		field.Time("start_date"),
		field.Float("distance").Default(0),
		field.Int32("moving_time").Default(0),
		field.Int32("elapsed_time").Default(0),
		field.Float("total_elevation_gain").Default(0),
		field.Float("average_speed").Default(0),
		field.Float("max_speed").Default(0),
		field.Time("created_at").Immutable().Default(time.Now()),
	}
}

// Edges of the StravaActivity.
func (StravaActivity) Edges() []ent.Edge {
	return []ent.Edge{
		util.One2ManyInverseRequired("cdc_auth_users", CdcAuthUsers.Type, "strava_activities", "user_id"),
	}
}

func (StravaActivity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "strava_activity",
		},
	}
}
