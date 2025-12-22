package pdcs

import (
	"strings"

	"github.com/containers/podman/v5/cmd/podman/registry"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"go.podman.io/common/libnetwork/types"
)

// Network implements network's basic information.
type Network struct {
	Name             string
	ID               string
	Driver           string
	NetworkInterface string
	Labels           string
}

type listPrintReports struct {
	types.Network
}

// Networks returns list of networks (Network).
func Networks() ([]Network, error) {
	networks := make([]Network, 0)

	reports, err := registry.ContainerEngine().NetworkList(registry.Context(), entities.NetworkListOptions{})
	if err != nil {
		return networks, err
	}

	for _, rep := range reports {
		networks = append(networks, Network{
			Name:             rep.Name,
			ID:               getID(rep.ID),
			Driver:           rep.Driver,
			NetworkInterface: rep.NetworkInterface,
			Labels:           listPrintReports{rep}.labels(),
		})
	}

	return networks, nil
}

func (n listPrintReports) labels() string {
	list := make([]string, 0, len(n.Labels))
	for k, v := range n.Labels {
		list = append(list, k+"="+v)
	}

	return strings.Join(list, ",")
}
