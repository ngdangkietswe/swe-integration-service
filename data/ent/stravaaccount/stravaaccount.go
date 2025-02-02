// Code generated by ent, DO NOT EDIT.

package stravaaccount

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the stravaaccount type in the database.
	Label = "strava_account"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldAthleteID holds the string denoting the athlete_id field in the database.
	FieldAthleteID = "athlete_id"
	// FieldAccessToken holds the string denoting the access_token field in the database.
	FieldAccessToken = "access_token"
	// FieldRefreshToken holds the string denoting the refresh_token field in the database.
	FieldRefreshToken = "refresh_token"
	// FieldExpiresAt holds the string denoting the expires_at field in the database.
	FieldExpiresAt = "expires_at"
	// FieldProfile holds the string denoting the profile field in the database.
	FieldProfile = "profile"
	// FieldUsername holds the string denoting the username field in the database.
	FieldUsername = "username"
	// FieldFirstName holds the string denoting the first_name field in the database.
	FieldFirstName = "first_name"
	// FieldLastName holds the string denoting the last_name field in the database.
	FieldLastName = "last_name"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeCdcAuthUsers holds the string denoting the cdc_auth_users edge name in mutations.
	EdgeCdcAuthUsers = "cdc_auth_users"
	// Table holds the table name of the stravaaccount in the database.
	Table = "strava_account"
	// CdcAuthUsersTable is the table that holds the cdc_auth_users relation/edge.
	CdcAuthUsersTable = "strava_account"
	// CdcAuthUsersInverseTable is the table name for the CdcAuthUsers entity.
	// It exists in this package in order to avoid circular dependency with the "cdcauthusers" package.
	CdcAuthUsersInverseTable = "cdc_auth_users"
	// CdcAuthUsersColumn is the table column denoting the cdc_auth_users relation/edge.
	CdcAuthUsersColumn = "user_id"
)

// Columns holds all SQL columns for stravaaccount fields.
var Columns = []string{
	FieldID,
	FieldUserID,
	FieldAthleteID,
	FieldAccessToken,
	FieldRefreshToken,
	FieldExpiresAt,
	FieldProfile,
	FieldUsername,
	FieldFirstName,
	FieldLastName,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// AccessTokenValidator is a validator for the "access_token" field. It is called by the builders before save.
	AccessTokenValidator func(string) error
	// RefreshTokenValidator is a validator for the "refresh_token" field. It is called by the builders before save.
	RefreshTokenValidator func(string) error
	// ProfileValidator is a validator for the "profile" field. It is called by the builders before save.
	ProfileValidator func(string) error
	// UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	UsernameValidator func(string) error
	// FirstNameValidator is a validator for the "first_name" field. It is called by the builders before save.
	FirstNameValidator func(string) error
	// LastNameValidator is a validator for the "last_name" field. It is called by the builders before save.
	LastNameValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the StravaAccount queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUserID orders the results by the user_id field.
func ByUserID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserID, opts...).ToFunc()
}

// ByAthleteID orders the results by the athlete_id field.
func ByAthleteID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAthleteID, opts...).ToFunc()
}

// ByAccessToken orders the results by the access_token field.
func ByAccessToken(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAccessToken, opts...).ToFunc()
}

// ByRefreshToken orders the results by the refresh_token field.
func ByRefreshToken(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRefreshToken, opts...).ToFunc()
}

// ByExpiresAt orders the results by the expires_at field.
func ByExpiresAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpiresAt, opts...).ToFunc()
}

// ByProfile orders the results by the profile field.
func ByProfile(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProfile, opts...).ToFunc()
}

// ByUsername orders the results by the username field.
func ByUsername(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUsername, opts...).ToFunc()
}

// ByFirstName orders the results by the first_name field.
func ByFirstName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFirstName, opts...).ToFunc()
}

// ByLastName orders the results by the last_name field.
func ByLastName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastName, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByCdcAuthUsersField orders the results by cdc_auth_users field.
func ByCdcAuthUsersField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCdcAuthUsersStep(), sql.OrderByField(field, opts...))
	}
}
func newCdcAuthUsersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CdcAuthUsersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, CdcAuthUsersTable, CdcAuthUsersColumn),
	)
}
