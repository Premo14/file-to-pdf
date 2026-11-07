package converter

import (
	"os"
	"strings"
	"testing"
)

func TestGetConverter(t *testing.T) {
	officeExtensions := []string{
		"test.docx", "test.doc", "test.odt", "test.wps", "test.rtf",
		"TEST.DOCX", "test.ODT",
	}

	for _, file := range officeExtensions {
		conv, err := GetConverter(file)
		if err != nil {
			t.Errorf("GetConverter(%q) unexpected error: %v", file, err)
		}
		if _, ok := conv.(*OfficeConverter); !ok {
			t.Errorf("GetConverter(%q) expected *OfficeConverter, got %T", file, conv)
		}
	}

	webExtensions := []string{
		"test.html", "test.htm", "test.md", "test.markdown", "test.xml", "test.txt",
		"TEST.MD", "test.HTML", "test.TXT",
	}

	for _, file := range webExtensions {
		conv, err := GetConverter(file)
		if err != nil {
			t.Errorf("GetConverter(%q) unexpected error: %v", file, err)
		}
		if _, ok := conv.(*WebConverter); !ok {
			t.Errorf("GetConverter(%q) expected *WebConverter, got %T", file, conv)
		}
	}

	unsupported := []string{"test.pdf", "test.jpg", "test.json", "test"}
	for _, file := range unsupported {
		conv, err := GetConverter(file)
		if err == nil {
			t.Errorf("GetConverter(%q) expected error for unsupported type, got converter %T", file, conv)
		}
		if conv != nil {
			t.Errorf("GetConverter(%q) expected nil converter on error", file)
		}
	}
}

func TestWrapInHTML(t *testing.T) {
	snippet := "<p>Hello World</p>"
	html := WrapInHTML(snippet)

	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Errorf("WrapInHTML missing DOCTYPE")
	}
	if !strings.Contains(html, snippet) {
		t.Errorf("WrapInHTML missing input snippet")
	}
	if !strings.Contains(html, "font-family: Arial, Helvetica, sans-serif;") {
		t.Errorf("WrapInHTML missing ResumeCSS styling")
	}
}

func TestFindChromeExecutable(t *testing.T) {
	// Test with non-existent path
	t.Setenv("CHROME_EXECUTABLE", "/non/existent/browser/path")
	t.Setenv("CHROME_PATH", "")
	if path := findChromeExecutable(); path != "" {
		t.Errorf("expected empty string for non-existent CHROME_EXECUTABLE, got %q", path)
	}

	// Test with an existing temporary file
	tmpFile, err := os.CreateTemp("", "dummy-chrome-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	t.Setenv("CHROME_EXECUTABLE", tmpFile.Name())
	if path := findChromeExecutable(); path != tmpFile.Name() {
		t.Errorf("expected %q, got %q", tmpFile.Name(), path)
	}

	// Test CHROME_PATH fallback
	t.Setenv("CHROME_EXECUTABLE", "")
	t.Setenv("CHROME_PATH", tmpFile.Name())
	if path := findChromeExecutable(); path != tmpFile.Name() {
		t.Errorf("expected %q from CHROME_PATH, got %q", tmpFile.Name(), path)
	}
}

func TestFindLibreOfficeEnvOverride(t *testing.T) {
	// Test non-existent path
	t.Setenv("LIBREOFFICE_PATH", "/non/existent/libreoffice")
	// If standard libreoffice isn't in PATH or standard dirs, it shouldn't return this non-existent path
	path, _ := findLibreOffice()
	if path == "/non/existent/libreoffice" {
		t.Errorf("findLibreOffice should not return non-existent LIBREOFFICE_PATH")
	}

	// Test valid file
	tmpFile, err := os.CreateTemp("", "dummy-soffice-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	t.Setenv("LIBREOFFICE_PATH", tmpFile.Name())
	found, err := findLibreOffice()
	if err != nil {
		t.Errorf("findLibreOffice error with valid LIBREOFFICE_PATH: %v", err)
	}
	if found != tmpFile.Name() {
		t.Errorf("expected %q, got %q", tmpFile.Name(), found)
	}
}
