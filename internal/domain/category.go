package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	FindAll(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CategoryService interface {
	CreateCategory(ctx context.Context, name, description string) (*Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
	GetAllCategories(ctx context.Context) ([]*Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name, description string) (*Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}
