package converter

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func MarkdownToHTML(inputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(data, &buf); err != nil {
		return "", err
	}

	return WrapInHTML(buf.String()), nil
}
