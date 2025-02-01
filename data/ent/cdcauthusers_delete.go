// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/cdcauthusers"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/predicate"
)

// CdcAuthUsersDelete is the builder for deleting a CdcAuthUsers entity.
type CdcAuthUsersDelete struct {
	config
	hooks    []Hook
	mutation *CdcAuthUsersMutation
}

// Where appends a list predicates to the CdcAuthUsersDelete builder.
func (caud *CdcAuthUsersDelete) Where(ps ...predicate.CdcAuthUsers) *CdcAuthUsersDelete {
	caud.mutation.Where(ps...)
	return caud
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (caud *CdcAuthUsersDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, caud.sqlExec, caud.mutation, caud.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (caud *CdcAuthUsersDelete) ExecX(ctx context.Context) int {
	n, err := caud.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (caud *CdcAuthUsersDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(cdcauthusers.Table, sqlgraph.NewFieldSpec(cdcauthusers.FieldID, field.TypeUUID))
	if ps := caud.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, caud.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	caud.mutation.done = true
	return affected, err
}

// CdcAuthUsersDeleteOne is the builder for deleting a single CdcAuthUsers entity.
type CdcAuthUsersDeleteOne struct {
	caud *CdcAuthUsersDelete
}

// Where appends a list predicates to the CdcAuthUsersDelete builder.
func (caudo *CdcAuthUsersDeleteOne) Where(ps ...predicate.CdcAuthUsers) *CdcAuthUsersDeleteOne {
	caudo.caud.mutation.Where(ps...)
	return caudo
}

// Exec executes the deletion query.
func (caudo *CdcAuthUsersDeleteOne) Exec(ctx context.Context) error {
	n, err := caudo.caud.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{cdcauthusers.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (caudo *CdcAuthUsersDeleteOne) ExecX(ctx context.Context) {
	if err := caudo.Exec(ctx); err != nil {
		panic(err)
	}
}
