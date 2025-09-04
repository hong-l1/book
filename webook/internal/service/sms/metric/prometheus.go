package metric

import (
	"context"
	"github.com/hong-l1/project/webook/internal/service/sms"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusDecorator struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewPrometheusDecorator(svc sms.Service) *PrometheusDecorator {
	v := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "gobook",
		Subsystem: "sms",
		Name:      "sms_requset_time",
		Help:      "监控短信服务",
	}, []string{"biz"})
	prometheus.MustRegister(v)
	return &PrometheusDecorator{
		svc:    svc,
		vector: v,
	}
}
func (p *PrometheusDecorator) SendSMS(ctx context.Context, biz string, args []string, numbers ...string) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		p.vector.WithLabelValues(biz).Observe(float64(duration))
	}()
	return p.SendSMS(ctx, biz, args, numbers...)
}
