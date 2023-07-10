package interpolate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInterpolate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interpolate Suite")
}
