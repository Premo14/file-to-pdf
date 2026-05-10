package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/premo14/file-to-pdf/internal/converter"
	"github.com/spf13/cobra"
)

var (
	inputPath  string
	outputPath string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "file-to-pdf",
		Short: "Convert various file formats to PDF",
		Long:  `A generic Go application to convert docx, odt, wps, md, html, and xml files to PDF.`,
		Run: func(cmd *cobra.Command, args []string) {
			if inputPath == "" {
				fmt.Println("Error: input file is required")
				cmd.Help()
				os.Exit(1)
			}

			if outputPath == "" {
				// Default output: input basename + .pdf
				outputPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".pdf"
			} else if filepath.Ext(outputPath) != ".pdf" {
				// Ensure output has .pdf extension
				outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".pdf"
			}

			absInput, _ := filepath.Abs(inputPath)
			absOutput, _ := filepath.Abs(outputPath)
			if absInput == absOutput {
				fmt.Println("Error: output path cannot be the same as input path")
				os.Exit(1)
			}

			conv, err := converter.GetConverter(inputPath)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Converting %s to %s...\n", inputPath, outputPath)
			err = conv.Convert(inputPath, outputPath)
			if err != nil {
				fmt.Printf("Conversion failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Success!")
		},
	}

	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file path (required)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output PDF path (optional)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
