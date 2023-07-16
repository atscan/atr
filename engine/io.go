package engine

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/atscan/atr/repo"
	"github.com/ipfs/go-cid"
	"github.com/klauspost/compress/zstd"
)

func Load(fn string) (repo.RepoSnapshot, error) {
	var ss repo.RepoSnapshot
	var err error
	if strings.HasSuffix(fn, ".car.zst") {
		ss, err = LoadCompressed(fn)
	} else if strings.HasSuffix(fn, ".car") {
		ss, err = LoadRaw(fn)
	}
	if err != nil {
		log.Fatal("Cannot load: ", fn)
		return ss, errors.New("Cannot load")
	}
	ss.File = fn
	return ss, nil
}

func LoadRaw(fn string) (repo.RepoSnapshot, error) {
	dat, err := os.Open(fn)
	defer dat.Close()

	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	return LoadFromStream(dat)
}

func LoadCompressed(fn string) (repo.RepoSnapshot, error) {
	dat, err := os.Open(fn)
	defer dat.Close()

	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	decoder, _ := zstd.NewReader(dat)
	return LoadFromStream(decoder)
}

func LoadFromStream(input io.Reader) (repo.RepoSnapshot, error) {
	rctx := context.TODO()
	ss := repo.RepoSnapshot{}

	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, input)
	if err != nil {
		log.Fatal(err)
	}
	ss.Size = int(size)

	r, err := repo.ReadRepoFromCar(rctx, buf)
	if err != nil {
		return ss, err
	}
	ss.Root = r.Head()
	var out []repo.RepoItem
	if err := r.ForEach(rctx, "", func(k string, v cid.Cid) error {
		_, rec, err := r.GetRecord(rctx, k)
		if err != nil {
			log.Println("Cannot get record:", v.String())
		}
		out = append(out, repo.RepoItem{Cid: v, Path: k, Body: rec})
		ss.Items = out
		return nil
	}); err != nil {
		return ss, err
	}
	return ss, nil
}
