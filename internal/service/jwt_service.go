package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidAccessToken là lỗi domain của service, dùng chung mọi nơi access token
// bị từ chối (sai chữ ký, hết hạn, sai format).
var ErrInvalidAccessToken = apperrors.NewAuthorization("Access token không hợp lệ hoặc đã hết hạn")

// Claims tùy chỉnh gắn kèm trong access token.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type JWTService struct {
	cfg *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{cfg: cfg}
}

// GenerateAccessToken tạo JWT ngắn hạn (mặc định 15 phút), client gửi kèm mỗi request
// qua header Authorization: Bearer <token>. Cách này giống nhau cho web (Next.js) và mobile.
func (s *JWTService) GenerateAccessToken(userID uuid.UUID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWTAccessTTLMinutes) * time.Minute)

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTAccessSecret))
	return signed, expiresAt, err
}

// ParseAccessToken verify và giải mã access token.
func (s *JWTService) ParseAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}
		return []byte(s.cfg.JWTAccessSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidAccessToken
	}
	return claims, nil
}

// GenerateRefreshToken tạo 1 chuỗi ngẫu nhiên (UUID), KHÔNG phải JWT.
// Lý do: refresh token cần thu hồi được (revoke) trong DB, JWT thuần thì không revoke được
// trừ khi có thêm blacklist - dùng random string + lưu hash trong DB đơn giản và an toàn hơn.
func (s *JWTService) GenerateRefreshToken() (raw string, hash string) {
	raw = uuid.NewString() + "-" + uuid.NewString()
	hash = s.HashToken(raw)
	return raw, hash
}

func (s *JWTService) HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (s *JWTService) RefreshTokenTTL() time.Duration {
	return time.Duration(s.cfg.JWTRefreshTTLDays) * 24 * time.Hour
}
