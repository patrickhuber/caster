package cast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCast(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cast Suite")
}
