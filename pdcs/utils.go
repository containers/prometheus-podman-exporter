package pdcs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/pkg/errors"
)

const (
	podStateCreated = 0 + iota
	podStateError
	podStateExited
	podStatePaused
	podStateRunning
	podStateDegraded
	podStateStopped
)

const (
	containerStateCreated = 0 + iota
	containerStateInitialized
	containerStateRunning
	containerStateStopped
	containerStatePaused
	containerStateExited
	containerStateRemoving
	containerStateStopping
)

const (
	stateUnknown  = -1
	noneReference = "<none>"
)

// ErrDeadline deadline exceeded error.
var ErrDeadline = errors.New("deadline exceeded")

type podReporter struct {
	*entities.ListPodsReport
}

func (pod podReporter) status() int {
	// nolint:typecheck,nolintlint
	state := strings.ToLower(pod.Status)

	switch state {
	case "created":
		return podStateCreated
	case "error":
		return podStateError
	case "exited":
		return podStateExited
	case "paused":
		return podStatePaused
	case "running":
		return podStateRunning
	case "degraded":
		return podStateDegraded
	case "stopped":
		return podStateStopped
	}

	return stateUnknown
}

type conReporter struct {
	entities.ListContainer
}

func (con conReporter) ports() string {
	if len(con.ListContainer.Ports) < 1 {
		return ""
	}

	return portsToString(con.ListContainer.Ports)
}

func (con conReporter) state() int {
	// nolint:typecheck,nolintlint
	state := strings.ToLower(con.State)

	switch state {
	case "created":
		return containerStateCreated
	case "initialized":
		return containerStateInitialized
	case "running":
		return containerStateRunning
	case "stopped":
		return containerStateStopped
	case "paused":
		return containerStatePaused
	case "exited":
		return containerStateExited
	case "removing":
		return containerStateRemoving
	case "stopping":
		return containerStateStopping
	}

	return stateUnknown
}

// Following code are from https://github.com/containers/podman/

// RemoveScientificNotationFromFloat returns a float without any
// scientific notation if the number has any.
// golang does not handle conversion of float64s that have scientific
// notation in them and otherwise stinks.  please replace this if you have
// a better implementation.
func RemoveScientificNotationFromFloat(x float64) (float64, error) {
	bitSize := 64
	bigNum := strconv.FormatFloat(x, 'g', -1, bitSize)
	breakPoint := strings.IndexAny(bigNum, "Ee")

	if breakPoint > 0 {
		bigNum = bigNum[:breakPoint]
	}

	result, err := strconv.ParseFloat(bigNum, bitSize)
	if err != nil {
		return x, errors.Wrapf(err, "unable to remove scientific number from calculations")
	}

	return result, nil
}

// Following code are from https://github.com/containers/podman/blob/main/cmd/podman/containers/ps.go

// portsToString converts the ports used to a string of the from "port1, port2"
// and also groups a continuous list of ports into a readable format.
// The format is IP:HostPort(-Range)->ContainerPort(-Range)/Proto.
func portsToString(ports []types.PortMapping) string {
	if len(ports) == 0 {
		return ""
	}

	sb := &strings.Builder{}

	for _, port := range ports {
		hostIP := port.HostIP
		if hostIP == "" {
			hostIP = "0.0.0.0"
		}

		protocols := strings.Split(port.Protocol, ",")

		for _, protocol := range protocols {
			if port.Range > 1 {
				fmt.Fprintf(sb, "%s:%d-%d->%d-%d/%s, ",
					hostIP, port.HostPort, port.HostPort+port.Range-1,
					port.ContainerPort, port.ContainerPort+port.Range-1, protocol)
			} else {
				fmt.Fprintf(sb, "%s:%d->%d/%s, ",
					hostIP, port.HostPort,
					port.ContainerPort, protocol)
			}
		}
	}

	display := sb.String()

	// make sure to trim the last ", " of the string
	return display[:len(display)-2]
}
