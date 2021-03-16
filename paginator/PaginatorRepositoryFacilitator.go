package paginator

import (
	"fmt"
)

// PaginatorRepositoryFacilitator is the facilitator that will add pagination handling capabilities to the repository.
type PaginatorRepositoryFacilitator struct {
	db string
}

// NewPaginatorRepositoryFacilitator creates a new instance of the facilitator.
func NewPaginatorRepositoryFacilitator(database string) *PaginatorRepositoryFacilitator {

	return &PaginatorRepositoryFacilitator{
		db: database,
	}
}

// withPagination generates the pagination clause for the query.
func (repo *PaginatorRepositoryFacilitator) withPagination(q string, p Paginator) string {

	return fmt.Sprintf("%s LIMIT %d OFFSET %d", q, p.Size, (p.Page-1)*p.Size)
}
