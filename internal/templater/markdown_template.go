package templater

import (
	"bytes"
	"context"
	"io/fs"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownHTML struct {
	fs fs.FS
}

func NewMarkdownHTML(fs fs.FS) *MarkdownHTML {
	return &MarkdownHTML{fs: fs}
}

func (p *MarkdownHTML) Render(_ context.Context, path string, params map[string]any) (*Content, error) {
	tpl, err := template.ParseFS(p.fs, path)
	if err != nil {
		return nil, err
	}

	data := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(data, params); err != nil {
		return nil, err
	}

	meta := Meta{
		Subject: "test",
	}
	parse := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs,
	)

	return &Content{
		Body: string(markdown.ToHTML(data.Bytes(), parse, nil)),
		Meta: meta,
	}, nil
}
