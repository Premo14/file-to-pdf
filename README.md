# File to PDF Converter

A generic Golang CLI application to convert various file formats to high-quality PDF files. Designed primarily for resume building and integration into web applications.

## Supported Formats

- **Office**: `.docx`, `.doc`, `.odt`, `.wps`, `.rtf` (via LibreOffice)
- **Web/Text**: `.html`, `.htm`, `.md`, `.markdown`, `.txt`, `.xml` (via Headless Chrome/LibreOffice)

## Features

- **High Fidelity**: Uses Headless Chrome for pixel-perfect HTML/Markdown rendering.
- **Markdown Support**: Automatically converts Markdown to HTML with clean, resume-friendly styling.
- **Office Integration**: Leverages LibreOffice headless mode for robust conversion of complex document formats.
- **Easy CLI**: Simple flag-based interface.

## Requirements

The application requires the following to be installed on the system:
- **Go** (to build and run the application)
- **Google Chrome**, **Chromium**, or **Brave** (for HTML, Markdown, XML, and text formats)
- **LibreOffice** (for Office document formats: `.docx`, `.doc`, `.odt`, `.wps`, `.rtf`)

### Installation on macOS
```bash
# Install Go (if not already installed)
brew install go

# Install LibreOffice (for Office formats)
brew install --cask libreoffice

# Install Google Chrome (if not already installed)
brew install --cask google-chrome
```

### Environment Variable Overrides (Optional)
- `CHROME_EXECUTABLE` / `CHROME_PATH`: Override path to a custom Chrome/Chromium/Brave executable.
- `LIBREOFFICE_PATH`: Override path to LibreOffice binary (`soffice` or `libreoffice`). Automatically discovers `/Applications/LibreOffice.app/Contents/MacOS/soffice` on macOS.

## Usage

### Install Globally
To use `file-to-pdf` anywhere from your terminal:
```bash
go install
```
*(Ensure `$GOPATH/bin` or `~/go/bin` is in your `$PATH`)*

### Or Build Locally
```bash
go build -o file-to-pdf main.go
```

### Convert a File
```bash
./file-to-pdf -i resume.docx -o resume.pdf
./file-to-pdf -i resume.md -o resume.pdf
```

If the `-o` flag is omitted, it will default to the input filename with a `.pdf` extension.

## How it works

1. **Markdown/HTML**: The tool uses `goldmark` to convert Markdown to HTML (with built-in CSS), then uses `chromedp` to drive a headless Chrome instance to print the page to PDF.
2. **Office Docs**: The tool wraps the `libreoffice --headless` command to perform conversion for binary formats like `.docx` and `.odt`.

> [!NOTE]
> **Plain Text (.txt) Files**: While `.txt` files are supported and will be converted to PDF, they do not support link styling or interactive hyperlinks like HTML or Markdown. They are treated as simple text streams.
