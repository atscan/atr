package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/atscan/atr/cli"
	"github.com/atscan/atr/engine"
	"github.com/atscan/atr/repo"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(LsCmd)
}

var LsCmd = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"i"},
	Short:   "Inspect repo(s)",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cli.Context{
			WorkingDir: workingDir,
			Args:       args,
		}
		fmt.Println("")
		walk := func(ss repo.RepoSnapshot, err error) {
			if ss.Root.String() == "b" {
				return
			}
			if ss.File == "" {
				ss.File = "(pipe)"
			}
			yellow := color.New(color.FgYellow).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			boldCyan := color.New(color.FgCyan, color.Bold).SprintFunc()
			green := color.New(color.FgGreen).SprintFunc()

			fmt.Printf("%v:\n", yellow(ss.File))
			fmt.Printf("  DID: %s  Repo Version: %v\n", boldCyan(ss.Repo.SignedCommit().Did), cyan(ss.Repo.SignedCommit().Version))
			fmt.Printf("  Head: %s\n", cyan(ss.Root.String()))
			fmt.Printf("  Sig: %s\n", cyan(hex.EncodeToString(ss.Repo.SignedCommit().Sig)))

			stats, _ := ss.GetCollectionStats("")
			keys := make([]string, 0, len(stats))
			total := 0
			for k := range stats {
				keys = append(keys, k)
				total += stats[k]
			}
			sort.Strings(keys)
			cp, _ := ss.Repo.GetCommitsPath(-1)
			_, v, _ := ss.Repo.GetRecord(context.TODO(), "app.bsky.actor.profile/self")

			dn := v["displayName"]
			if dn != nil {
				dn = boldCyan(dn)
			} else {
				dn = "(empty)"
			}
			desc := v["description"]
			if desc != nil {
				desc = cyan(strings.Replace(desc.(string), "\n", "\n      ", -1))
			} else {
				desc = "(empty)"
			}

			fmt.Printf("  Size: %s  Blocks: %v  Commits: %v  Objects: %v\n", cyan(humanize.Bytes(uint64(ss.Size))), cyan(humanize.Comma(int64(ss.Repo.Blocks))), cyan(humanize.Comma(int64(len(cp)))), cyan(humanize.Comma(int64(total))))
			fmt.Printf("  Profile:\n")
			fmt.Printf("    Display Name: %v\n", dn)
			fmt.Printf("    Description: %v\n", desc)

			fmt.Printf("  Collections:\n")
			for _, k := range keys {
				fmt.Printf("    %s: %v\n", green(k), cyan(humanize.Comma(int64(stats[k]))))
			}
			fmt.Printf("  Last 5 commits:\n")
			for i, cid := range cp {
				fmt.Printf("    %v\n", cyan(cid.String()))
				if i >= 5 {
					break
				}
			}
			if len(cp) > 5 {
				fmt.Printf("    %s\n", cyan("..."))
			}

			fmt.Println("")
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
