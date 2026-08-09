package postgres

import (
	"context"
	"errors"
	"yt-clone-server/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(token).Error
}
func (r *refreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var t domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND revoked = false", tokenHash).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *refreshTokenRepository) RevokeByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}
func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
