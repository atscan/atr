package engine

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/atscan/atr/cli"
	"github.com/atscan/atr/repo"
	"github.com/klauspost/compress/zstd"
	"github.com/schollz/progressbar/v3"
)

func Load(ctx *cli.Context, fn string) (repo.RepoSnapshot, error) {
	var ss repo.RepoSnapshot
	var err error
	if strings.HasSuffix(fn, ".car.zst") {
		ss, err = LoadCompressed(ctx, fn)
	} else if strings.HasSuffix(fn, ".car") {
		ss, err = LoadRaw(ctx, fn)
	}
	if err != nil {
		log.Fatal("Cannot load: ", fn)
		return ss, errors.New("Cannot load")
	}
	ss.File = fn
	return ss, nil
}

func LoadRaw(ctx *cli.Context, fn string) (repo.RepoSnapshot, error) {
	stat, err := os.Stat(fn)
	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	dat, err := os.Open(fn)
	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	defer dat.Close()

	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	return LoadFromStream(ctx, dat, stat.Size())
}

func LoadCompressed(ctx *cli.Context, fn string) (repo.RepoSnapshot, error) {
	dat, err := os.Open(fn)
	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	defer dat.Close()

	if err != nil {
		return repo.RepoSnapshot{}, err
	}
	decoder, _ := zstd.NewReader(dat)
	return LoadFromStream(ctx, decoder, -1)
}

func LoadFromStream(ctx *cli.Context, input io.Reader, sz int64) (repo.RepoSnapshot, error) {
	ss, err := LoadRepoFromStream(ctx, input, sz)
	if err != nil {
		return ss, err
	}

	ss.Root = ss.Repo.Head()

	return ss, nil
}

func LoadRepoFromStream(ctx *cli.Context, input io.Reader, sz int64) (repo.RepoSnapshot, error) {
	rctx := context.TODO()
	ss := repo.RepoSnapshot{}

	var bar *progressbar.ProgressBar
	if sz > 150_000_000 && !ctx.Raw {
		bar = progressbar.NewOptions(int(sz),
			progressbar.OptionSetDescription("loading"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionShowBytes(true),
			//progressbar.OptionSetWidth(10),
			progressbar.OptionThrottle(35*time.Millisecond),
			progressbar.OptionShowCount(),
			//progressbar.OptionOnCompletion(func() {
			//}),
			//progressbar.OptionSpinnerType(14),
			//progressbar.OptionFullWidth(),
			//progressbar.OptionUseANSICodes(true),
			progressbar.OptionSetRenderBlankState(false),
			/*progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),*/
		)
		ctx.ProgressBar = bar
	}

	buf := new(bytes.Buffer)
	var err error
	var size int64
	if bar != nil {
		size, err = io.Copy(io.MultiWriter(buf, bar), input)
	} else {
		size, err = io.Copy(buf, input)
	}
	if err != nil {
		log.Fatal(err)
	}
	ss.Size = int(size)

	if bar != nil {
		defer bar.Finish()
		defer bar.Clear()
	}

	r, err := repo.ReadRepoFromCar(rctx, buf)
	if err != nil {
		return ss, err
	}
	ss.Repo = *r
	return ss, nil
}
