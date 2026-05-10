package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type OfficeConverter struct{}

func NewOfficeConverter() *OfficeConverter {
	return &OfficeConverter{}
}

func (c *OfficeConverter) Convert(inputPath, outputPath string) error {
	// 1. Convert Office doc to HTML using LibreOffice
	tempDir, err := os.MkdirTemp("", "office-conv-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "html", "--outdir", tempDir, absInput)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("libreoffice error: %v, output: %s", err, string(output))
	}

	// 2. Find the generated HTML file
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	var htmlFile string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".html") {
			htmlFile = filepath.Join(tempDir, file.Name())
			break
		}
	}

	if htmlFile == "" {
		return fmt.Errorf("failed to find converted HTML file in %s", tempDir)
	}

	// 3. Use the WebConverter logic to convert this HTML to PDF with our styles
	webConv := NewWebConverter()
	return webConv.Convert(htmlFile, outputPath)
}
