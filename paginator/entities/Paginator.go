package entities

// Paginator entity
type Paginator struct {
	Page uint32
	Size uint32
}

// NewPaginator creates a new paginator object with default values.
func NewPaginator() Paginator {

	return Paginator{
		Page: 1,
		Size: 10,
	}
}
