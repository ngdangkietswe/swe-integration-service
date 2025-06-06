// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/predicate"
	"github.com/ngdangkietswe/swe-integration-service/data/ent/stravaaccount"
)

// StravaAccountDelete is the builder for deleting a StravaAccount entity.
type StravaAccountDelete struct {
	config
	hooks    []Hook
	mutation *StravaAccountMutation
}

// Where appends a list predicates to the StravaAccountDelete builder.
func (sad *StravaAccountDelete) Where(ps ...predicate.StravaAccount) *StravaAccountDelete {
	sad.mutation.Where(ps...)
	return sad
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sad *StravaAccountDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, sad.sqlExec, sad.mutation, sad.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (sad *StravaAccountDelete) ExecX(ctx context.Context) int {
	n, err := sad.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sad *StravaAccountDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(stravaaccount.Table, sqlgraph.NewFieldSpec(stravaaccount.FieldID, field.TypeUUID))
	if ps := sad.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, sad.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	sad.mutation.done = true
	return affected, err
}

// StravaAccountDeleteOne is the builder for deleting a single StravaAccount entity.
type StravaAccountDeleteOne struct {
	sad *StravaAccountDelete
}

// Where appends a list predicates to the StravaAccountDelete builder.
func (sado *StravaAccountDeleteOne) Where(ps ...predicate.StravaAccount) *StravaAccountDeleteOne {
	sado.sad.mutation.Where(ps...)
	return sado
}

// Exec executes the deletion query.
func (sado *StravaAccountDeleteOne) Exec(ctx context.Context) error {
	n, err := sado.sad.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{stravaaccount.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sado *StravaAccountDeleteOne) ExecX(ctx context.Context) {
	if err := sado.Exec(ctx); err != nil {
		panic(err)
	}
}
