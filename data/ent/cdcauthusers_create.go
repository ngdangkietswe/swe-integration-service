// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/cdcauthusers"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/stravaaccount"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/stravaactivity"
)

// CdcAuthUsersCreate is the builder for creating a CdcAuthUsers entity.
type CdcAuthUsersCreate struct {
	config
	mutation *CdcAuthUsersMutation
	hooks    []Hook
}

// SetUsername sets the "username" field.
func (cauc *CdcAuthUsersCreate) SetUsername(s string) *CdcAuthUsersCreate {
	cauc.mutation.SetUsername(s)
	return cauc
}

// SetEmail sets the "email" field.
func (cauc *CdcAuthUsersCreate) SetEmail(s string) *CdcAuthUsersCreate {
	cauc.mutation.SetEmail(s)
	return cauc
}

// SetID sets the "id" field.
func (cauc *CdcAuthUsersCreate) SetID(u uuid.UUID) *CdcAuthUsersCreate {
	cauc.mutation.SetID(u)
	return cauc
}

// AddStravaAccountIDs adds the "strava_accounts" edge to the StravaAccount entity by IDs.
func (cauc *CdcAuthUsersCreate) AddStravaAccountIDs(ids ...uuid.UUID) *CdcAuthUsersCreate {
	cauc.mutation.AddStravaAccountIDs(ids...)
	return cauc
}

// AddStravaAccounts adds the "strava_accounts" edges to the StravaAccount entity.
func (cauc *CdcAuthUsersCreate) AddStravaAccounts(s ...*StravaAccount) *CdcAuthUsersCreate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cauc.AddStravaAccountIDs(ids...)
}

// AddStravaActivityIDs adds the "strava_activities" edge to the StravaActivity entity by IDs.
func (cauc *CdcAuthUsersCreate) AddStravaActivityIDs(ids ...uuid.UUID) *CdcAuthUsersCreate {
	cauc.mutation.AddStravaActivityIDs(ids...)
	return cauc
}

// AddStravaActivities adds the "strava_activities" edges to the StravaActivity entity.
func (cauc *CdcAuthUsersCreate) AddStravaActivities(s ...*StravaActivity) *CdcAuthUsersCreate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cauc.AddStravaActivityIDs(ids...)
}

// Mutation returns the CdcAuthUsersMutation object of the builder.
func (cauc *CdcAuthUsersCreate) Mutation() *CdcAuthUsersMutation {
	return cauc.mutation
}

// Save creates the CdcAuthUsers in the database.
func (cauc *CdcAuthUsersCreate) Save(ctx context.Context) (*CdcAuthUsers, error) {
	return withHooks(ctx, cauc.sqlSave, cauc.mutation, cauc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cauc *CdcAuthUsersCreate) SaveX(ctx context.Context) *CdcAuthUsers {
	v, err := cauc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cauc *CdcAuthUsersCreate) Exec(ctx context.Context) error {
	_, err := cauc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cauc *CdcAuthUsersCreate) ExecX(ctx context.Context) {
	if err := cauc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cauc *CdcAuthUsersCreate) check() error {
	if _, ok := cauc.mutation.Username(); !ok {
		return &ValidationError{Name: "username", err: errors.New(`ent: missing required field "CdcAuthUsers.username"`)}
	}
	if v, ok := cauc.mutation.Username(); ok {
		if err := cdcauthusers.UsernameValidator(v); err != nil {
			return &ValidationError{Name: "username", err: fmt.Errorf(`ent: validator failed for field "CdcAuthUsers.username": %w`, err)}
		}
	}
	if _, ok := cauc.mutation.Email(); !ok {
		return &ValidationError{Name: "email", err: errors.New(`ent: missing required field "CdcAuthUsers.email"`)}
	}
	if v, ok := cauc.mutation.Email(); ok {
		if err := cdcauthusers.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf(`ent: validator failed for field "CdcAuthUsers.email": %w`, err)}
		}
	}
	return nil
}

func (cauc *CdcAuthUsersCreate) sqlSave(ctx context.Context) (*CdcAuthUsers, error) {
	if err := cauc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cauc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cauc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cauc.mutation.id = &_node.ID
	cauc.mutation.done = true
	return _node, nil
}

func (cauc *CdcAuthUsersCreate) createSpec() (*CdcAuthUsers, *sqlgraph.CreateSpec) {
	var (
		_node = &CdcAuthUsers{config: cauc.config}
		_spec = sqlgraph.NewCreateSpec(cdcauthusers.Table, sqlgraph.NewFieldSpec(cdcauthusers.FieldID, field.TypeUUID))
	)
	if id, ok := cauc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cauc.mutation.Username(); ok {
		_spec.SetField(cdcauthusers.FieldUsername, field.TypeString, value)
		_node.Username = value
	}
	if value, ok := cauc.mutation.Email(); ok {
		_spec.SetField(cdcauthusers.FieldEmail, field.TypeString, value)
		_node.Email = value
	}
	if nodes := cauc.mutation.StravaAccountsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   cdcauthusers.StravaAccountsTable,
			Columns: []string{cdcauthusers.StravaAccountsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(stravaaccount.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cauc.mutation.StravaActivitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   cdcauthusers.StravaActivitiesTable,
			Columns: []string{cdcauthusers.StravaActivitiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(stravaactivity.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CdcAuthUsersCreateBulk is the builder for creating many CdcAuthUsers entities in bulk.
type CdcAuthUsersCreateBulk struct {
	config
	err      error
	builders []*CdcAuthUsersCreate
}

// Save creates the CdcAuthUsers entities in the database.
func (caucb *CdcAuthUsersCreateBulk) Save(ctx context.Context) ([]*CdcAuthUsers, error) {
	if caucb.err != nil {
		return nil, caucb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(caucb.builders))
	nodes := make([]*CdcAuthUsers, len(caucb.builders))
	mutators := make([]Mutator, len(caucb.builders))
	for i := range caucb.builders {
		func(i int, root context.Context) {
			builder := caucb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CdcAuthUsersMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, caucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, caucb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, caucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (caucb *CdcAuthUsersCreateBulk) SaveX(ctx context.Context) []*CdcAuthUsers {
	v, err := caucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (caucb *CdcAuthUsersCreateBulk) Exec(ctx context.Context) error {
	_, err := caucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (caucb *CdcAuthUsersCreateBulk) ExecX(ctx context.Context) {
	if err := caucb.Exec(ctx); err != nil {
		panic(err)
	}
}
