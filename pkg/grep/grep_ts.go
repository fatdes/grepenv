package grep

import (
	"bytes"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
)

type GrepTs struct {
	*Grep
}

func NewGrepTs(grep *Grep) *GrepTs {
	return &GrepTs{
		grep,
	}
}

func (g *GrepTs) Execute() error {
	envPattern := regexp.MustCompile(`(?m).*process\.env\.\s*([-a-zA-Z_0-9]+)(?:.*(?:\|\||\?\?)\s*([^,)\s]+))?`)

	return fs.WalkDir(g.fs, ".", func(path string, info fs.DirEntry, _ error) error {
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		if !strings.HasSuffix(name, ".config.ts") {
			return nil
		}

		b := &bytes.Buffer{}
		fmt.Fprintf(b, "# ./%s\n", path)

		fb, err := fs.ReadFile(g.fs, path)
		if err != nil {
			panic(err)
		}

		for _, submatches := range envPattern.FindAllStringSubmatch(string(fb), -1) {
			b.WriteString(fmt.Sprintf("%s=%s\n", submatches[1], submatches[2]))
		}

		g.output.WriteString(b.String())
		g.output.Flush()

		return nil
	})
}
