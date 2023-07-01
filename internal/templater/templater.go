package templater

import (
	"bytes"
	"io/fs"
	"text/template"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownHTML struct {
	fs fs.FS
}

type Content struct {
	Body string
	Meta map[string]string
}

func NewMarkdown(fs fs.FS) *MarkdownHTML {
	return &MarkdownHTML{fs: fs}
}

func (p *MarkdownHTML) Render(path string, params map[string]string) (*Content, error) {
	tpl, err := template.ParseFS(p.fs, path)
	if err != nil {
		return nil, err
	}

	data := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(data, params); err != nil {
		return nil, err
	}

	meta := map[string]string{}
	md, err := frontmatter.Parse(data, meta)
	if err != nil {
		return nil, err
	}

	parse := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs,
	)

	return &Content{
		Body: string(markdown.ToHTML([]byte(md), parse, nil)),
		Meta: meta,
	}, nil
}
