package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hong-l1/project/webook/internal/service"
	"github.com/hong-l1/project/webook/internal/service/oauth2/wechat"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WeChatHandle struct {
	svc    wechat.WechatService
	uersvc service.UserService
	ijwt.Handle
	statekey []byte
}

func NewOAuth2WeChatHandle(svc wechat.WechatService, uersvc service.UserService, ijwthandle ijwt.Handle) *OAuth2WeChatHandle {
	return &OAuth2WeChatHandle{
		svc:      svc,
		uersvc:   uersvc,
		statekey: []byte("hPJesV2bzzJKEpLQzhgfozn0fZZXqL18"),
		Handle:   ijwthandle,
	}
}
func (h *OAuth2WeChatHandle) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.Auth2URL)
	g.Any("/callback", h.Callback)
}
func (h *OAuth2WeChatHandle) Auth2URL(c *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(c.Request.Context(), state)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "URL构造失败",
		})
		return
	}
	err = h.SetStateCookie(c, state)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: url,
	})
}

func (h *OAuth2WeChatHandle) SetStateCookie(c *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenstr, err := token.SignedString(h.statekey)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return err
	}
	c.SetCookie("jwt-state", tokenstr, 600,
		"/oauth2/wechat/callback", "",
		false, true)
	return nil
}
func (h *OAuth2WeChatHandle) Callback(c *gin.Context) {
	code := c.Query("code")
	err := h.Verifystate(c)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	info, err := h.svc.VerifyCode(c, code)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u, err := h.uersvc.FindORCreateBywechat(c.Request.Context(), info)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.SetLogintoken(c, u.Id)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
}

func (h *OAuth2WeChatHandle) Verifystate(c *gin.Context) error {
	state := c.Query("state")
	ck, err := c.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到state的cookie，%w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.statekey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token已经过期%w", err)
	}
	if state != sc.State {
		return fmt.Errorf("state不相等%w", err)
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
