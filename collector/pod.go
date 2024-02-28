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
			nil, prometheus.GaugeValue,
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
		infoMetric, infoValues := c.getPodInfoDesc(rep)
		c.info.desc = infoMetric

		ch <- c.info.mustNewConstMetric(1, infoValues...)
		ch <- c.state.mustNewConstMetric(float64(rep.State), rep.ID)
		ch <- c.numOfContainers.mustNewConstMetric(float64(rep.NumOfContainers), rep.ID)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID)
	}

	return nil
}

func (c *podCollector) getPodInfoDesc(rep pdcs.Pod) (*prometheus.Desc, []string) {
	podLabels := []string{"id", "name", "infra_id"}
	podLabelsValue := []string{rep.ID, rep.Name, rep.InfraID}

	extraLabels, extraValues := c.getExtraLabelsAndValues(podLabels, rep)

	podLabels = append(podLabels, extraLabels...)
	podLabelsValue = append(podLabelsValue, extraValues...)

	infoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "pod", "info"),
		"Pod information",
		podLabels, nil,
	)

	return infoDesc, podLabelsValue
}

func (c *podCollector) getExtraLabelsAndValues(collectorLabels []string, rep pdcs.Pod) ([]string, []string) {
	extraLabels := make([]string, 0)
	extraValues := make([]string, 0)

	for label, value := range rep.Labels {
		if slicesContains(collectorLabels, label) {
			continue
		}

		validLabel := sanitizeLabelName(label)
		if storeLabels {
			extraLabels = append(extraLabels, validLabel)
			extraValues = append(extraValues, value)
		} else if whitelistContains(label) {
			extraLabels = append(extraLabels, validLabel)
			extraValues = append(extraValues, value)
		}
	}

	return extraLabels, extraValues
}
