package cmd

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/atscan/atr/cli"
	"github.com/atscan/atr/engine"
	"github.com/atscan/atr/repo"
	"github.com/itchyny/gojq"
	"github.com/jmespath/go-jmespath"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

var (
	Query     string
	QueryJmes string
	Type      string
	Raw       bool
	Root      string
	GoBack    int
)

func init() {
	rootCmd.AddCommand(ShowCmd)
	ShowCmd.Flags().StringVarP(&Root, "root", "r", "", "Use specific root")
	ShowCmd.Flags().IntVarP(&GoBack, "back", "b", 0, "Go back (n) commits")
	ShowCmd.Flags().StringVarP(&Type, "type", "t", "", "Filter by item type")
	ShowCmd.Flags().StringVarP(&Query, "query", "q", "", "Query results (jq)")
	ShowCmd.Flags().StringVarP(&QueryJmes, "query-jmes", "x", "", "Query results (jmespath)")
	ShowCmd.Flags().BoolVar(&Raw, "raw", false, "Do not use colors (faster)")
}

var ShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show repo(s) documents",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cli.Context{
			WorkingDir: workingDir,
			Args:       args,
			Raw:        Raw,
		}

		//q := ctx.Args().First()
		q := Query
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
		qq := QueryJmes
		var queryJmes *jmespath.JMESPath
		if qq != "" {
			jc, err := jmespath.Compile(qq)
			if err != nil {
				return
			}
			queryJmes = jc
		}

		style := "paraiso-dark"
		if runtime.GOOS == "darwin" {
			eo, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").Output()
			if err != nil {
				log.Fatal(err)
			}
			if strings.Index(string(eo), "Dark") != 0 {
				style = "paraiso-light"
			}
		}
		hg := cli.Highlight(style)

		walk := func(ss repo.RepoSnapshot, err error) {

			rr := Root
			if GoBack > 0 {
				gb, _ := ss.Repo.GetCommitsPath(GoBack)
				rr = gb[len(gb)-1].String()

			}
			ss.LoadItems(rr)

			for _, e := range ss.Items {
				tf := Type
				if tf != "" {
					if m := regexp.MustCompile(tf).Match([]byte(e.Path)); !m {
						continue
					}
				}
				var out interface{}
				if q != "" || qq != "" {
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
					if q != "" {
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
					}
					if qq != "" {
						r, err := queryJmes.Search(pv)
						if err != nil {
							log.Fatalln("jmespath error:", err)
						}
						out = r
					}
				} else {
					out = e.Body
				}
				stat, _ := os.Stdout.Stat()
				if !Raw && (stat.Mode()&os.ModeCharDevice) != 0 {
					cli.PrettyPrint(out, hg)
				} else {
					cli.Print(out)
				}
			}
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// data is being piped to stdin
			engine.WalkStream(&ctx, os.Stdin, walk)
		} else {
			//stdin is from a terminal
			engine.WalkFiles(&ctx, walk)
		}
	},
}
