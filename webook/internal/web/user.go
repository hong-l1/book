package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/pkg/wrapper"
	"github.com/hong-l1/project/webook/internal/service"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// UseHandle 定义有关用户的路由
type UserHandle struct {
	codesvc          service.CodeService
	svc              service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	phoneNumberExp   *regexp.Regexp
	ijwt.Handle
	cmd redis.Cmdable
	l   logger.Loggerv1
}
type LohInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const biz = "login"

func NewUserHandle(svc service.UserService, codesvc service.CodeService, ijwthandle ijwt.Handle, l logger.Loggerv1) *UserHandle {
	const (
		emailRegexPattern    = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		passwordRegexPattern = "^(?=.*[A-Za-z])(?=.*\\d).{8,}$"
		phoneNumberPattern   = `^1[3-9]\d{9}$`
	)
	return &UserHandle{
		svc:              svc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneNumberExp:   regexp.MustCompile(phoneNumberPattern, regexp.None),
		codesvc:          codesvc,
		Handle:           ijwthandle,
		l:                l,
	}
}
func (u *UserHandle) RegisterUsersRoutes(server *gin.Engine) {
	ug := server.Group("users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.LogIn)
	ug.POST("/login", wrapper.Wrapper[LohInReq](u.LogInJwt, u.l.With(logger.String("method", "LoginJWt"))))
	ug.POST("/edit", u.Edit)
	ug.POST("/logout", u.Logout)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/refresh_token", u.RefreshToken)
}
func (u *UserHandle) Logout(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登成功",
	})
}
func (u *UserHandle) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(400, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err := u.codesvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("校验验证码出错", zap.Error(err))
		return
	}
	user, err := u.svc.FindORCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = u.SetLogintoken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证成功",
	})
}
func (u *UserHandle) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.phoneNumberExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入合法的手机号",
		})
		return
	}
	err = u.codesvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}
func (u *UserHandle) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}
	ok, err = u.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于八位，包含数字，特数字符")
		return
	}
	//在这调service方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicated) {
		ctx.String(http.StatusOK, "冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Println(req)
}
func (u *UserHandle) LogIn(ctx *gin.Context) {
	type LohInReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LohInReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	user, err := u.svc.LogIn(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//登录成功
	//设置session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}
func (u *UserHandle) LogInJwt(ctx *gin.Context, req LohInReq) (Result, error) {
	user, err := u.svc.LogIn(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		return Result{Msg: "账号或密码不对", Code: 4}, fmt.Errorf("账号或密码不对 %w", err)
	}
	if err != nil {
		return Result{Msg: "系统错误", Code: 5}, fmt.Errorf("系统错误 %w", err)
	}
	//创建token结构体
	if err := u.SetLogintoken(ctx, user.Id); err != nil {
		return Result{Msg: "系统错误", Code: 5}, fmt.Errorf("系统错误 %w", err)
	}
	return Result{Msg: "登录成功", Code: 5}, nil
}
func (u *UserHandle) RefreshToken(ctx *gin.Context) {
	refreshtoken := u.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshtoken, &rc, func(*jwt.Token) (interface{}, error) {
		return ijwt.Refresh_token_key, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.SetJWTtoken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		zap.L().Error("设置JWT出现异常", zap.Error(err),
			zap.String("method", "RefreshToken"))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}
func (u *UserHandle) Edit(ctx *gin.Context) {
	type EditReq struct {
		Id           int64
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	if _, err := time.Parse("2006-01-02", req.Birthday); err != nil {
		ctx.String(http.StatusBadRequest, "生日格式应为 YYYY-MM-DD")
		return
	}
	c, ok := ctx.Get("claim")
	if !ok {
		ctx.String(http.StatusUnauthorized, "认证失败")
		return
	}
	claim, ok := c.(*ijwt.Claim)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//格式符合要求，根据邮箱来查找到用户，并且补全信息
	err := u.svc.Edit(ctx, domain.User{
		Id:           claim.UserId,
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "修改成功")
}
func (u *UserHandle) Profile(ctx *gin.Context) {

	ctx.String(http.StatusOK, "登录界面")
	return
}
func (u *UserHandle) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claim")
	if !ok {
		ctx.String(http.StatusUnauthorized, "认证失败")
		return
	}
	claim, ok := c.(*ijwt.Claim)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	t, err := u.svc.Profile(ctx, domain.User{
		Id: claim.UserId,
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, "查询用户信息失败")
		return
	}
	t.Id = claim.UserId
	ctx.JSON(http.StatusOK, t)
	return
}
