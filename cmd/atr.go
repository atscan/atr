package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	repo "github.com/atscan/atr/repo"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/ipfs/go-cid"
	"github.com/itchyny/gojq"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/zstd"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "atr",
		Usage: "AT Protocol IPLD-CAR Repository toolkit",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "C",
				Usage: "Working directory",
				Value: ".",
			},
		},
		Commands: []*cli.Command{
			ShowCommand,
			InspectCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type RepoItem struct {
	Cid  cid.Cid
	Path string
	Body interface{}
}

type RepoSnapshot struct {
	Root  cid.Cid
	File  string
	Size  int
	Items []RepoItem
}

var InspectCommand = &cli.Command{
	Name:    "inspect",
	Aliases: []string{"i"},
	Usage:   "Inspect repo(s)",
	Action: func(ctx *cli.Context) error {

		walk := func(ss RepoSnapshot, err error) {
			if ss.Root.String() == "b" {
				return
			}
			if ss.File == "" {
				ss.File = "(pipe)"
			}
			yellow := color.New(color.FgYellow).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			fmt.Printf("%v\n  Head: %s\n  Size: %s  Commits: %v\n\n", yellow(ss.File), cyan(ss.Root.String()), cyan(humanize.Bytes(uint64(ss.Size))), cyan(humanize.Comma(int64(len(ss.Items)))))
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// data is being piped to stdin
			WalkStream(ctx, os.Stdin, walk)
		} else {
			//stdin is from a terminal
			WalkFiles(ctx, walk)
		}
		return nil
	},
}

var ShowCommand = &cli.Command{
	Name:      "show",
	Aliases:   []string{"s"},
	ArgsUsage: "[<target>]",
	Usage:     "Show repo(s) documents",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "Filter by item type",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Query results (jq syntax)",
			Value:   "",
		},
		&cli.BoolFlag{
			Name:  "raw",
			Usage: "Do not use pretty print (faster)",
			Value: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		//q := ctx.Args().First()
		q := ctx.String("query")
		var query *gojq.Query
		if q != "" {
			tq, err := gojq.Parse(q)
			if err != nil {
				log.Fatalln("gojq parse error:", err)
				//return err
			} else {
				query = tq
			}
		}

		eo, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").Output()
		if err != nil {
			log.Fatal(err)
		}
		style := "paraiso-dark"
		if strings.Index(string(eo), "Dark") != 0 {
			style = "paraiso-light"
		}
		hg := highlight(style)

		walk := func(ss RepoSnapshot, err error) {

			for _, e := range ss.Items {
				tf := ctx.String("type")
				if tf != "" {
					if m := regexp.MustCompile(tf).Match([]byte(e.Path)); !m {
						continue
					}
				}
				var out interface{}
				if q != "" {
					json, err := jsoniter.Marshal(e.Body)
					if err != nil {
						log.Fatal("jsoniter error:", err)
						continue
					}
					var pv interface{}
					err = jsoniter.Unmarshal(json, &pv)
					if err != nil {
						log.Fatal("jsoniter error:", err)
						continue
					}
					iter := query.Run(interface{}(pv))
					for {
						v, ok := iter.Next()
						if !ok {
							break
						}
						if err, ok := v.(error); ok {
							log.Fatalln("gojq iter error:", err)
							continue
						}
						if v == nil {
							continue
						}
						out = v
					}
				} else {
					out = e.Body
				}
				stat, _ := os.Stdout.Stat()
				if ((stat.Mode() & os.ModeCharDevice) != 0) && ctx.Bool("raw") == false {
					prettyPrint(out, hg)
				} else {
					print(out)
				}
			}
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// data is being piped to stdin
			WalkStream(ctx, os.Stdin, walk)
		} else {
			//stdin is from a terminal
			WalkFiles(ctx, walk)
		}
		return nil
	},
}

func WalkFiles(ctx *cli.Context, cb func(RepoSnapshot, error)) error {
	dir := ctx.Args().First()
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

func WalkStream(ctx *cli.Context, input io.Reader, cb func(RepoSnapshot, error)) error {
	ss, err := LoadFromStream(input)
	if err != nil {
		log.Println("Cannot load from stream:", err)
		return nil
	}
	cb(ss, nil)
	return nil
}

func Load(fn string) (RepoSnapshot, error) {
	var ss RepoSnapshot
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

func LoadRaw(fn string) (RepoSnapshot, error) {
	dat, err := os.Open(fn)
	defer dat.Close()

	if err != nil {
		return RepoSnapshot{}, err
	}
	return LoadFromStream(dat)
}

func LoadCompressed(fn string) (RepoSnapshot, error) {
	dat, err := os.Open(fn)
	defer dat.Close()

	if err != nil {
		return RepoSnapshot{}, err
	}
	decoder, _ := zstd.NewReader(dat)
	return LoadFromStream(decoder)
}

func LoadFromStream(input io.Reader) (RepoSnapshot, error) {
	rctx := context.TODO()
	ss := RepoSnapshot{}

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
	var out []RepoItem
	if err := r.ForEach(rctx, "", func(k string, v cid.Cid) error {
		_, rec, err := r.GetRecord(rctx, k)
		if err != nil {
			log.Println("Cannot get record:", v.String())
		}
		out = append(out, RepoItem{Cid: v, Path: k, Body: rec})
		ss.Items = out
		return nil
	}); err != nil {
		return ss, err
	}
	return ss, nil
}

func print(v interface{}) error {
	json, err := jsoniter.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
	return nil
}

func prettyPrint(v interface{}, hg func(io.Writer, string)) error {
	json, err := jsoniter.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	hg(os.Stdout, string(json))
	fmt.Println("")
	return nil
}

func highlight(style string) func(w io.Writer, source string) {
	// Determine lexer.
	l := lexers.Get("json")
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := formatters.Get("terminal")
	if f == nil {
		f = formatters.Fallback
	}
	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}
	return func(w io.Writer, source string) {
		it, _ := l.Tokenise(nil, source)
		f.Format(w, s, it)
	}
}
