package cmd

import (
	"fmt"
	"os"

	"github.com/atscan/atr/cli"
	"github.com/atscan/atr/engine"
	"github.com/atscan/atr/repo"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	LogRaw bool
)

func init() {
	rootCmd.AddCommand(LogCmd)
	LogCmd.Flags().BoolVar(&LogRaw, "raw", false, "Do not use colors (faster)")
}

var LogCmd = &cobra.Command{
	Use:     "log [target] [--raw]",
	Aliases: []string{"l"},
	Short:   "Show commit history (path)",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cli.Context{
			WorkingDir: workingDir,
			Args:       args,
		}
		walk := func(ss repo.RepoSnapshot, err error) {
			if ss.Root.String() == "b" {
				return
			}
			cp, _ := ss.Repo.GetCommitsPath(-1)
			//cp = util.ReverseCidSlice(cp)

			if LogRaw {
				for _, cid := range cp {
					fmt.Printf("%v\n", cid.String())
				}
			} else {
				yellow := color.New(color.FgYellow).SprintFunc()
				cyan := color.New(color.FgCyan).SprintFunc()
				green := color.New(color.FgGreen).SprintFunc()

				fmt.Printf("[%v]\n", yellow(ss.File))
				for i, cid := range cp {
					stats, _ := ss.GetCollectionStats(cid.String())

					sum := 0
					for _, v := range stats {
						sum += v
					}
					fmt.Printf("%v [#%v] %v objects\n", cyan(cid.String()), len(cp)-i, green(sum))
				}
			}
			//fmt.Printf("\n")
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
