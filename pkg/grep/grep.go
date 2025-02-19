package grep

import (
	"bufio"
	"io/fs"
)

type Grep struct {
	fs     fs.FS
	output *bufio.Writer
}

func NewGrep(fs fs.FS, output *bufio.Writer) *Grep {
	return &Grep{
		fs:     fs,
		output: output,
	}
}

func (g *Grep) Execute() error {
	isGo, err := fs.Glob(g.fs, "go.mod")
	if err != nil {
		return err
	}
	if len(isGo) > 0 {
		if err := NewGrepGo(g).Execute(); err != nil {
			return err
		}
	}

	isTs, err := fs.Glob(g.fs, "package.json")
	if err != nil {
		return err
	}
	if len(isTs) > 0 {
		if err := NewGrepTs(g).Execute(); err != nil {
			return err
		}
	}

	return nil
}
