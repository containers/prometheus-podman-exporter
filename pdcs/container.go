package pdcs

import (
	"context"
	"time"

	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

const (
	nano float64 = 1e+9
)

// Container implements container's basic information and its state.
type Container struct {
	ID       string
	PodID    string // if container is part of pod
	Name     string
	Labels   map[string]string
	Image    string
	Created  int64
	Started  int64
	Exited   int64
	ExitCode int32
	Ports    string
	State    int
	Health   int
}

// ContainerStat implements container's stat.
type ContainerStat struct {
	ID          string
	Name        string
	PIDs        uint64
	CPU         float64
	CPUSystem   float64
	MemUsage    uint64
	MemLimit    uint64
	NetInput    uint64
	NetOutput   uint64
	BlockInput  uint64
	BlockOutput uint64
}

// Containers returns list of containers (Container).
func Containers() ([]Container, error) {
	containers := make([]Container, 0)

	reports, err := registry.ContainerEngine().ContainerList(registry.Context(), entities.ContainerListOptions{All: true})
	if err != nil {
		return containers, err
	}

	for _, rep := range reports {
		containers = append(containers, Container{
			ID:       getID(rep.ID),
			PodID:    getID(rep.Pod),
			Name:     rep.Names[0],
			Image:    rep.Image,
			Created:  rep.Created.Unix(),
			Started:  rep.StartedAt,
			Exited:   rep.ExitedAt,
			ExitCode: rep.ExitCode,
			State:    conReporter{rep}.state(),
			Health:   conReporter{rep}.health(),
			Ports:    conReporter{rep}.ports(),
			Labels:   rep.Labels,
		})
	}

	return containers, nil
}

// ContainersStats returns list of containers stats (ContainerStat).
func ContainersStats() ([]ContainerStat, error) {
	stat := make([]ContainerStat, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	reports, err := registry.ContainerEngine().ContainerStats(
		registry.Context(),
		[]string{},
		entities.ContainerStatsOptions{Stream: false, Interval: 1})
	if err != nil {
		return stat, err
	}

	getStat := func() ([]define.ContainerStats, error) {
		for {
			select {
			case <-ctx.Done():
				return nil, ErrDeadline
			case s := <-reports:
				return s.Stats, nil
			}
		}
	}

	statReport, err := getStat()
	if err != nil {
		return nil, err
	}

	for _, rep := range statReport {
		stat = append(stat, ContainerStat{
			ID:          getID(rep.ContainerID),
			Name:        rep.Name,
			PIDs:        rep.PIDs,
			CPU:         float64(rep.CPUNano) / nano,
			CPUSystem:   float64(rep.CPUSystemNano) / nano,
			MemUsage:    rep.MemUsage,
			MemLimit:    rep.MemLimit,
			NetInput:    rep.NetInput,
			NetOutput:   rep.NetOutput,
			BlockInput:  rep.BlockInput,
			BlockOutput: rep.BlockOutput,
		})
	}

	return stat, nil
}
