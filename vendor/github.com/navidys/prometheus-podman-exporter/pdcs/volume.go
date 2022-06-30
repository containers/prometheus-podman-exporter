package pdcs

import (
	"context"

	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

// Volume implements volume's basic information.
type Volume struct {
	Name       string
	Driver     string
	Created    int64
	MountPoint string
}

// Volumes returns list of volumes (Volume).
func Volumes() ([]Volume, error) {
	volumes := make([]Volume, 0)

	reports, err := registry.ContainerEngine().VolumeList(context.Background(), entities.VolumeListOptions{})
	if err != nil {
		return volumes, err
	}

	for _, rep := range reports {
		volumes = append(volumes, Volume{
			Name:       rep.Name,
			Driver:     rep.Driver,
			MountPoint: rep.Mountpoint,
			Created:    rep.CreatedAt.Unix(),
		})
	}

	return volumes, nil
}
