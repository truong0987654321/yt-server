package postgres

import (
	"context"
	"errors"
	"yt-clone-server/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var c domain.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Category{}, "id = ?", id).Error
}
