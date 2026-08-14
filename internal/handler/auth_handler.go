package handler

import (
	"net/http"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/config"
	"yt-clone-server/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg         *config.Config
	authService *service.AuthService
}

func NewAuthHandler(cfg *config.Config, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{cfg: cfg, authService: authService}
}

// GET /api/auth/google/login?state=xxx
// Dùng cho WEB: Next.js gọi endpoint này (hoặc redirect thẳng người dùng tới đây),
// backend trả về URL Google, browser redirect sang Google để đăng nhập.
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		state = "default-state"
	}
	url := h.authService.GetGoogleAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"auth_url": url})
}

// GET /api/auth/google/callback?code=xxx
// Google redirect trình duyệt về đây sau khi user đồng ý đăng nhập.
// Backend đổi code lấy token, tạo user, rồi redirect tiếp về Next.js kèm token qua query param.
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		appErr := apperrors.NewBadRequest("Missing code parameter")
		c.JSON(appErr.Status(), appErr)
		return
	}

	tokens, err := h.authService.HandleGoogleCallback(c.Request.Context(), code, "web")
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	// Redirect về Next.js kèm token. Production nên cân nhắc dùng httpOnly cookie
	// hoặc 1 mã tạm (one-time code) thay vì đặt thẳng token lên URL để tránh lộ qua log/lịch sử trình duyệt.
	redirectURL := h.cfg.FrontendCallbackURL +
		"?access_token=" + tokens.AccessToken +
		"&refresh_token=" + tokens.RefreshToken

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)

}

// POST /api/auth/google/mobile
// Body: { "id_token": "..." }
// Dùng cho MOBILE (iOS/Android): app native dùng Google Sign-In SDK để lấy id_token
// trên máy, rồi gửi thẳng lên đây để backend verify và cấp JWT riêng của hệ thống.
func (h *AuthHandler) GoogleMobileLogin(c *gin.Context) {
	var req struct {
		IDToken    string `json:"id_token" binding:"required"`
		DeviceInfo string `json:"device_info"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("thiếu id_token")
		c.JSON(appErr.Status(), appErr)
		return
	}
	if req.DeviceInfo == "" {
		req.DeviceInfo = "mobile"
	}

	tokens, err := h.authService.VerifyGoogleIDToken(c.Request.Context(), req.IDToken, req.DeviceInfo)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// POST /api/auth/refresh
// Body: { "refresh_token": "..." }
// Dùng chung cho web và mobile để lấy access token mới khi cái cũ hết hạn.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("Missing refresh token")
		c.JSON(appErr.Status(), appErr)
		return
	}

	tokens, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("Missing refresh token")
		c.JSON(appErr.Status(), appErr)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
