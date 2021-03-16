package facilitators

import (
	"github.com/kosatnkn/req/paginator/entities"
)

// PaginatorControllerFacilitator is the facilitator that will add pagination handling capabilities to the controller.
type PaginatorControllerFacilitator struct{}

// NewPaginatorControllerFacilitator creates a new instance of the facilitator.
func NewPaginatorControllerFacilitator() *PaginatorControllerFacilitator {

	return &PaginatorControllerFacilitator{}
}

// getPaginator extracts pagination data from query parameters
func (ctl *PaginatorControllerFacilitator) getPaginator(page, size uint32) entities.Paginator {

	// create default paginator
	paginator := entities.NewPaginator()

	paginator.Page = page
	paginator.Size = size

	return paginator
}
