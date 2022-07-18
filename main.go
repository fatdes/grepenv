package main

import (
	"bufio"
	"os"

	"github.com/fatdes/grepenv/pkg/grep"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fs := os.DirFS(cwd)
	output := bufio.NewWriter(os.Stdout)
	grep := grep.NewGrep(fs, output)
	if err := grep.Execute(); err != nil {
		panic(err)
	}
}
