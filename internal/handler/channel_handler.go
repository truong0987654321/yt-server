package handler

import (
	"net/http"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/domain"
	"yt-clone-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChannelHandler struct {
	channelService domain.ChannelService
}

func NewChannelHandler(channelService domain.ChannelService) *ChannelHandler {
	return &ChannelHandler{channelService: channelService}
}

type createChannelRequest struct {
	Name        string `json:"name" binding:"required"`
	Handle      string `json:"handle" binding:"required"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

type updateChannelRequest struct {
	Name        string `json:"name"`
	Handle      string `json:"handle"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	BannerURL   string `json:"banner_url"`
}

// POST /api/channels
func (h *ChannelHandler) Create(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		appErr := apperrors.NewAuthorization("Unauthenticated")
		c.JSON(appErr.Status(), appErr)
		return
	}
	userID := userIDVal.(uuid.UUID)
	var req createChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("Name and Handle are required")
		c.JSON(appErr.Status(), appErr)
		return
	}

	channel, err := h.channelService.CreateChannel(c.Request.Context(), userID, req.Name, req.Handle, req.Description, req.AvatarURL)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// GET /api/channels/me
func (h *ChannelHandler) GetMyChannels(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		appErr := apperrors.NewAuthorization("Unauthenticated")
		c.JSON(appErr.Status(), appErr)
		return
	}
	userID := userIDVal.(uuid.UUID)

	channels, err := h.channelService.GetMyChannels(c.Request.Context(), userID)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, channels)
}

// GET /api/channels/:id
func (h *ChannelHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("Invalid channel ID format")
		c.JSON(appErr.Status(), appErr)
		return
	}

	channel, err := h.channelService.GetChannelByID(c.Request.Context(), id)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, channel)
}

// GET /api/channels/handle/:handle
func (h *ChannelHandler) GetByHandle(c *gin.Context) {
	handle := c.Param("handle")
	channel, err := h.channelService.GetChannelByHandle(c.Request.Context(), handle)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, channel)
}

// PUT /api/channels/:id
func (h *ChannelHandler) Update(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		appErr := apperrors.NewAuthorization("Unauthenticated")
		c.JSON(appErr.Status(), appErr)
		return
	}
	userID := userIDVal.(uuid.UUID)

	idStr := c.Param("id")
	channelID, err := uuid.Parse(idStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("Invalid channel ID format")
		c.JSON(appErr.Status(), appErr)
		return
	}

	var req updateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("Invalid request payload")
		c.JSON(appErr.Status(), appErr)
		return
	}

	channel, err := h.channelService.UpdateChannel(
		c.Request.Context(),
		userID,
		channelID,
		req.Name,
		req.Handle,
		req.Description,
		req.AvatarURL,
		req.BannerURL,
	)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, channel)
}

// DELETE /api/channels/:id
func (h *ChannelHandler) Delete(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		appErr := apperrors.NewAuthorization("Unauthenticated")
		c.JSON(appErr.Status(), appErr)
		return
	}
	userID := userIDVal.(uuid.UUID)

	idStr := c.Param("id")
	channelID, err := uuid.Parse(idStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("Invalid channel ID format")
		c.JSON(appErr.Status(), appErr)
		return
	}

	if err := h.channelService.DeleteChannel(
		c.Request.Context(),
		userID,
		channelID,
	); err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}
