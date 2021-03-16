package paginator

// PaginatorControllerFacilitator is the facilitator that will add pagination handling capabilities to the controller.
type PaginatorControllerFacilitator struct{}

// NewPaginatorControllerFacilitator creates a new instance of the facilitator.
func NewPaginatorControllerFacilitator() *PaginatorControllerFacilitator {

	return &PaginatorControllerFacilitator{}
}

// getPaginator extracts pagination data from query parameters
func (ctl *PaginatorControllerFacilitator) getPaginator(page, size uint32) Paginator {

	// create default paginator
	paginator := NewPaginator()

	paginator.Page = page
	paginator.Size = size

	return paginator
}
