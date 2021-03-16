package filter

import (
	"context"
	"fmt"
	"reflect"
)

const (
	selectEqual string = "="
	selectLike  string = "LIKE"
	selectIn    string = "IN"
)

// extendedFilter extension of the normal filter object with additional repository related fields.
type extendedFilter struct {
	Filter
	Field    string
	Operator string
}

// FilterRepositoryFacilitator is the facilitator that will add filter handling capabilities to the repository.
type FilterRepositoryFacilitator struct {
	db        string
	filterMap map[string][]string
}

// NewFilterRepositoryFacilitator creates a new instance of the facilitator.
func NewFilterRepositoryFacilitator(database string, filterMap map[string][]string) *FilterRepositoryFacilitator {

	return &FilterRepositoryFacilitator{
		db:        database,
		filterMap: filterMap,
	}
}

// withFilters generates the WHERE clause for the query.
func (repo *FilterRepositoryFacilitator) withFilters(q string, fts []extendedFilter) (string, map[string]interface{}) {

	params := map[string]interface{}{}

	if len(fts) == 0 {
		return q, params
	}

	var w string

	for _, f := range fts {

		qp, vs := repo.getConditionQueryPart(f)

		w += qp

		for k, v := range vs {
			params[k] = v
		}
	}

	return fmt.Sprintf("%s WHERE%s", q, w[4:]), params
}

// extendFilters sets additional filter parameters like table field and operator for filters.
func (repo *FilterRepositoryFacilitator) extendFilters(ctx context.Context, filters []Filter) []extendedFilter {

	efs := make([]extendedFilter, 0)

	for _, filter := range filters {

		fm := repo.filterMap[filter.Name]

		if len(fm) == 0 {
			// repo.logger.Warn(ctx, fmt.Sprintf("No field mapping found for key '%s'", filter.Name))
			continue
		}

		efs = append(efs, extendedFilter{
			Filter:   filter,
			Field:    fm[0],
			Operator: repo.getOperatorFor(filter.Name),
		})
	}

	return efs
}

// getOperatorFor returns the operator from field mapping if one is set, otherwise
// will return 'selectEqual' as the default.
func (repo *FilterRepositoryFacilitator) getOperatorFor(name string) string {

	m := repo.filterMap[name]

	if len(m) == 1 {
		return selectEqual
	}

	return m[1]
}

// getConditionQueryPart returns the query part needed to add the filter condition to the query.
func (repo *FilterRepositoryFacilitator) getConditionQueryPart(f extendedFilter) (string, map[string]interface{}) {

	switch f.Operator {
	case selectLike:
		return repo.getSelectLikeQueryPart(f)
	case selectIn:
		return repo.getSelectInQueryPart(f)
	default:
		return repo.getSelectEqualQueryPart(f)
	}
}

// getSelectEqualQueryPart creates the query part for an 'equal' operation.
//
// ex: AND `field` = `value`
func (repo *FilterRepositoryFacilitator) getSelectEqualQueryPart(f extendedFilter) (string, map[string]interface{}) {

	m := make(map[string]interface{}, 0)

	m[f.Name] = f.Value
	return fmt.Sprintf(" AND %s %s ?%s", f.Field, f.Operator, f.Name), m
}

// getSelectLikeQueryPart creates the query part for a 'like' operation.
//
// ex: AND `field` LIKE `%value%`
func (repo *FilterRepositoryFacilitator) getSelectLikeQueryPart(f extendedFilter) (string, map[string]interface{}) {

	m := make(map[string]interface{}, 0)

	m[f.Name] = fmt.Sprintf("%%%s%%", f.Value)
	return fmt.Sprintf(" AND %s %s ?%s", f.Field, f.Operator, f.Name), m
}

// getSelectInQueryPart creates the query part for an 'in' operation.
//
// ex: AND `field` IN (`value1`, `value2`, `value3`)
func (repo *FilterRepositoryFacilitator) getSelectInQueryPart(f extendedFilter) (string, map[string]interface{}) {

	m := make(map[string]interface{}, 0)

	// placeholders
	var phs string
	var vs []interface{}

	switch reflect.TypeOf(f.Value).Kind() {

	case reflect.Slice:

		rvs := reflect.ValueOf(f.Value)

		if rvs.Len() == 0 {
			return "", m
		}

		for i := 0; i < rvs.Len(); i++ {
			vs = append(vs, rvs.Index(i).Interface())
		}

	default:
		return "", m
	}

	// iterate through the interface{} slice to build the `in` clause
	for i, v := range vs {

		ph := fmt.Sprintf("%s%d", f.Name, i)
		phs += fmt.Sprintf(",?%s", ph)
		m[ph] = v
	}

	return fmt.Sprintf(" AND %s %s (%s)", f.Field, f.Operator, phs[1:]), m
}
