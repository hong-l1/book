package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MidddlewareBuilder struct {
	NameSpace string
	Subsystem string
	Name      string
	Help      string
}

func (m *MidddlewareBuilder) Build() gin.HandlerFunc {
	labels := []string{"method", "path", "code"}
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.NameSpace,
		Subsystem: m.Subsystem,
		Name:      m.Name,
		Help:      m.Help,
		ConstLabels: map[string]string{},
	}, labels)
	prometheus.MustRegister(summary)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			pattern := ctx.FullPath()
			if len(pattern) == 0 {
				pattern = "unknown"
			}
			summary.WithLabelValues(ctx.Request.Method,
				pattern,
				strconv.Itoa(ctx.Writer.Status()),
			).Observe(float64(duration.Milliseconds()))
		}()
		ctx.Next()
	}
}
