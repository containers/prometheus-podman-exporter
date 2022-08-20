package pdcs

import (
	"fmt"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

// Image implements image's basic information.
type Image struct {
	ID         string
	Repository string
	Tag        string
	Created    int64
	Size       int64
}

// Images returns list of images (Image).
func Images() ([]Image, error) {
	images := make([]Image, 0)

	reports, err := registry.ImageEngine().List(registry.GetContext(), entities.ImageListOptions{All: true})
	if err != nil {
		return images, err
	}

	for _, rep := range reports {
		if len(rep.RepoTags) > 0 {
			for i := 0; i < len(rep.RepoTags); i++ {
				repository, tag := repoTagDecompose(rep.RepoTags[i])

				images = append(images, Image{
					ID:         rep.ID[0:12],
					Repository: repository,
					Tag:        tag,
					Size:       rep.Size,
					Created:    rep.Created,
				})
			}
		} else {
			images = append(images, Image{
				ID:         rep.ID[0:12],
				Repository: "<none>",
				Tag:        "<none>",
				Created:    rep.Created,
				Size:       rep.Size,
			})
		}
	}

	return images, nil
}

func repoTagDecompose(repoTags string) (string, string) {
	noneName := fmt.Sprintf("%s:%s", noneReference, noneReference)
	if repoTags == noneName {
		return noneReference, noneReference
	}

	repo, err := reference.Parse(repoTags)
	if err != nil {
		return noneReference, noneReference
	}

	named, ok := repo.(reference.Named)
	if !ok {
		return repoTags, noneReference
	}

	name := named.Name()
	if name == "" {
		name = noneReference
	}

	tagged, ok := repo.(reference.Tagged)
	if !ok {
		return name, noneReference
	}

	tag := tagged.Tag()
	if tag == "" {
		tag = noneReference
	}

	return name, tag
}
