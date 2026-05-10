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
- **LibreOffice** (for Office document formats)
- **Google Chrome** or **Chromium** (for HTML and Markdown formats)

## Usage

### Build
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
