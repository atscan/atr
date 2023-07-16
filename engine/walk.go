package engine

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/atscan/atr/cli"
	"github.com/atscan/atr/repo"
)

func WalkFiles(ctx *cli.Context, cb func(repo.RepoSnapshot, error)) error {
	wd := ctx.WorkingDir
	if wd != "." {
		syscall.Chdir(wd)
	}

	dir := "."
	if len(ctx.Args) > 0 {
		dir = ctx.Args[0]
	}
	if dir == "" {
		dir = "."
	}
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !info.IsDir() {
		cb(Load(dir))
		return nil
	}
	err = filepath.Walk(dir,
		func(fn string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			cb(Load(fn))
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func WalkStream(ctx *cli.Context, input io.Reader, cb func(repo.RepoSnapshot, error)) error {
	ss, err := LoadFromStream(input)
	if err != nil {
		log.Println("Cannot load from stream:", err)
		return nil
	}
	cb(ss, nil)
	return nil
}
