package collector

import (
	"log/slog"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/prometheus/client_golang/prometheus"
)

type podCollector struct {
	info            typedDesc
	state           typedDesc
	numOfContainers typedDesc
	created         typedDesc
	logger          *slog.Logger
}

type podDescLabels struct {
	labels      []string
	labelsValue []string
}

func init() {
	registerCollector("pod", defaultDisabled, NewPodStatsCollector)
}

// NewPodStatsCollector returns a Collector exposing pod stats information.
func NewPodStatsCollector(logger *slog.Logger) (Collector, error) {
	return &podCollector{
		info: typedDesc{
			nil, prometheus.GaugeValue,
		},
		state: typedDesc{
			nil, prometheus.GaugeValue,
		},
		numOfContainers: typedDesc{
			nil, prometheus.GaugeValue,
		},
		created: typedDesc{
			nil, prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes pod stats.
func (c *podCollector) Update(ch chan<- prometheus.Metric) error {
	defaultPodLabels := []string{"id"}

	reports, err := pdcs.Pods()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		podLabelsInfo := c.getPodDescLabels(rep)

		if enhanceAllMetrics {
			defaultPodLabels = podLabelsInfo.labels
		}

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pod", "info"),
			"Pod information",
			podLabelsInfo.labels, nil,
		)

		stateDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pod", "state"),
			"Pods current state current state (-1=unknown,0=created,1=error,2=exited,3=paused,4=running,5=degraded,6=stopped).",
			defaultPodLabels, nil)

		numOfCntDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pod", "containers"),
			"Number of containers in a pod.",
			defaultPodLabels, nil,
		)

		createdDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pod", "created_seconds"),
			"Pods creation time in unixtime.",
			defaultPodLabels, nil,
		)

		c.info.desc = infoDesc
		c.state.desc = stateDesc
		c.numOfContainers.desc = numOfCntDesc
		c.created.desc = createdDesc

		ch <- c.info.mustNewConstMetric(1, podLabelsInfo.labelsValue...)

		if enhanceAllMetrics {
			ch <- c.state.mustNewConstMetric(float64(rep.State), podLabelsInfo.labelsValue...)
			ch <- c.numOfContainers.mustNewConstMetric(float64(rep.NumOfContainers), podLabelsInfo.labelsValue...)
			ch <- c.created.mustNewConstMetric(float64(rep.Created), podLabelsInfo.labelsValue...)

			continue
		}

		ch <- c.state.mustNewConstMetric(float64(rep.State), rep.ID)
		ch <- c.numOfContainers.mustNewConstMetric(float64(rep.NumOfContainers), rep.ID)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID)
	}

	return nil
}

func (c *podCollector) getPodDescLabels(rep pdcs.Pod) *podDescLabels {
	podLabels := []string{"id", "name", "infra_id"}
	podLabelsValue := []string{rep.ID, rep.Name, rep.InfraID}

	extraLabels, extraValues := c.getExtraLabelsAndValues(podLabels, rep)

	podLabels = append(podLabels, extraLabels...)
	podLabelsValue = append(podLabelsValue, extraValues...)

	pDescLabels := podDescLabels{
		labels:      podLabels,
		labelsValue: podLabelsValue,
	}

	return &pDescLabels
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
