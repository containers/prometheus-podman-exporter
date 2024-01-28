package pdcs

import (
	"fmt"
	"sync"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

var imageRep ImageReport

type ImageReport struct {
	images    []Image
	updateErr error
	repLock   sync.Mutex
}

// Image implements image's basic information.
type Image struct {
	ID         string
	ParentID   string
	Repository string
	Tag        string
	Created    int64
	Size       int64
	Labels     map[string]string
}

// Images returns list of images (Image).
func Images() ([]Image, error) {
	imageRep.repLock.Lock()
	defer imageRep.repLock.Unlock()

	if imageRep.updateErr != nil {
		return nil, imageRep.updateErr
	}

	images := imageRep.images

	return images, nil
}

func updateImages() {
	images := make([]Image, 0)
	reports, err := registry.ImageEngine().List(registry.GetContext(), entities.ImageListOptions{All: true})

	imageRep.repLock.Lock()
	defer imageRep.repLock.Unlock()

	if err != nil {
		imageRep.updateErr = err
		imageRep.images = nil

		return
	}

	imageRep.updateErr = nil

	for _, rep := range reports {
		if len(rep.RepoTags) > 0 {
			for i := 0; i < len(rep.RepoTags); i++ {
				repository, tag := repoTagDecompose(rep.RepoTags[i])

				images = append(images, Image{
					ID:         getID(rep.ID),
					ParentID:   getID(rep.ParentId),
					Repository: repository,
					Tag:        tag,
					Size:       rep.Size,
					Created:    rep.Created,
					Labels:     rep.Labels,
				})
			}
		} else {
			images = append(images, Image{
				ID:         getID(rep.ID),
				ParentID:   getID(rep.ParentId),
				Repository: "<none>",
				Tag:        "<none>",
				Created:    rep.Created,
				Size:       rep.Size,
				Labels:     rep.Labels,
			})
		}
	}

	imageRep.images = images
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
