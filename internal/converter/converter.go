package converter

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Converter defines the interface for all file converters
type Converter interface {
	Convert(inputPath, outputPath string) error
}

// GetConverter returns the appropriate converter based on the file extension
func GetConverter(inputPath string) (Converter, error) {
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".docx", ".doc", ".odt", ".wps", ".txt", ".rtf":
		return NewOfficeConverter(), nil
	case ".html", ".htm", ".md", ".markdown", ".xml":
		return NewWebConverter(), nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
