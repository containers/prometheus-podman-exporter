package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type podCollector struct {
	info            typedDesc
	state           typedDesc
	numOfContainers typedDesc
	created         typedDesc
	logger          log.Logger
}

func init() {
	registerCollector("pod", defaultDisabled, NewPodStatsCollector)
}

// NewPodStatsCollector returns a Collector exposing pod stats information.
func NewPodStatsCollector(logger log.Logger) (Collector, error) {
	return &podCollector{
		info: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "pod", "info"),
				"Pod information",
				[]string{"id", "name", "infra_id"}, nil,
			), prometheus.GaugeValue,
		},
		state: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "pod", "state"),
				"Pods current state current state (-1=unknown,0=created,1=error,2=exited,3=paused,4=running,5=degraded,6=stopped).",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		numOfContainers: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "pod", "containers"),
				"Number of containers in a pod.",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		created: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "pod", "created_seconds"),
				"Pods creation time in unixtime.",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes pod stats.
func (c *podCollector) Update(ch chan<- prometheus.Metric) error {
	reports, err := pdcs.Pods()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		ch <- c.info.mustNewConstMetric(1, rep.ID, rep.Name, rep.InfraID)
		ch <- c.state.mustNewConstMetric(float64(rep.State), rep.ID)
		ch <- c.numOfContainers.mustNewConstMetric(float64(rep.NumOfContainers), rep.ID)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID)
	}

	return nil
}
