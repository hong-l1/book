package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.Path, "/login") || strings.HasSuffix(ctx.Request.URL.Path, "/signup") || strings.HasSuffix(ctx.Request.URL.Path, "/code/send") || strings.HasSuffix(ctx.Request.URL.Path, "/login_sms") {
			return
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime := sess.Get("updateTime")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		if updateTime == nil {
			sess.Set("updateTime", now)
			sess.Save()
			return
		}
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > 30*time.Second {
			sess.Set("updateTime", now)
			sess.Save()
			return
		}
	}
}
