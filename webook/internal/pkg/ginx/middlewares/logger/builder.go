package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLogger)
}
type AccessLogger struct {
	DUration time.Duration
	Method   string
	Url      string
	Reqbody  string
	Resbody  string
	Status   int
}

func NewBuilder(fn func(ctx context.Context, al *AccessLogger)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}
func (b *MiddlewareBuilder) AllowReBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}
func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}
func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLogger{
			Method: c.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && c.Request.Body != nil {
			body, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.Reqbody = string(body)
		}
		if b.allowRespBody && c.Request.Body != nil {
			c.Writer = responseWriter{
				al:             al,
				ResponseWriter: c.Writer,
			}
		}
		defer func() {
			al.DUration = time.Since(start)
			b.loggerFunc(c, al)
		}()
		c.Next()

	}
}

type responseWriter struct {
	al *AccessLogger
	gin.ResponseWriter
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.Resbody = string(data)
	return w.ResponseWriter.Write(data)
}
func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w responseWriter) WriteString(data string) (int, error) {
	w.al.Resbody = data
	return w.ResponseWriter.WriteString(data)
}
