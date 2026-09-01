package postgres

import (
	"context"
	"errors"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) domain.ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) Create(ctx context.Context, channel *domain.Channel) error {
	if channel.ID == uuid.Nil {
		channel.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(channel).Error
}

func (r *channelRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Channel, error) {
	var c domain.Channel
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).Take(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *channelRepository) FindByHandle(ctx context.Context, handle string) (*domain.Channel, error) {
	var c domain.Channel
	err := r.db.WithContext(ctx).Preload("User").Where("handle = ?", handle).Take(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *channelRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Channel, error) {
	var channels []*domain.Channel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (r *channelRepository) Update(ctx context.Context, channel *domain.Channel) error {
	return r.db.WithContext(ctx).Save(channel).Error
}

func (r *channelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Channel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperrors.NewNotFound("channel", id.String())
	}
	return nil
}
