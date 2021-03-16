package filter

import (
	"reflect"
)

// FilterControllerFacilitator is the facilitator that will add filter handling capabilities to the controller.
type FilterControllerFacilitator struct{}

// NewFilterControllerFacilitator creates a new instance of the facilitator.
func NewFilterControllerFacilitator() *FilterControllerFacilitator {

	return &FilterControllerFacilitator{}
}

// GetFilters return the struct passed in to it as a slice of filters.
//
// The 'data' parameter should always be a struct.
func (ctl *FilterControllerFacilitator) GetFilters(data interface{}) (filters []Filter, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	// populate filters slice using data
	elem := reflect.ValueOf(data).Elem()
	elemType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {

		f := elem.Field(i)

		// prevent ignoring of the filter if it is of the type `bool`,
		// in which case both `true` and `false`(zero value for bool type)
		// need to be captured as valid values for the filter.
		if f.Kind() == reflect.Bool {

			filters = append(filters, Filter{
				Name:  elemType.Field(i).Name,
				Value: f.Interface(),
			})

			continue
		}

		if !f.IsZero() {

			filters = append(filters, Filter{
				Name:  elemType.Field(i).Name,
				Value: f.Interface(),
			})
		}
	}

	return filters, nil
}
