package templater

import (
	"bytes"
	"context"
	"encoding/json"
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
	mdBuf, err := p.renderBuffer(path, "md", params)
	if err != nil {
		return nil, err
	}

	meta := Meta{
		Subject: "empty subject",
	}

	metaBuf, err := p.renderBuffer(path, "json", params)
	if err == nil {
		if jerr := json.Unmarshal(metaBuf, &meta); jerr != nil {
			return nil, err
		}
	}

	parse := parser.NewWithExtensions(
		parser.CommonExtensions | parser.AutoHeadingIDs,
	)

	return &Content{
		Body: string(markdown.ToHTML(mdBuf, parse, nil)),
		Meta: meta,
	}, nil
}

func (p *MarkdownHTML) renderBuffer(path, ext string, params map[string]any) ([]byte, error) {
	tpl, err := template.ParseFS(p.fs, path+"."+ext)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, params); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
