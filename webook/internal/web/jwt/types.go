package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handle interface {
	ClearToken(*gin.Context) error
	ExtractToken(ctx *gin.Context) string
	SetJWTtoken(ctx *gin.Context, id int64, ssid string) error
	SetLogintoken(ctx *gin.Context, id int64) error
	CheckSession(ctx *gin.Context, ssid string) error
}
type RefreshClaims struct {
	Ssid string
	Uid  int64
	jwt.RegisteredClaims
}
type Claim struct {
	Ssid string
	jwt.RegisteredClaims
	UserId    int64 `json:"user_id"`
	UserAgent string
}
