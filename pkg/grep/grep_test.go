package grep_test

import (
	"bufio"
	"bytes"
	"embed"

	"github.com/fatdes/grepenv/pkg/grep"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed _tests/non_recursive
var nonRecursiveFS embed.FS

//go:embed _tests/l1
var recursiveFS embed.FS

//go:embed _tests/not_related
var notRelatedFS embed.FS

var _ = Describe("Grep", func() {

	It("Non recursive", func() {
		b := &bytes.Buffer{}
		output := bufio.NewWriter(b)

		g := grep.NewGrep(nonRecursiveFS, output)
		err := g.Execute()

		Expect(err).To(BeNil())

		Expect(b.String()).To(Equal(`# ./_tests/non_recursive/l1_1.go
## L1_1_1
K1_1_1_1=V1_1_1_1
# ./_tests/non_recursive/l1_2.go
## L1_2_1
K1_2_1_1=
`))
	})

	It("Recursive", func() {
		b := &bytes.Buffer{}
		output := bufio.NewWriter(b)

		g := grep.NewGrep(recursiveFS, output)
		err := g.Execute()

		Expect(err).To(BeNil())

		Expect(b.String()).To(Equal(`# ./_tests/l1/l1_1.go
## L1_1_1
K1_1_1_1=V1_1_1_1
# ./_tests/l1/l1_2.go
## L1_2_1
K1_2_1_1=V1_2_1_1
# ./_tests/l1/l2/l2_1.go
## L2_1_1
K2_1_1_1=V2_1_1_1
K2_1_1_2=V2_1_1_2
## L2_1_2
K2_1_2_1=V2_1_2_1
`))
	})

	It("Not Related", func() {
		b := &bytes.Buffer{}
		output := bufio.NewWriter(b)

		g := grep.NewGrep(notRelatedFS, output)
		err := g.Execute()

		Expect(err).To(BeNil())

		Expect(b.String()).To(Equal(""))
	})

})
