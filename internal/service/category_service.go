package service

import (
	"context"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/domain"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(ctx context.Context, name, description string) (*domain.Category, error) {
	if name == "" {
		return nil, apperrors.NewBadRequest("Category name cannot be empty")
	}

	category := &domain.Category{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, apperrors.NewNotFound("category", id.String())
	}

	return category, nil
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, name, description string) (*domain.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, apperrors.NewNotFound("category", id.String())
	}
	if name != "" {
		category.Name = name
	}
	category.Description = description

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return apperrors.NewNotFound("category", id.String())
	}

	return s.categoryRepo.Delete(ctx, id)
}
