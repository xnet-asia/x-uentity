package domain

// Entity is the base interface for all domain entities
type Entity interface {
	GetID() string
}

// BaseEntity provides common fields for all entities
type BaseEntity struct {
	ID        string
	CreatedAt int64
	UpdatedAt int64
}

func (e *BaseEntity) GetID() string {
	return e.ID
}
