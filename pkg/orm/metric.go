package orm

import (
	"github.com/zeromicro/go-zero/core/metric"
)

const gormNamespace = "gorm_client"

var (
	//Histogram 统计QPS和相关的耗时
	metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: gormNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "gorm client requests duration(ms).",
		Labels:    []string{"table", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	//统计错误数量
	metricClientReqErrTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: gormNamespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "gorm client requests error count.",
		Labels:    []string{"table", "method", "is_error"},
	})
)
