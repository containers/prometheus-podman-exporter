package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type containerCollector struct {
	info        typedDesc
	state       typedDesc
	health      typedDesc
	created     typedDesc
	started     typedDesc
	exited      typedDesc
	exitCode    typedDesc
	pids        typedDesc
	cpu         typedDesc
	cpuSystem   typedDesc
	memUsage    typedDesc
	memLimit    typedDesc
	netInput    typedDesc
	netOutput   typedDesc
	blockInput  typedDesc
	blockOutput typedDesc
	logger      log.Logger
}

func init() {
	registerCollector("container", defaultEnabled, NewContainerStatsCollector)
}

// NewContainerStatsCollector returns a Collector exposing containers stats information.
func NewContainerStatsCollector(logger log.Logger) (Collector, error) {
	return &containerCollector{
		info: typedDesc{
			nil, prometheus.GaugeValue,
		},
		state: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "state"),
				//nolint:lll
				"Container current state (-1=unknown,0=created,1=initialized,2=running,3=stopped,4=paused,5=exited,6=removing,7=stopping).",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		health: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "health"),
				"Container current health (-1=unknown,0=healthy,1=unhealthy,2=starting).",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		created: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "created_seconds"),
				"Container creation time in unixtime.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		started: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "started_seconds"),
				"Container started time in unixtime.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		exited: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "exited_seconds"),
				"Container exited time in unixtime.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		exitCode: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "exit_code"),
				"Container exit code, if the container has not exited or restarted then the exit code will be 0.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		pids: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "pids"),
				"Container pid number.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		cpu: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "cpu_seconds_total"),
				"total CPU time spent for container in seconds.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.CounterValue,
		},
		cpuSystem: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "cpu_system_seconds_total"),
				"total system CPU time spent for container in seconds.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.CounterValue,
		},
		memUsage: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "mem_usage_bytes"),
				"Container memory usage.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		memLimit: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "mem_limit_bytes"),
				"Container memory limit.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.GaugeValue,
		},
		netInput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "net_input_total"),
				"Container network input in bytes.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.CounterValue,
		},
		netOutput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "net_output_total"),
				"Container network output in bytes.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.CounterValue,
		},
		blockInput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "block_input_total"),
				"Container block input in bytes.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.CounterValue,
		},
		blockOutput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "block_output_total"),
				"Container block output in bytes.",
				[]string{"id", "pod_id", "pod_name"}, nil,
			), prometheus.CounterValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes container stats.
func (c *containerCollector) Update(ch chan<- prometheus.Metric) error {
	reports, err := pdcs.Containers()
	if err != nil {
		return err
	}

	statReports, err := pdcs.ContainersStats()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		infoMetric, infoValues := c.getContainerInfoDesc(rep)
		c.info.desc = infoMetric
		cntStat := getContainerStat(rep.ID, statReports)

		ch <- c.info.mustNewConstMetric(1, infoValues...)
		ch <- c.state.mustNewConstMetric(float64(rep.State), rep.ID, rep.PodID, rep.PodName)
		ch <- c.health.mustNewConstMetric(float64(rep.Health), rep.ID, rep.PodID, rep.PodName)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID, rep.PodID, rep.PodName)
		ch <- c.started.mustNewConstMetric(float64(rep.Started), rep.ID, rep.PodID, rep.PodName)
		ch <- c.exited.mustNewConstMetric(float64(rep.Exited), rep.ID, rep.PodID, rep.PodName)
		ch <- c.exitCode.mustNewConstMetric(float64(rep.ExitCode), rep.ID, rep.PodID, rep.PodName)

		if cntStat != nil {
			ch <- c.pids.mustNewConstMetric(float64(cntStat.PIDs), rep.ID, rep.PodID, rep.PodName)
			ch <- c.cpu.mustNewConstMetric(cntStat.CPU, rep.ID, rep.PodID, rep.PodName)
			ch <- c.cpuSystem.mustNewConstMetric(cntStat.CPUSystem, rep.ID, rep.PodID, rep.PodName)
			ch <- c.memUsage.mustNewConstMetric(float64(cntStat.MemUsage), rep.ID, rep.PodID, rep.PodName)
			ch <- c.memLimit.mustNewConstMetric(float64(cntStat.MemLimit), rep.ID, rep.PodID, rep.PodName)
			ch <- c.netInput.mustNewConstMetric(float64(cntStat.NetInput), rep.ID, rep.PodID, rep.PodName)
			ch <- c.netOutput.mustNewConstMetric(float64(cntStat.NetOutput), rep.ID, rep.PodID, rep.PodName)
			ch <- c.blockInput.mustNewConstMetric(float64(cntStat.BlockInput), rep.ID, rep.PodID, rep.PodName)
			ch <- c.blockOutput.mustNewConstMetric(float64(cntStat.BlockOutput), rep.ID, rep.PodID, rep.PodName)
		}
	}

	return nil
}

func (c *containerCollector) getContainerInfoDesc(rep pdcs.Container) (*prometheus.Desc, []string) {
	containerLabels := []string{"id", "name", "image", "ports", "pod_id", "pod_name"}
	containerLabelsValue := []string{rep.ID, rep.Name, rep.Image, rep.Ports, rep.PodID, rep.PodName}

	extraLabels, extraValues := c.getExtraLabelsAndValues(containerLabels, rep)

	containerLabels = append(containerLabels, extraLabels...)
	containerLabelsValue = append(containerLabelsValue, extraValues...)

	infoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "info"),
		"Container information.",
		containerLabels, nil,
	)

	return infoDesc, containerLabelsValue
}

func (c *containerCollector) getExtraLabelsAndValues(
	collectorLabels []string,
	rep pdcs.Container,
) ([]string, []string) {
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

func getContainerStat(containerID string, statReport []pdcs.ContainerStat) *pdcs.ContainerStat {
	for _, cstat := range statReport {
		if cstat.ID == containerID {
			return &cstat
		}
	}

	return nil
}
