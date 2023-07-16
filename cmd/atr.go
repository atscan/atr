package main

import (
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
			{
				Name:    "grep",
				Aliases: []string{"g"},
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
			{
				Name:      "show",
				ArgsUsage: "[<target>]",
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

					walk := func(list []RepoItem, err error) {

						for _, e := range list {
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
			},
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

func WalkFiles(ctx *cli.Context, cb func([]RepoItem, error)) error {
	dir := ctx.Args().First()
	if dir == "" {
		dir = "."
	}
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(info.IsDir())

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

func WalkStream(ctx *cli.Context, input io.Reader, cb func([]RepoItem, error)) error {
	arr, err := LoadFromStream(input)
	if err != nil {
		return nil
	}
	cb(arr, nil)
	return nil
}

func Load(fn string) ([]RepoItem, error) {
	var list []RepoItem
	var err error
	if strings.HasSuffix(fn, ".car.zst") {
		list, err = LoadCompressed(fn)
	} else if strings.HasSuffix(fn, ".car") {
		list, err = LoadRaw(fn)
	}
	if err != nil {
		log.Fatal("Cannot load: ", fn)
		return list, errors.New("Cannot load")
	}
	return list, nil
}

func LoadRaw(fn string) ([]RepoItem, error) {
	dat, err := os.Open(fn)
	defer dat.Close()

	if err != nil {
		return nil, err
	}
	return LoadFromStream(dat)
}

func LoadCompressed(fn string) ([]RepoItem, error) {
	dat, err := os.Open(fn)
	defer dat.Close()

	if err != nil {
		return nil, err
	}
	decoder, _ := zstd.NewReader(dat)
	return LoadFromStream(decoder)
}

func LoadFromStream(input io.Reader) ([]RepoItem, error) {
	rctx := context.TODO()
	r, err := repo.ReadRepoFromCar(rctx, input)
	if err != nil {
		return nil, err
	}
	var out []RepoItem
	if err := r.ForEach(rctx, "", func(k string, v cid.Cid) error {
		_, rec, err := r.GetRecord(rctx, k)
		if err != nil {
			log.Println("Cannot get record:", v.String())
		}
		out = append(out, RepoItem{Cid: v, Path: k, Body: rec})
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
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
