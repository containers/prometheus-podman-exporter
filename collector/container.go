package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type containerCollector struct {
	info             typedDesc
	state            typedDesc
	health           typedDesc
	created          typedDesc
	started          typedDesc
	exited           typedDesc
	exitCode         typedDesc
	pids             typedDesc
	cpu              typedDesc
	cpuSystem        typedDesc
	memUsage         typedDesc
	memLimit         typedDesc
	netInput         typedDesc
	netInputDropped  typedDesc
	netInputErrors   typedDesc
	netInputPackets  typedDesc
	netOutput        typedDesc
	netOutputDropped typedDesc
	netOutputErrors  typedDesc
	netOutputPackets typedDesc
	blockInput       typedDesc
	blockOutput      typedDesc
	rwSize           typedDesc
	rootFsSize       typedDesc
	logger           log.Logger
}

type containerDescLabels struct {
	labels      []string
	labelsValue []string
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
			nil, prometheus.GaugeValue,
		},
		health: typedDesc{
			nil, prometheus.GaugeValue,
		},
		created: typedDesc{
			nil, prometheus.GaugeValue,
		},
		started: typedDesc{
			nil, prometheus.GaugeValue,
		},
		exited: typedDesc{
			nil, prometheus.GaugeValue,
		},
		exitCode: typedDesc{
			nil, prometheus.GaugeValue,
		},
		pids: typedDesc{
			nil, prometheus.GaugeValue,
		},
		cpu: typedDesc{
			nil, prometheus.CounterValue,
		},
		cpuSystem: typedDesc{
			nil, prometheus.CounterValue,
		},
		memUsage: typedDesc{
			nil, prometheus.GaugeValue,
		},
		memLimit: typedDesc{
			nil, prometheus.GaugeValue,
		},
		netInput: typedDesc{
			nil, prometheus.CounterValue,
		},
		netInputDropped: typedDesc{
			nil, prometheus.CounterValue,
		},
		netInputErrors: typedDesc{
			nil, prometheus.CounterValue,
		},
		netInputPackets: typedDesc{
			nil, prometheus.CounterValue,
		},
		netOutput: typedDesc{
			nil, prometheus.CounterValue,
		},
		netOutputDropped: typedDesc{
			nil, prometheus.CounterValue,
		},
		netOutputErrors: typedDesc{
			nil, prometheus.CounterValue,
		},
		netOutputPackets: typedDesc{
			nil, prometheus.CounterValue,
		},
		blockInput: typedDesc{
			nil, prometheus.CounterValue,
		},
		blockOutput: typedDesc{
			nil, prometheus.CounterValue,
		},
		rwSize: typedDesc{
			nil, prometheus.GaugeValue,
		},
		rootFsSize: typedDesc{
			nil, prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes container stats.
func (c *containerCollector) Update(ch chan<- prometheus.Metric) error {
	defaultContainersLabel := []string{"id", "pod_id", "pod_name"}

	reports, err := pdcs.Containers()
	if err != nil {
		return err
	}

	statReports, err := pdcs.ContainersStats()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		cntLabelsInfo := c.getContainerDescLabel(rep)

		if enhanceAllMetrics {
			defaultContainersLabel = cntLabelsInfo.labels
		}

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "container", "info"),
			"Container information.",
			cntLabelsInfo.labels, nil,
		)

		c.info.desc = infoDesc

		cntStat := getContainerStat(rep.ID, statReports)

		ch <- c.info.mustNewConstMetric(1, cntLabelsInfo.labelsValue...)
		c.updateInfo(ch, &rep, cntLabelsInfo, defaultContainersLabel, enhanceAllMetrics)
		c.updateStats(ch, &rep, cntStat, cntLabelsInfo, defaultContainersLabel, enhanceAllMetrics)
	}

	return nil
}

func (c *containerCollector) updateInfo(
	ch chan<- prometheus.Metric,
	rep *pdcs.Container,
	cntLabelsInfo *containerDescLabels,
	defaultLabels []string,
	enhance bool,
) {
	stateDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "state"),
		//nolint:lll
		"Container current state (-1=unknown,0=created,1=initialized,2=running,3=stopped,4=paused,5=exited,6=removing,7=stopping).",
		defaultLabels, nil,
	)

	healthDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "health"),
		"Container current health (-1=unknown,0=healthy,1=unhealthy,2=starting).",
		defaultLabels, nil,
	)

	createdDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "created_seconds"),
		"Container creation time in unixtime.",
		defaultLabels, nil,
	)

	startedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "started_seconds"),
		"Container started time in unixtime.",
		defaultLabels, nil,
	)

	exitedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "exited_seconds"),
		"Container exited time in unixtime.",
		defaultLabels, nil,
	)

	exitedCodeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "exit_code"),
		"Container exit code, if the container has not exited or restarted then the exit code will be 0.",
		defaultLabels, nil,
	)

	rwSizeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "rw_size_bytes"),
		"Container top read-write layer size in bytes.",
		defaultLabels, nil,
	)

	rootFsSizeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "rootfs_size_bytes"),
		"Container root filesystem size in bytes.",
		defaultLabels, nil,
	)

	c.state.desc = stateDesc
	c.health.desc = healthDesc
	c.created.desc = createdDesc
	c.started.desc = startedDesc
	c.exited.desc = exitedDesc
	c.exitCode.desc = exitedCodeDesc
	c.rwSize.desc = rwSizeDesc
	c.rootFsSize.desc = rootFsSizeDesc

	if enhance {
		ch <- c.state.mustNewConstMetric(float64(rep.State), cntLabelsInfo.labelsValue...)
		ch <- c.health.mustNewConstMetric(float64(rep.Health), cntLabelsInfo.labelsValue...)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), cntLabelsInfo.labelsValue...)
		ch <- c.started.mustNewConstMetric(float64(rep.Started), cntLabelsInfo.labelsValue...)
		ch <- c.exited.mustNewConstMetric(float64(rep.Exited), cntLabelsInfo.labelsValue...)
		ch <- c.exitCode.mustNewConstMetric(float64(rep.ExitCode), cntLabelsInfo.labelsValue...)
		ch <- c.rwSize.mustNewConstMetric(float64(rep.RwSize), cntLabelsInfo.labelsValue...)
		ch <- c.rootFsSize.mustNewConstMetric(float64(rep.RootFsSize), cntLabelsInfo.labelsValue...)

		return
	}

	ch <- c.state.mustNewConstMetric(float64(rep.State), rep.ID, rep.PodID, rep.PodName)
	ch <- c.health.mustNewConstMetric(float64(rep.Health), rep.ID, rep.PodID, rep.PodName)
	ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID, rep.PodID, rep.PodName)
	ch <- c.started.mustNewConstMetric(float64(rep.Started), rep.ID, rep.PodID, rep.PodName)
	ch <- c.exited.mustNewConstMetric(float64(rep.Exited), rep.ID, rep.PodID, rep.PodName)
	ch <- c.exitCode.mustNewConstMetric(float64(rep.ExitCode), rep.ID, rep.PodID, rep.PodName)
	ch <- c.rwSize.mustNewConstMetric(float64(rep.RwSize), rep.ID, rep.PodID, rep.PodName)
	ch <- c.rootFsSize.mustNewConstMetric(float64(rep.RootFsSize), rep.ID, rep.PodID, rep.PodName)
}

func (c *containerCollector) updateStats(
	ch chan<- prometheus.Metric,
	rep *pdcs.Container,
	cntStat *pdcs.ContainerStat,
	cntLabelsInfo *containerDescLabels,
	defaultLabels []string,
	enhance bool,
) {
	pidsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "pids"),
		"Container pid number.",
		defaultLabels, nil,
	)

	cpuDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "cpu_seconds_total"),
		"total CPU time spent for container in seconds.",
		defaultLabels, nil,
	)

	cpuSystemDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "cpu_system_seconds_total"),
		"total system CPU time spent for container in seconds.",
		defaultLabels, nil,
	)

	memUsageDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "mem_usage_bytes"),
		"Container memory usage.",
		defaultLabels, nil,
	)

	memLimitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "mem_limit_bytes"),
		"Container memory limit.",
		defaultLabels, nil,
	)

	netInputDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_input_total"),
		"Container network input in bytes.",
		defaultLabels, nil,
	)

	netInputDroppedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_input_dropped_total"),
		"Container network input dropped.",
		defaultLabels, nil,
	)

	netInputErrorsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_input_errors_total"),
		"Container network input errors.",
		defaultLabels, nil,
	)

	netInputPacketsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_input_packets_total"),
		"Container network input packets.",
		defaultLabels, nil,
	)

	netOutputDroppedDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_output_dropped_total"),
		"Container network output dropped.",
		defaultLabels, nil,
	)

	netOutputErrorsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_output_errors_total"),
		"Container network output errors.",
		defaultLabels, nil,
	)

	netOutputPacketsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_output_packets_total"),
		"Container network output packets.",
		defaultLabels, nil,
	)

	netOutputDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "net_output_total"),
		"Container network output in bytes.",
		defaultLabels, nil,
	)

	blockInputDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "block_input_total"),
		"Container block input in bytes.",
		defaultLabels, nil,
	)

	blockOutputDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "container", "block_output_total"),
		"Container block output in bytes.",
		defaultLabels, nil,
	)

	c.pids.desc = pidsDesc
	c.cpu.desc = cpuDesc
	c.cpuSystem.desc = cpuSystemDesc
	c.memUsage.desc = memUsageDesc
	c.memLimit.desc = memLimitDesc
	c.netInput.desc = netInputDesc
	c.netInputDropped.desc = netInputDroppedDesc
	c.netInputErrors.desc = netInputErrorsDesc
	c.netInputPackets.desc = netInputPacketsDesc
	c.netOutput.desc = netOutputDesc
	c.netOutputDropped.desc = netOutputDroppedDesc
	c.netOutputErrors.desc = netOutputErrorsDesc
	c.netOutputPackets.desc = netOutputPacketsDesc
	c.blockInput.desc = blockInputDesc
	c.blockOutput.desc = blockOutputDesc

	if enhance {
		if cntStat != nil {
			ch <- c.pids.mustNewConstMetric(float64(cntStat.PIDs), cntLabelsInfo.labelsValue...)
			ch <- c.cpu.mustNewConstMetric(cntStat.CPU, cntLabelsInfo.labelsValue...)
			ch <- c.cpuSystem.mustNewConstMetric(cntStat.CPUSystem, cntLabelsInfo.labelsValue...)
			ch <- c.memUsage.mustNewConstMetric(float64(cntStat.MemUsage), cntLabelsInfo.labelsValue...)
			ch <- c.memLimit.mustNewConstMetric(float64(cntStat.MemLimit), cntLabelsInfo.labelsValue...)
			ch <- c.netInput.mustNewConstMetric(float64(cntStat.NetInput), cntLabelsInfo.labelsValue...)
			ch <- c.netInputDropped.mustNewConstMetric(float64(cntStat.NetInputDropped), cntLabelsInfo.labelsValue...)
			ch <- c.netInputErrors.mustNewConstMetric(float64(cntStat.NetInputErrors), cntLabelsInfo.labelsValue...)
			ch <- c.netInputPackets.mustNewConstMetric(float64(cntStat.NetInputPackets), cntLabelsInfo.labelsValue...)
			ch <- c.netOutput.mustNewConstMetric(float64(cntStat.NetOutput), cntLabelsInfo.labelsValue...)
			ch <- c.netOutputDropped.mustNewConstMetric(float64(cntStat.NetOutputDropped), cntLabelsInfo.labelsValue...)
			ch <- c.netOutputErrors.mustNewConstMetric(float64(cntStat.NetOutputErrors), cntLabelsInfo.labelsValue...)
			ch <- c.netOutputPackets.mustNewConstMetric(float64(cntStat.NetOutputPackets), cntLabelsInfo.labelsValue...)
			ch <- c.blockInput.mustNewConstMetric(float64(cntStat.BlockInput), cntLabelsInfo.labelsValue...)
			ch <- c.blockOutput.mustNewConstMetric(float64(cntStat.BlockOutput), cntLabelsInfo.labelsValue...)
		}

		return
	}

	if cntStat != nil {
		ch <- c.pids.mustNewConstMetric(float64(cntStat.PIDs), rep.ID, rep.PodID, rep.PodName)
		ch <- c.cpu.mustNewConstMetric(cntStat.CPU, rep.ID, rep.PodID, rep.PodName)
		ch <- c.cpuSystem.mustNewConstMetric(cntStat.CPUSystem, rep.ID, rep.PodID, rep.PodName)
		ch <- c.memUsage.mustNewConstMetric(float64(cntStat.MemUsage), rep.ID, rep.PodID, rep.PodName)
		ch <- c.memLimit.mustNewConstMetric(float64(cntStat.MemLimit), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netInput.mustNewConstMetric(float64(cntStat.NetInput), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netInputDropped.mustNewConstMetric(float64(cntStat.NetInputDropped), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netInputErrors.mustNewConstMetric(float64(cntStat.NetInputErrors), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netInputPackets.mustNewConstMetric(float64(cntStat.NetInputPackets), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netOutput.mustNewConstMetric(float64(cntStat.NetOutput), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netOutputDropped.mustNewConstMetric(float64(cntStat.NetOutputDropped), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netOutputErrors.mustNewConstMetric(float64(cntStat.NetOutputErrors), rep.ID, rep.PodID, rep.PodName)
		ch <- c.netOutputPackets.mustNewConstMetric(float64(cntStat.NetOutputPackets), rep.ID, rep.PodID, rep.PodName)
		ch <- c.blockInput.mustNewConstMetric(float64(cntStat.BlockInput), rep.ID, rep.PodID, rep.PodName)
		ch <- c.blockOutput.mustNewConstMetric(float64(cntStat.BlockOutput), rep.ID, rep.PodID, rep.PodName)
	}
}

func (c *containerCollector) getContainerDescLabel(rep pdcs.Container) *containerDescLabels {
	containerLabels := []string{"id", "name", "image", "ports", "pod_id", "pod_name"}
	containerLabelsValue := []string{rep.ID, rep.Name, rep.Image, rep.Ports, rep.PodID, rep.PodName}

	extraLabels, extraValues := c.getExtraLabelsAndValues(containerLabels, rep)

	containerLabels = append(containerLabels, extraLabels...)
	containerLabelsValue = append(containerLabelsValue, extraValues...)

	cntDescLabels := containerDescLabels{
		labels:      containerLabels,
		labelsValue: containerLabelsValue,
	}

	return &cntDescLabels
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
