package pdcs

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/containers/podman/v5/cmd/podman/registry"
	"github.com/containers/podman/v5/libpod/define"
	bindingsContainers "github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/domain/infra/tunnel"
)

const (
	nano float64 = 1e+9
)

var (
	errInvalidContainerStatsTimeout = errors.New("container stats timeout must be greater than zero")
	cntSizeCache                    containerSizeCache
	containerStatsTimeoutMtx        sync.RWMutex
	containerStatsTimeout           = time.Second
)

// SetContainerStatsTimeout configures how long container statistics collection may take.
func SetContainerStatsTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return errInvalidContainerStatsTimeout
	}

	containerStatsTimeoutMtx.Lock()
	defer containerStatsTimeoutMtx.Unlock()

	containerStatsTimeout = timeout

	return nil
}

func getContainerStatsTimeout() time.Duration {
	containerStatsTimeoutMtx.RLock()
	defer containerStatsTimeoutMtx.RUnlock()

	return containerStatsTimeout
}

// Container implements container's basic information and its state.
type Container struct {
	ID         string
	PodID      string // if container is part of pod
	PodName    string // if container is part of pod
	Name       string
	Labels     map[string]string
	Image      string
	ImageID    string
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
	ID               string
	Name             string
	PIDs             uint64
	CPU              float64
	CPUSystem        float64
	MemUsage         uint64
	MemLimit         uint64
	NetInput         uint64
	NetOutput        uint64
	NetInputDropped  uint64
	NetInputErrors   uint64
	NetInputPackets  uint64
	NetOutputDropped uint64
	NetOutputErrors  uint64
	NetOutputPackets uint64
	BlockInput       uint64
	BlockOutput      uint64
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
		return nil, err //nolint:nilnesserr
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
			ImageID:    getID(rep.ImageID),
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
	engine := registry.ContainerEngine()
	parentCtx := registry.Context()

	remoteEngine, remote := engine.(*tunnel.ContainerEngine)
	if remote {
		parentCtx = remoteEngine.ClientCtx
	}

	ctx, cancel := context.WithTimeout(parentCtx, getContainerStatsTimeout())

	defer cancel()

	var (
		reports chan entities.ContainerStatsReport
		err     error
	)

	if remote {
		reports, err = bindingsContainers.Stats(
			ctx,
			[]string{},
			new(bindingsContainers.StatsOptions).WithStream(false).WithInterval(1),
		)
	} else {
		reports, err = engine.ContainerStats(
			ctx,
			[]string{},
			entities.ContainerStatsOptions{Stream: false, Interval: 1},
		)
	}

	if err != nil {
		return stat, err
	}

	statReport, err := waitForContainerStats(ctx, reports)
	if err != nil {
		return nil, err
	}

	for _, rep := range statReport {
		var (
			netInput         uint64
			netInputDropped  uint64
			netInputErrors   uint64
			netInputPackets  uint64
			netOutput        uint64
			netOutputDropped uint64
			netOutputErrors  uint64
			netOutputPackets uint64
		)

		for _, net := range rep.Network {
			netInput += net.RxBytes
			netInputDropped += net.RxDropped
			netInputErrors += net.RxErrors
			netInputPackets += net.RxPackets
			netOutput += net.TxBytes
			netOutputDropped += net.TxDropped
			netOutputErrors += net.TxErrors
			netOutputPackets += net.TxPackets
		}

		stat = append(stat, ContainerStat{
			ID:               getID(rep.ContainerID),
			Name:             rep.Name,
			PIDs:             rep.PIDs,
			CPU:              float64(rep.CPUNano) / nano,
			CPUSystem:        float64(rep.CPUSystemNano) / nano,
			MemUsage:         rep.MemUsage,
			MemLimit:         rep.MemLimit,
			NetInput:         netInput,
			NetInputDropped:  netInputDropped,
			NetInputErrors:   netInputErrors,
			NetInputPackets:  netInputPackets,
			NetOutput:        netOutput,
			NetOutputDropped: netOutputDropped,
			NetOutputErrors:  netOutputErrors,
			NetOutputPackets: netOutputPackets,
			BlockInput:       rep.BlockInput,
			BlockOutput:      rep.BlockOutput,
		})
	}

	return stat, nil
}

func waitForContainerStats(
	ctx context.Context,
	reports <-chan entities.ContainerStatsReport,
) ([]define.ContainerStats, error) {
	select {
	case <-ctx.Done():
		go func() {
			for range reports {
			}
		}()

		return nil, ErrDeadline
	case report, ok := <-reports:
		if !ok {
			return nil, ErrDeadline
		}

		if report.Error != nil {
			return nil, report.Error
		}

		return report.Stats, nil
	}
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
func StartCacheSizeTicker(logger *slog.Logger, duration int64) {
	logger.Info("starting container size cache ticker", "duration", duration)
	logger.Info("update container size cache")

	updateContainerSize()

	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	go func() {
		for {
			<-ticker.C
			logger.Info("update container size cache")
			updateContainerSize()
		}
	}()
}
