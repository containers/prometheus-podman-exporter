package utils

import (
	"encoding/json"
	"os/exec"

	"go.podman.io/podman/v6/pkg/domain/entities"
)

type PodInfo struct {
	ID      string
	InfraID string
	Name    string
}

func PodInformation(name string) (*PodInfo, error) {
	podInspectResult, err := exec.Command("podman", "pod", "inspect", name).Output()
	if err != nil {
		return nil, err
	}

	var podInspect []entities.PodInspectReport

	err = json.Unmarshal(podInspectResult, &podInspect)
	if err != nil {
		return nil, err
	}

	if len(podInspect) == 0 {
		return nil, nil
	}

	return &PodInfo{
		Name:    podInspect[0].Name,
		InfraID: podInspect[0].InfraContainerID[0:12],
		ID:      podInspect[0].ID[0:12],
	}, nil
}
