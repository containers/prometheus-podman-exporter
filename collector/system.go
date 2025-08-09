package collector

import (
	"log/slog"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/prometheus/client_golang/prometheus"
)

type systemCollector struct {
	podmanVer  typedDesc
	runtimeVer typedDesc
	conmonVer  typedDesc
	buildahVer typedDesc
	logger     *slog.Logger
}

func init() {
	registerCollector("system", defaultDisabled, NewSystemCollector)
}

// NewSystemCollector returns a Collector exposing podman system information.
func NewSystemCollector(logger *slog.Logger) (Collector, error) {
	return &systemCollector{
		podmanVer: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "system", "api_version"),
				"Podman system api version.",
				[]string{"version"}, nil,
			), prometheus.GaugeValue,
		},
		runtimeVer: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "system", "runtime_version"),
				"Podman system runtime version.",
				[]string{"version"}, nil,
			), prometheus.GaugeValue,
		},
		conmonVer: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "system", "conmon_version"),
				"Podman system conmon version.",
				[]string{"version"}, nil,
			), prometheus.GaugeValue,
		},
		buildahVer: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "system", "buildah_version"),
				"Podman system buildahVer version.",
				[]string{"version"}, nil,
			), prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes system information.
func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	info, err := pdcs.SystemInfo()
	if err != nil {
		return err
	}

	ch <- c.podmanVer.mustNewConstMetric(1, info.Podman)

	ch <- c.runtimeVer.mustNewConstMetric(1, info.Runtime)

	ch <- c.conmonVer.mustNewConstMetric(1, info.Conmon)

	ch <- c.buildahVer.mustNewConstMetric(1, info.Buildah)

	return nil
}
