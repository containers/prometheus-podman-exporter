package pdcs

import (
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

// Pod implements pod's basic information.
type Pod struct {
	ID              string
	InfraID         string
	Name            string
	Created         int64
	State           int
	NumOfContainers int
}

// Pods returns list of pods (Pod).
func Pods() ([]Pod, error) {
	pods := make([]Pod, 0)

	reports, err := registry.ContainerEngine().PodPs(registry.GetContext(), entities.PodPSOptions{})
	if err != nil {
		return pods, err
	}

	for _, rep := range reports {
		pods = append(pods, Pod{
			ID:              rep.Id[0:12],
			InfraID:         rep.InfraId[0:12],
			Name:            rep.Name,
			Created:         rep.Created.Unix(),
			NumOfContainers: len(rep.Containers),
			State:           podReporter{rep}.status(),
		})
	}

	return pods, nil
}
