package paginator

import (
	"fmt"
)

// PaginatorRepositoryFacilitator is the facilitator that will add pagination handling capabilities to the repository.
type PaginatorRepositoryFacilitator struct {
	dbType string
}

// NewPaginatorRepositoryFacilitator creates a new instance of the facilitator.
func NewPaginatorRepositoryFacilitator(dbType string) *PaginatorRepositoryFacilitator {

	return &PaginatorRepositoryFacilitator{
		dbType: dbType,
	}
}

// WithPagination attaches the pagination clause to the query.
func (repo *PaginatorRepositoryFacilitator) WithPagination(query string, p Paginator) string {

	return fmt.Sprintf("%s LIMIT %d OFFSET %d", query, p.Size, (p.Page-1)*p.Size)
}
