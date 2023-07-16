package cli

import (
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func Highlight(style string) func(w io.Writer, source string) {
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
