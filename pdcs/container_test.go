package pdcs_test

import (
	"errors"
	"io"
	"log/slog"
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.podman.io/podman/v6/pkg/domain/entities"
	psdefine "go.podman.io/podman/v6/pkg/ps/define"
)

var _ = Describe("Pdcs/Container", func() {
	It("StoreContainerSize keeps last known sizes after a failed refresh", func() {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		cntID := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

		pdcs.StoreContainerSize(logger, []entities.ListContainer{
			{ID: cntID, Size: &psdefine.ContainerSize{RwSize: 10, RootFsSize: 20}},
		}, nil)

		rwSize, rootFsSize, ok := pdcs.CachedContainerSize(cntID)
		Expect(ok).To(BeTrue())
		Expect(rwSize).To(Equal(int64(10)))
		Expect(rootFsSize).To(Equal(int64(20)))

		// a failed refresh (e.g. "database is locked") keeps the cached sizes.
		pdcs.StoreContainerSize(logger, nil, errors.New("database is locked"))

		rwSize, rootFsSize, ok = pdcs.CachedContainerSize(cntID)
		Expect(ok).To(BeTrue())
		Expect(rwSize).To(Equal(int64(10)))
		Expect(rootFsSize).To(Equal(int64(20)))

		// the next successful refresh updates the cache again.
		pdcs.StoreContainerSize(logger, []entities.ListContainer{
			{ID: cntID, Size: &psdefine.ContainerSize{RwSize: 30, RootFsSize: 40}},
		}, nil)

		rwSize, rootFsSize, ok = pdcs.CachedContainerSize(cntID)
		Expect(ok).To(BeTrue())
		Expect(rwSize).To(Equal(int64(30)))
		Expect(rootFsSize).To(Equal(int64(40)))
	})

	It("Containers", func() {
		testImage := "quay.io/quay/busybox"

		_, err := exec.Command("podman", "image", "pull", testImage).Output()
		Expect(err).To(BeNil())

		testCntName := "exp_pdcs_test_container01"
		_, err = exec.Command("podman", "container", "create", "--name", testCntName, testImage).Output()
		Expect(err).To(BeNil())

		podmanContainers, err := pdcs.Containers()
		Expect(err).To(BeNil())

		cntFound := false
		for _, container := range podmanContainers {
			if container.Name == testCntName {
				cntFound = true

				break
			}
		}

		Expect(cntFound).To(BeTrue())
	})
})
