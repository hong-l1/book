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
		Name:      m.Name + "_request_time",
		Help:      m.Help,
		//ConstLabels: map[string]string{},
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)
	prometheus.MustRegister(summary)
	guage := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: m.NameSpace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_active_cnt",
		Help:      m.Help,
	})
	prometheus.MustRegister(guage)
	return func(ctx *gin.Context) {
		start := time.Now()
		guage.Inc()
		defer func() {
			duration := time.Since(start)
			pattern := ctx.FullPath()
			guage.Dec()
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
func NewMidddlewareBuilder(NameSpace string, Subsystem string, Name string, Help string) *MidddlewareBuilder {
	return &MidddlewareBuilder{
		NameSpace: NameSpace,
		Subsystem: Subsystem,
		Name:      Name,
		Help:      Help,
	}
}
