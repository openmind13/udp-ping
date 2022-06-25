package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "info",
	}, []string{"version"})
)
