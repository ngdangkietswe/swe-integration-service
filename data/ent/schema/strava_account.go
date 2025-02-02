package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// StravaAccount holds the schema definition for the StravaAccount entity.
type StravaAccount struct {
	ent.Schema
}

// Fields of the StravaAccount.
func (StravaAccount) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user_id", uuid.UUID{}),
		field.Int64("athlete_id"),
		field.String("access_token").NotEmpty(),
		field.String("refresh_token").NotEmpty(),
		field.Time("expires_at"),
		field.String("profile").NotEmpty(),
		field.String("username").NotEmpty(),
		field.String("first_name").NotEmpty(),
		field.String("last_name").NotEmpty(),
		field.Time("created_at").Immutable().Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
	}
}

// Edges of the StravaAccount.
func (StravaAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cdc_auth_users", CdcAuthUsers.Type).
			Ref("strava_accounts").
			Unique().
			Required().
			Field("user_id"),
	}
}

func (StravaAccount) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "strava_account",
		},
	}
}
