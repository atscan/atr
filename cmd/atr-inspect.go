package cmd

import (
	"fmt"
	"os"

	"github.com/atscan/atr/cli"
	"github.com/atscan/atr/engine"
	"github.com/atscan/atr/repo"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(InspectCmd)
}

var InspectCmd = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"i"},
	Short:   "Inspect repo(s)",
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
			if ss.File == "" {
				ss.File = "(pipe)"
			}
			yellow := color.New(color.FgYellow).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			fmt.Printf("%v:\n  Head: %s\n  Size: %s  Commits: %v\n\n", yellow(ss.File), cyan(ss.Root.String()), cyan(humanize.Bytes(uint64(ss.Size))), cyan(humanize.Comma(int64(len(ss.Items)))))
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
