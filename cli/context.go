package cli

import "github.com/schollz/progressbar/v3"

type Context struct {
	Args        []string
	WorkingDir  string
	Raw         bool
	ProgressBar *progressbar.ProgressBar
}
