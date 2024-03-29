package filter

import (
	"fmt"
	"reflect"

	"github.com/kosatnkn/req"
)

// extendedFilter extension of the normal filter object with additional repository related fields.
type extendedFilter struct {
	Filter
	Field    string
	Operator string
}

// FilterRepositoryFacilitator is the facilitator that will add filter handling capabilities to the repository.
type FilterRepositoryFacilitator struct {
	dbType    string
	filterMap map[string][]string
}

// NewFilterRepositoryFacilitator creates a new instance of the facilitator.
func NewFilterRepositoryFacilitator(dbType string, filterMap map[string][]string) *FilterRepositoryFacilitator {
	return &FilterRepositoryFacilitator{
		dbType:    dbType,
		filterMap: filterMap,
	}
}

// WithFilters attaches filters as a WHERE clause to the query.
func (repo *FilterRepositoryFacilitator) WithFilters(query string, fts []Filter) (string, map[string]interface{}, error) {
	params := make(map[string]interface{})

	if len(fts) == 0 {
		return query, params, nil
	}

	var where string

	efs, err := repo.extendFilters(fts)
	if err != nil {
		return query, params, err
	}
	if len(efs) == 0 {
		return query, params, nil
	}

	for _, f := range efs {
		qp, vs := repo.getConditionQueryPart(f)

		// attach query part to the where clause
		where += qp

		// add params needed for the query part
		for k, v := range vs {
			params[k] = v
		}
	}

	return fmt.Sprintf("%s WHERE%s", query, where[4:]), params, nil
}

// extendFilters sets additional filter parameters like table field and operator for filters.
func (repo *FilterRepositoryFacilitator) extendFilters(filters []Filter) ([]extendedFilter, error) {
	efs := make([]extendedFilter, 0)

	if repo.filterMap == nil {
		return efs, fmt.Errorf("filter-repo-facilitator: a filter mapping is not declared")
	}

	for _, filter := range filters {
		fm := repo.filterMap[filter.Name]
		if len(fm) == 0 {
			continue
		}

		efs = append(efs, extendedFilter{
			Filter:   filter,
			Field:    fm[0],
			Operator: repo.getOperatorFor(filter.Name),
		})
	}

	return efs, nil
}

// getOperatorFor returns the operator from field mapping if one is set, otherwise
// will return 'SelectEqual' as the default.
func (repo *FilterRepositoryFacilitator) getOperatorFor(name string) string {
	m := repo.filterMap[name]
	if len(m) == 1 {
		return req.SelectEqual
	}

	return m[1]
}

// getConditionQueryPart returns the query part needed to add the filter condition to the query.
func (repo *FilterRepositoryFacilitator) getConditionQueryPart(f extendedFilter) (string, map[string]interface{}) {
	switch f.Operator {
	case req.SelectLike:
		return repo.getSelectLikeQueryPart(f)
	case req.SelectIn:
		return repo.getSelectInQueryPart(f)
	default:
		return repo.getSelectEqualQueryPart(f)
	}
}

// getSelectEqualQueryPart creates the query part for an 'equal' operation.
//
// ex: AND `field` = `value`
func (repo *FilterRepositoryFacilitator) getSelectEqualQueryPart(f extendedFilter) (string, map[string]interface{}) {
	m := make(map[string]interface{})
	m[f.Name] = f.Value

	return fmt.Sprintf(" AND %s %s ?%s", f.Field, f.Operator, f.Name), m
}

// getSelectLikeQueryPart creates the query part for a 'like' operation.
//
// ex: AND `field` LIKE `%value%`
func (repo *FilterRepositoryFacilitator) getSelectLikeQueryPart(f extendedFilter) (string, map[string]interface{}) {
	m := make(map[string]interface{})
	m[f.Name] = fmt.Sprintf("%%%s%%", f.Value)

	return fmt.Sprintf(" AND %s %s ?%s", f.Field, f.Operator, f.Name), m
}

// getSelectInQueryPart creates the query part for an 'in' operation.
//
// ex: AND `field` IN (`value1`, `value2`, `value3`)
func (repo *FilterRepositoryFacilitator) getSelectInQueryPart(f extendedFilter) (string, map[string]interface{}) {
	m := make(map[string]interface{})

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
