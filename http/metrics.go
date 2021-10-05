package http

import (
	"strconv"
	"time"

	"github.com/pretcat/ugc_test_task/logger"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	httpConnTotalMetric       = "ugc_test_task_http_connections_total"
	httpRequestDurationMetric = "ugc_test_task_http_request_duration_seconds"

	httpStatusCodeKey = "http_status_code"
	httpPathKey       = "http_path"
	httpMethodKey     = "http_method"
)

type requestsDurationMetric struct {
	*prometheus.HistogramVec
}

func newRequestsDurationMetric() *requestsDurationMetric {
	c := new(requestsDurationMetric)
	opts := prometheus.HistogramOpts{
		Name: httpRequestDurationMetric,
		Help: "duration http requests",
	}
	c.HistogramVec = prometheus.NewHistogramVec(opts, []string{httpStatusCodeKey, httpPathKey, httpMethodKey})
	if err := prometheus.Register(c.HistogramVec); err != nil {
		logger.Msg(opts.Name).AddMsg("not registered in prometheus").Warning(err.Error())
	}
	return c
}

func (c *requestsDurationMetric) write(req Request, statusCode int) {
	c.WithLabelValues(
		strconv.FormatInt(int64(statusCode), 10),
		req.Path(),
		req.Method).Observe(time.Since(req.Time()).Seconds())
}

type connectionsCounterMetric struct {
	prometheus.Gauge
}

func newConnectionsCounterMetric() *connectionsCounterMetric {
	c := new(connectionsCounterMetric)
	opts := prometheus.GaugeOpts{
		Name: httpConnTotalMetric,
		Help: "current amount open http connections",
	}
	c.Gauge = prometheus.NewGauge(opts)
	if err := prometheus.Register(c.Gauge); err != nil {
		logger.Msg(opts.Name).AddMsg("not registered in prometheus").Warning(err.Error())
	}
	return c
}
