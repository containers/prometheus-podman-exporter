package utils

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
)

type PodInfo struct {
	ID      string
	InfraID string
	Name    string
}

func PodInformation(name string) (*PodInfo, error) {
	var podmanVersion types.SystemVersionReport

	podmanVersionReport, err := exec.Command("podman", "version", "-f", "json").Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(podmanVersionReport, &podmanVersion)
	if err != nil {
		return nil, err
	}

	podInspectResult, err := exec.Command("podman", "pod", "inspect", name).Output()
	if err != nil {
		return nil, err
	}

	if strings.Index(podmanVersion.Client.Version, "5") == 0 {
		var podInspect []entities.PodInspectReport

		err = json.Unmarshal(podInspectResult, &podInspect)
		if err != nil {
			return nil, err
		}

		return &PodInfo{
			Name:    podInspect[0].Name,
			InfraID: podInspect[0].InfraContainerID[0:12],
			ID:      podInspect[0].ID[0:12],
		}, nil
	}

	var podInspect entities.PodInspectReport

	err = json.Unmarshal(podInspectResult, &podInspect)
	if err != nil {
		return nil, err
	}

	return &PodInfo{
		Name:    podInspect.Name,
		InfraID: podInspect.InfraContainerID[0:12],
		ID:      podInspect.ID[0:12],
	}, nil
}
