package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name             string    `gorm:"not null" json:"name"`
	Handle           string    `gorm:"uniqueIndex;not null" json:"handle"`
	AvatarURL        string    `json:"avatar_url"`
	BannerURL        string    `json:"banner_url"`
	Description      string    `json:"description"`
	SubscribersCount int64     `gorm:"default:0" json:"subscribers_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type ChannelRepository interface {
	Create(ctx context.Context, Channel *Channel) error
	FindByID(ctx context.Context, id uuid.UUID) (*Channel, error)
	FindByHandle(ctx context.Context, handle string) (*Channel, error)
	FindAllByUserID(ctx context.Context, UserID uuid.UUID) ([]*Channel, error)
	Update(ctx context.Context, channel *Channel) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChannelService interface {
	CreateChannel(ctx context.Context, userID uuid.UUID, name, handle, description, avatarURL string) (*Channel, error)
	GetMyChannels(ctx context.Context, UserID uuid.UUID) ([]*Channel, error)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*Channel, error)
	GetChannelByHandle(ctx context.Context, handle string) (*Channel, error)
	UpdateChannel(ctx context.Context, userID, channelID uuid.UUID, name, handle, description, avatarURL, bannerURL string) (*Channel, error)
	DeleteChannel(ctx context.Context, userID, channelID uuid.UUID) error
}
