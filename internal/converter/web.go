package converter

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type WebConverter struct{}

func NewWebConverter() *WebConverter {
	return &WebConverter{}
}

func (c *WebConverter) Convert(inputPath, outputPath string) error {
	var htmlContent string
	var err error

	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".md", ".markdown":
		htmlContent, err = MarkdownToHTML(inputPath)
	case ".xml":
		htmlContent, err = XMLToHTML(inputPath)
	case ".html", ".htm":
		data, readErr := os.ReadFile(inputPath)
		if readErr != nil {
			return readErr
		}
		content := string(data)
		// Try to extract body content to generalize styles
		bodyStart := strings.Index(strings.ToLower(content), "<body")
		if bodyStart != -1 {
			// Find the end of the opening body tag
			tagEnd := strings.Index(content[bodyStart:], ">")
			if tagEnd != -1 {
				bodyContent := content[bodyStart+tagEnd+1:]
				bodyEnd := strings.Index(strings.ToLower(bodyContent), "</body>")
				if bodyEnd != -1 {
					bodyContent = bodyContent[:bodyEnd]
				}
				htmlContent = WrapInHTML(bodyContent)
			}
		}
		if htmlContent == "" {
			// Fallback: just wrap the whole thing or inject style
			htmlContent = WrapInHTML(content)
		}
	default:
		data, readErr := os.ReadFile(inputPath)
		if readErr != nil {
			return readErr
		}
		htmlContent = WrapInHTML(string(data))
	}

	if err != nil {
		return err
	}

	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	dataURL := "data:text/html;base64," + base64.StdEncoding.EncodeToString([]byte(htmlContent))

	err = chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.5).
				WithPaperHeight(11).
				WithGenerateTaggedPDF(true).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return fmt.Errorf("chromedp error: %v", err)
	}

	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return fmt.Errorf("failed to write pdf: %v", err)
	}

	return nil
}
