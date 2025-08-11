package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var Access_token_key = []byte("uLhPJesQ6V2bJKEpLQzhgfozn0fZZXqL")
var Refresh_token_key = []byte("uLhPJesz6V3bJKEpLQzhgfozn8fZZXqL")

type RedisJWT struct {
	cmd redis.Cmdable
}

func (h *RedisJWT) ClearToken(ctx *gin.Context) error {
	ctx.Header("jwt-token", "")
	ctx.Header("refesh-token", "")
	claim := ctx.MustGet("claims").(*Claim)
	return h.cmd.Exists(ctx, fmt.Sprintf("users:ssid%s", claim.Ssid)).Err()
}

func (h *RedisJWT) ExtractToken(ctx *gin.Context) string {
	token := ctx.GetHeader("authorization")
	seg := strings.Split(token, " ")
	if len(seg) != 2 {
		return ""
	}
	return seg[1]
}

func (h *RedisJWT) SetLogintoken(ctx *gin.Context, id int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTtoken(ctx, id, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefeshtoken(ctx, id, ssid)
	return err
}

func (h *RedisJWT) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if cnt == 0 {
			return nil
		}
		return fmt.Errorf("session expired")
	default:
		return err
	}
}

func NewRedisJWT(cmd redis.Cmdable) Handle {
	return &RedisJWT{
		cmd: cmd,
	}
}

func (h *RedisJWT) SetJWTtoken(ctx *gin.Context, id int64, ssid string) error {
	claim := Claim{
		Ssid:      ssid,
		UserId:    id,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(Access_token_key)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("jwt-token", tokenStr)
	return nil
}
func (h *RedisJWT) SetRefeshtoken(ctx *gin.Context, uid int64, ssid string) error {
	claim := RefreshClaims{
		Ssid: ssid,
		Uid:  uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(Refresh_token_key)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("refesh-token", tokenStr)
	return nil
}
