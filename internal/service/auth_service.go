package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/config"
	"yt-clone-server/internal/domain"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
)

// ErrGoogleAuthFailed: Google từ chối/verify thất bại (code sai, id_token giả, hết hạn...).
var ErrGoogleAuthFailed = apperrors.NewAuthorization("Google authentication failed")

// ErrInvalidRefreshToken: refresh token không có trong DB, đã bị revoke, hoặc hết hạn.
var ErrInvalidRefreshToken = apperrors.NewAuthorization("The refresh token is invalid or has been revoked")

var ErrCreateUser = apperrors.NewAuthorization("The refresh token is invalid or has been revoked")

// googleUserInfo là dữ liệu Google trả về từ endpoint userinfo.
type googleUserInfo struct {
	Sub     string `json:"sub"` // Google ID duy nhất
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// TokenPair là cặp token trả về cho client (web hoặc mobile) sau khi login thành công.
type TokenPair struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         *domain.User `json:"user"`
}

type AuthService struct {
	cfg         *config.Config
	oauthConfig *oauth2.Config
	userRepo    domain.UserRepository
	refreshRepo domain.RefreshTokenRepository
	jwtService  *JWTService
}

func NewAuthService(cfg *config.Config, userRepo domain.UserRepository, refreshRepo domain.RefreshTokenRepository, jwtService *JWTService) *AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     googleoauth.Endpoint,
	}
	return &AuthService{
		cfg:         cfg,
		oauthConfig: oauthConfig,
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtService:  jwtService,
	}
}

func (s *AuthService) GetGoogleAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code, deviceInfo string) (*TokenPair, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	info, err := s.fetchGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	return s.loginOrRegister(ctx, info, deviceInfo)
}

func (s *AuthService) VerifyGoogleIDToken(ctx context.Context, idToken, deviceInfo string) (*TokenPair, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://oauth2.googleapis.com/tokeninfo?id_token="+idToken, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrGoogleAuthFailed
	}

	var payload struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
		Aud     string `json:"aud"` // phải khớp GoogleClientID (hoặc client id riêng của app mobile)
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	// QUAN TRỌNG: production nên kiểm tra payload.Aud nằm trong danh sách các Client ID
	// hợp lệ (web + iOS + Android đều có Client ID riêng trên Google Cloud Console).

	info := &googleUserInfo{
		Sub:     payload.Sub,
		Email:   payload.Email,
		Name:    payload.Name,
		Picture: payload.Picture,
	}

	return s.loginOrRegister(ctx, info, deviceInfo)
}

func (s *AuthService) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*googleUserInfo, error) {
	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var info googleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGoogleAuthFailed, err)
	}
	return &info, nil
}

// ===== FLOW 2: Mobile (iOS/Android) - Native Google Sign-In SDK =====
// Ở mobile, app dùng Google Sign-In SDK native để lấy id_token trực tiếp trên máy,
// sau đó CHỈ gửi id_token đó lên backend để verify, không cần chạy authorization code flow.
// Đây là lý do server này "sẵn sàng đa nền tảng": chỉ cần thêm 1 endpoint xác minh token.

// loginOrRegister: tìm user theo GoogleID, nếu chưa có thì tạo mới (đăng ký tự động),
// sau đó sinh access + refresh token mới cho thiết bị này.
func (s *AuthService) loginOrRegister(ctx context.Context, info *googleUserInfo, deviceInfo string) (*TokenPair, error) {
	user, err := s.userRepo.FindByGoogleID(ctx, info.Sub)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &domain.User{
			Email:     info.Email,
			Name:      info.Name,
			AvatarURL: info.Picture,
			GoogleID:  info.Sub,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}
	return s.issueTokenPair(ctx, user, deviceInfo)

}

func (s *AuthService) issueTokenPair(ctx context.Context, user *domain.User, deviceInfo string) (*TokenPair, error) {
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	rawRefresh, hashedRefresh := s.jwtService.GenerateRefreshToken()

	refreshRecord := &domain.RefreshToken{
		UserID:     user.ID,
		TokenHash:  hashedRefresh,
		DeviceInfo: deviceInfo,
		ExpiresAt:  time.Now().Add(s.jwtService.RefreshTokenTTL()),
	}
	if err := s.refreshRepo.Create(ctx, refreshRecord); err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresAt:    expiresAt,
		User:         user,
	}, nil
}

// RefreshAccessToken cấp access token mới từ 1 refresh token còn hợp lệ.
func (s *AuthService) RefreshAccessToken(ctx context.Context, rawRefreshToken string) (*TokenPair, error) {
	hash := s.jwtService.HashToken(rawRefreshToken)

	record, err := s.refreshRepo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if record == nil || record.Revoked || record.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(ctx, record.UserID)
	if err != nil || user == nil {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresAt:    expiresAt,
		User:         user,
	}, nil
}

// Logout thu hồi 1 refresh token cụ thể (đăng xuất 1 thiết bị).
func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	hash := s.jwtService.HashToken(rawRefreshToken)
	record, err := s.refreshRepo.FindByHash(ctx, hash)
	if err != nil {
		return err
	}
	if record == nil {
		return nil
	}
	return s.refreshRepo.RevokeByID(ctx, record.ID)
}
