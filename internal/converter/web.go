package converter

import (
	"context"
	"encoding/base64"
	"fmt"
	"html"
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

// findChromeExecutable looks for user-specified browser binaries via environment variables
func findChromeExecutable() string {
	for _, envVar := range []string{"CHROME_EXECUTABLE", "CHROME_PATH"} {
		if p := os.Getenv(envVar); p != "" {
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p
			}
		}
	}
	return ""
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
	case ".txt":
		data, readErr := os.ReadFile(inputPath)
		if readErr != nil {
			return readErr
		}
		escaped := html.EscapeString(string(data))
		htmlContent = WrapInHTML(fmt.Sprintf("<pre style=\"white-space: pre-wrap; word-wrap: break-word; font-family: inherit; font-size: 13px; line-height: 1.5; margin: 0;\">%s</pre>", escaped))
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

	// Create allocator context with custom browser path if provided
	allocOpts := chromedp.DefaultExecAllocatorOptions[:]
	if customChrome := findChromeExecutable(); customChrome != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(customChrome))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer allocCancel()

	// Create context
	ctx, cancel := chromedp.NewContext(allocCtx)
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
		return fmt.Errorf("chromedp error: %w (ensure Google Chrome, Chromium, or Brave is installed, or set CHROME_EXECUTABLE)", err)
	}

	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return fmt.Errorf("failed to write pdf: %v", err)
	}

	return nil
}
