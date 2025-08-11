package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareJwtBuilder struct {
	ijwt.Handle
	cmd redis.Cmdable
}

func NewLoginMiddlewareJwtBuilder(ijwthandle ijwt.Handle) *LoginMiddlewareJwtBuilder {
	return &LoginMiddlewareJwtBuilder{
		Handle: ijwthandle,
	}
}
func (l *LoginMiddlewareJwtBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.Path, "/login") || strings.HasSuffix(ctx.Request.URL.Path, "/signup") || strings.HasSuffix(ctx.Request.URL.Path, "/code/send") || strings.HasSuffix(ctx.Request.URL.Path, "/login_sms") {
			return
		}
		tokenStr := l.ExtractToken(ctx)
		claim := &ijwt.Claim{}
		t, err := jwt.ParseWithClaims(tokenStr, claim, func(token *jwt.Token) (interface{}, error) {
			return []byte("uLhPJesQ6V2bJKEpLQzhgfozn0fZZXqL"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !t.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claim.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = l.CheckSession(ctx, claim.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claim", claim)
	}

}
