package service

import (
	"context"
	"strings"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/domain"

	"github.com/google/uuid"
)

type channelService struct {
	channelRepo domain.ChannelRepository
	userRepo    domain.UserRepository
}

func NewChannelService(channelRepo domain.ChannelRepository, userRepo domain.UserRepository) domain.ChannelService {
	return &channelService{
		channelRepo: channelRepo,
		userRepo:    userRepo,
	}
}

func (s *channelService) CreateChannel(ctx context.Context, userID uuid.UUID, name, handle, description, avatarURL string) (*domain.Channel, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperrors.NewBadRequest("Channel name cannot be empty")
	}

	handle = strings.TrimSpace(handle)
	if handle == "" {
		return nil, apperrors.NewBadRequest("Channel name cannot be empty")

	}

	if !strings.HasPrefix(handle, "@") {
		handle = "@" + handle
	}

	// 1. Kiểm tra handle trùng lập
	existingHandle, err := s.channelRepo.FindByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}
	if existingHandle != nil {
		return nil, apperrors.NewBadRequest("Channel handle is already taken")
	}

	// 2. Lấy thông tin user để set default avatar nếu không truyền
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, apperrors.NewNotFound("user", userID.String())
	}

	if avatarURL == "" {
		avatarURL = user.AvatarURL
	}

	channel := &domain.Channel{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Handle:      handle,
		AvatarURL:   avatarURL,
		Description: description,
	}

	if err := s.channelRepo.Create(ctx, channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *channelService) GetMyChannels(ctx context.Context, userID uuid.UUID) ([]*domain.Channel, error) {
	return s.channelRepo.FindAllByUserID(ctx, userID)
}

func (s *channelService) GetChannelByID(ctx context.Context, id uuid.UUID) (*domain.Channel, error) {
	channel, err := s.channelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		return nil, apperrors.NewNotFound("channel", id.String())
	}
	return channel, nil
}

func (s *channelService) GetChannelByHandle(ctx context.Context, handle string) (*domain.Channel, error) {
	if !strings.HasPrefix(handle, "@") {
		handle = "@" + handle
	}
	channel, err := s.channelRepo.FindByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		return nil, apperrors.NewNotFound("channel", handle)
	}
	return channel, nil
}
func (s *channelService) UpdateChannel(ctx context.Context, userID, channelID uuid.UUID, name, handle, description, avatarURL, bannerURL string) (*domain.Channel, error) {
	channel, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		return nil, apperrors.NewNotFound("channel", channelID.String())
	}

	if channel.UserID != userID {
		return nil, apperrors.NewAuthorization("You do not have permission to update this channel")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		channel.Name = name
	}

	handle = strings.TrimSpace(handle)
	if handle != "" {
		if !strings.HasPrefix(handle, "@") {
			handle = "@" + handle
		}
		if handle != channel.Handle {
			existing, err := s.channelRepo.FindByHandle(ctx, handle)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != channelID {
				return nil, apperrors.NewBadRequest("Channel handle is already taken")
			}
			channel.Handle = handle
		}
	}

	channel.Description = description
	if avatarURL != "" {
		channel.AvatarURL = avatarURL
	}
	if bannerURL != "" {
		channel.BannerURL = bannerURL
	}

	if err := s.channelRepo.Update(ctx, channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *channelService) DeleteChannel(ctx context.Context, userID, channelID uuid.UUID) error {
	channel, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		return err
	}
	if channel == nil {
		return apperrors.NewNotFound("channel", channelID.String())
	}

	if channel.UserID != userID {
		return apperrors.NewAuthorization("You do not have permission to delete this channel")
	}

	return s.channelRepo.Delete(ctx, channelID)
}
