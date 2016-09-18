package blueprint

// Model represents a data model used by resource controllers to automate
// RESTful CRUD operations.
type Model interface {
	New() Model
	PathID() string
	SetID(id string) error
	FindAll() ([]Model, error)
	FindByID(id string) error
	FindAllByOwner(owner Model) ([]Model, error)
	BelongsTo(owner Model) error
	Save() error
	Delete() error
	Validate() error
}
