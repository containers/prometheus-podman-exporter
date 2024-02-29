package pdcs_test

import (
	"testing"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPdcs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pdcs Suite")
}

var _ = BeforeSuite(func() {
	pdcs.SetupRegistry()
})
