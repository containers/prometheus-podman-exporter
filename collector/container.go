package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type containerCollector struct {
	info        typedDesc
	state       typedDesc
	created     typedDesc
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
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "info"),
				"Container information.",
				[]string{"id", "name", "image", "ports", "pod_id"}, nil,
			), prometheus.GaugeValue,
		},
		state: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "state"),
				// nolint:lll
				"Container current state (-1=unknown,0=created,1=initialized,2=running,3=stopped,4=paused,5=exited,6=removing,7=stopping).",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		created: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "created_seconds"),
				"Container creation time in unixtime.",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		pids: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "pids"),
				"Container pid number.",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		cpu: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "cpu_seconds_total"),
				"total CPU time spent for container in seconds.",
				[]string{"id"}, nil,
			), prometheus.CounterValue,
		},
		cpuSystem: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "cpu_system_seconds_total"),
				"total system CPU time spent for container in seconds.",
				[]string{"id"}, nil,
			), prometheus.CounterValue,
		},
		memUsage: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "mem_usage_bytes"),
				"Container memory usage.",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		memLimit: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "mem_limit_bytes"),
				"Container memory limit.",
				[]string{"id"}, nil,
			), prometheus.GaugeValue,
		},
		netInput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "net_input_total"),
				"Container network input.",
				[]string{"id"}, nil,
			), prometheus.CounterValue,
		},
		netOutput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "net_output_total"),
				"Container network output.",
				[]string{"id"}, nil,
			), prometheus.CounterValue,
		},
		blockInput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "block_input_total"),
				"Container block input.",
				[]string{"id"}, nil,
			), prometheus.CounterValue,
		},
		blockOutput: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "container", "block_output_total"),
				"Container block output.",
				[]string{"id"}, nil,
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

	for _, rep := range reports {
		ch <- c.info.mustNewConstMetric(1, rep.ID, rep.Name, rep.Image, rep.Ports, rep.PodID)
		ch <- c.state.mustNewConstMetric(float64(rep.State), rep.ID)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID)
	}

	statReports, err := pdcs.ContainersStats()
	if err != nil {
		return err
	}

	for _, rep := range statReports {
		ch <- c.pids.mustNewConstMetric(float64(rep.PIDs), rep.ID)
		ch <- c.cpu.mustNewConstMetric(rep.CPU, rep.ID)
		ch <- c.cpuSystem.mustNewConstMetric(rep.CPUSystem, rep.ID)
		ch <- c.memUsage.mustNewConstMetric(float64(rep.MemUsage), rep.ID)
		ch <- c.memLimit.mustNewConstMetric(float64(rep.MemLimit), rep.ID)
		ch <- c.netInput.mustNewConstMetric(float64(rep.NetInput), rep.ID)
		ch <- c.netOutput.mustNewConstMetric(float64(rep.NetOutput), rep.ID)
		ch <- c.blockInput.mustNewConstMetric(float64(rep.BlockInput), rep.ID)
		ch <- c.blockOutput.mustNewConstMetric(float64(rep.BlockOutput), rep.ID)
	}

	return nil
}
