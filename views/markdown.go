package views

import (
	"bytes"
	"io/fs"

	"github.com/yuin/goldmark"
)

func MarkdownToHTML(fname string) ([]byte, error) {
	source, err := fs.ReadFile(sources, "markdown/"+fname)
	if err != nil {
		return nil, err
	}

	md := goldmark.New()
	var buf bytes.Buffer
	if err = md.Convert(source, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
