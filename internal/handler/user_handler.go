package handler

import (
	"net/http"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/domain"
	"yt-clone-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepo domain.UserRepository
}

func NewUserHandler(uRepo domain.UserRepository) *UserHandler {
	return &UserHandler{userRepo: uRepo}
}

func (h *UserHandler) Me(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)

	if !exists {
		appErr := apperrors.NewAuthorization("Unauthenticated")
		c.JSON(appErr.Status(), appErr)
		return
	}

	userID := userIDVal.(uuid.UUID)
	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		apperr := apperrors.NewNotFound("user", userID.String())
		c.JSON(apperr.Status(), apperr)
		return
	}
	c.JSON(http.StatusOK, user)
}
