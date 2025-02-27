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

	walkFunc := func(path string, info fs.DirEntry, _ error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		name := info.Name()
		if name != "config.ts" && !strings.HasSuffix(name, ".config.ts") {
			return nil
		}

		fb, err := fs.ReadFile(g.fs, path)
		if err != nil {
			panic(err)
		}

		b := &bytes.Buffer{}
		for _, submatches := range envPattern.FindAllStringSubmatch(string(fb), -1) {
			b.WriteString(fmt.Sprintf("%s=%s\n", submatches[1], submatches[2]))
		}

		if b.Len() == 0 {
			return nil
		}

		g.output.WriteString(fmt.Sprintf("# ./%s\n", path))
		g.output.WriteString(b.String())
		g.output.Flush()

		return nil
	}

	err := fs.WalkDir(g.fs, "apps", walkFunc)
	if err != nil {
		return err
	}

	return fs.WalkDir(g.fs, "libs", walkFunc)
}
