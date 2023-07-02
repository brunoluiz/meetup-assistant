package templater

import (
	"bytes"
	"context"
	"io/fs"
	"text/template"

	"github.com/adrg/frontmatter"
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
	buf, err := p.renderBuffer(path, "md", params)
	if err != nil {
		return nil, err
	}

	meta := Meta{
		Subject: "empty subject",
	}

	content, err := frontmatter.Parse(buf, &meta)
	if err != nil {
		return nil, err
	}

	parse := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs,
	)

	return &Content{
		Body: string(markdown.ToHTML(content, parse, nil)),
		Meta: meta,
	}, nil
}

func (p *MarkdownHTML) renderBuffer(path, ext string, params map[string]any) (*bytes.Buffer, error) {
	tpl, err := template.ParseFS(p.fs, path+"."+ext)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, params); err != nil {
		return nil, err
	}

	return buf, nil
}
