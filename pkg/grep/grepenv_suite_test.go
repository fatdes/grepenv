package grep_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrepenv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grepenv Suite")
}
