package pdcs

import (
	"context"
	"sync"
	"time"

	"github.com/containers/podman/v5/cmd/podman/registry"
	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/domain/entities"
	klog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	nano float64 = 1e+9
)

var cntSizeCache containerSizeCache

// Container implements container's basic information and its state.
type Container struct {
	ID         string
	PodID      string // if container is part of pod
	PodName    string // if container is part of pod
	Name       string
	Labels     map[string]string
	Image      string
	Created    int64
	Started    int64
	Exited     int64
	ExitCode   int32
	Ports      string
	State      int
	Health     int
	RwSize     int64
	RootFsSize int64
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

type containerSizeCache struct {
	cacheLock  sync.Mutex
	cacheError error
	cache      map[string]containerSize
}

type containerSize struct {
	rwSize     int64
	rootFsSize int64
}

// Containers returns list of containers (Container).
func Containers() ([]Container, error) {
	containers := make([]Container, 0)

	reports, err := registry.ContainerEngine().ContainerList(
		registry.Context(),
		entities.ContainerListOptions{All: true, Pod: true},
	)
	if err != nil {
		return nil, err
	}

	cntSizeCache.cacheLock.Lock()

	cacheSizeInfo := cntSizeCache.cache
	cacheErr := cntSizeCache.cacheError

	cntSizeCache.cacheLock.Unlock()

	if cacheErr != nil {
		return nil, err
	}

	for _, rep := range reports {
		var (
			rwSize     int64
			rootFsSize int64
		)

		cntID := getID(rep.ID)

		cntSizeInfo, ok := cacheSizeInfo[cntID]
		if ok {
			rwSize = cntSizeInfo.rwSize
			rootFsSize = cntSizeInfo.rootFsSize
		}

		containers = append(containers, Container{
			ID:         cntID,
			PodID:      getID(rep.Pod),
			PodName:    rep.PodName,
			Name:       rep.Names[0],
			Image:      rep.Image,
			Created:    rep.Created.Unix(),
			Started:    rep.StartedAt,
			Exited:     rep.ExitedAt,
			ExitCode:   rep.ExitCode,
			State:      conReporter{rep}.state(),
			Health:     conReporter{rep}.health(),
			Ports:      conReporter{rep}.ports(),
			Labels:     rep.Labels,
			RwSize:     rwSize,
			RootFsSize: rootFsSize,
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
		var (
			netInput  uint64
			netOutput uint64
		)

		for _, net := range rep.Network {
			netInput += net.RxBytes
			netOutput += net.TxBytes
		}

		stat = append(stat, ContainerStat{
			ID:          getID(rep.ContainerID),
			Name:        rep.Name,
			PIDs:        rep.PIDs,
			CPU:         float64(rep.CPUNano) / nano,
			CPUSystem:   float64(rep.CPUSystemNano) / nano,
			MemUsage:    rep.MemUsage,
			MemLimit:    rep.MemLimit,
			NetInput:    netInput,
			NetOutput:   netOutput,
			BlockInput:  rep.BlockInput,
			BlockOutput: rep.BlockOutput,
		})
	}

	return stat, nil
}

func updateContainerSize() {
	cntSizeCache.cacheLock.Lock()
	defer cntSizeCache.cacheLock.Unlock()

	reports, err := registry.ContainerEngine().ContainerList(
		registry.Context(),
		entities.ContainerListOptions{All: true, Pod: false, Size: true},
	)
	if err != nil {
		cntSizeCache.cacheError = err

		return
	}

	for _, cnt := range reports {
		cntID := getID(cnt.ID)

		var cntSz containerSize

		if cnt.Size != nil {
			cntSz.rwSize = cnt.Size.RwSize
			cntSz.rootFsSize = cnt.Size.RootFsSize
		}

		cntSizeCache.cache[cntID] = cntSz
	}
}

// StartCacheSizeTicker starts container cache refresh routine.
func StartCacheSizeTicker(logger klog.Logger, duration int64) {
	level.Debug(logger).Log("msg", "starting container size cache ticker", "duration", duration)
	level.Debug(logger).Log("msg", "update container size cache")
	updateContainerSize()

	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	go func() {
		for {
			<-ticker.C
			level.Debug(logger).Log("msg", "update container size cache")
			updateContainerSize()
		}
	}()
}
