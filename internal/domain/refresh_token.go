package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshToken lưu trong DB để có thể thu hồi (revoke) khi logout / đổi thiết bị.
// Quan trọng cho multi-platform: mỗi thiết bị (web, iOS, Android) có 1 refresh token riêng,
// nên user có thể "đăng xuất trên tất cả thiết bị" hoặc chỉ 1 thiết bị.
type RefreshToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	TokenHash  string    `gorm:"not null" json:"-"`
	DeviceInfo string    `json:"device_info"` // vd: "web-chrome", "ios-app", "android-app"
	ExpiresAt  time.Time `json:"expires_at"`
	Revoked    bool      `gorm:"default:false" json:"revoked"`
	CreatedAt  time.Time `json:"created_at"`
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeByID(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}
