package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type networkCollector struct {
	info   typedDesc
	logger log.Logger
}

func init() {
	registerCollector("network", defaultDisabled, NewNetworkStatsCollector)
}

// NewNetworkStatsCollector returns a Collector exposing network stats information.
func NewNetworkStatsCollector(logger log.Logger) (Collector, error) {
	return &networkCollector{
		info: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "info"),
				"Network information.",
				[]string{"name", "id", "driver", "interface", "labels"}, nil,
			), prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes networks stats.
func (c *networkCollector) Update(ch chan<- prometheus.Metric) error {
	reports, err := pdcs.Networks()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		ch <- c.info.mustNewConstMetric(1, rep.Name, rep.ID, rep.Driver, rep.NetworkInterface, rep.Labels)
	}

	return nil
}
