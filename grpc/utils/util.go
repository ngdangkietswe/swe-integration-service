package utils

import (
	"entgo.io/ent/dialect/sql"
	"github.com/ngdangkietswe/swe-integration-service/data/ent"
	"github.com/samber/lo"
)

// AsOrderSpecifier is a function that returns an order specifier based on the sort and order.
func AsOrderSpecifier(sort, direction string) func(selector *sql.Selector) {
	return lo.Ternary(direction == "asc", ent.Asc(sort), ent.Desc(sort))
}
