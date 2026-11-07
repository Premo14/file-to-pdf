package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type OfficeConverter struct{}

func NewOfficeConverter() *OfficeConverter {
	return &OfficeConverter{}
}

// findLibreOffice locates the LibreOffice / soffice executable across platforms
func findLibreOffice() (string, error) {
	// 1. Check custom environment variable
	if envPath := os.Getenv("LIBREOFFICE_PATH"); envPath != "" {
		if fi, err := os.Stat(envPath); err == nil && !fi.IsDir() {
			return envPath, nil
		}
	}

	// 2. Check PATH for standard binary names
	for _, name := range []string{"libreoffice", "soffice"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	// 3. Platform-specific known locations
	var candidates []string
	if runtime.GOOS == "darwin" {
		homeDir, _ := os.UserHomeDir()
		candidates = []string{
			"/Applications/LibreOffice.app/Contents/MacOS/soffice",
			"/Applications/LibreOffice.app/Contents/MacOS/libreoffice",
			filepath.Join(homeDir, "Applications/LibreOffice.app/Contents/MacOS/soffice"),
			filepath.Join(homeDir, "Applications/LibreOffice.app/Contents/MacOS/libreoffice"),
			"/opt/homebrew/bin/soffice",
			"/opt/homebrew/bin/libreoffice",
			"/usr/local/bin/soffice",
			"/usr/local/bin/libreoffice",
		}
	} else if runtime.GOOS == "linux" {
		candidates = []string{
			"/usr/bin/libreoffice",
			"/usr/bin/soffice",
			"/usr/local/bin/libreoffice",
			"/usr/local/bin/soffice",
			"/snap/bin/libreoffice",
		}
	} else if runtime.GOOS == "windows" {
		programFiles := os.Getenv("ProgramFiles")
		programFilesX86 := os.Getenv("ProgramFiles(x86)")
		candidates = []string{
			filepath.Join(programFiles, "LibreOffice", "program", "soffice.exe"),
			filepath.Join(programFilesX86, "LibreOffice", "program", "soffice.exe"),
		}
	}

	for _, path := range candidates {
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			return path, nil
		}
	}

	if runtime.GOOS == "darwin" {
		return "", errors.New("libreoffice executable not found. On macOS, install it via 'brew install --cask libreoffice' or download from https://www.libreoffice.org/ (or set LIBREOFFICE_PATH)")
	}
	return "", errors.New("libreoffice/soffice executable not found in PATH or standard installation directories (or set LIBREOFFICE_PATH)")
}

func (c *OfficeConverter) Convert(inputPath, outputPath string) error {
	binPath, err := findLibreOffice()
	if err != nil {
		return err
	}

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

	// Use an isolated UserInstallation profile so it does not conflict with active LibreOffice sessions
	profileURL := fmt.Sprintf("-env:UserInstallation=file://%s", filepath.ToSlash(filepath.Join(tempDir, "profile")))

	cmd := exec.Command(binPath, profileURL, "--headless", "--convert-to", "html", "--outdir", tempDir, absInput)
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
