#!/bin/bash

# Configuration
INPUT_DIR="inputs"
OUTPUT_DIR="outputs"
MAIN_GO="../main.go"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

echo "Starting batch conversion..."

# Iterate over all files in the input directory
for input_file in "$INPUT_DIR"/*; do
    # Skip directories and hidden files
    if [ -d "$input_file" ] || [[ $(basename "$input_file") == .* ]]; then
        continue
    fi

    # Get the filename and extension
    filename=$(basename -- "$input_file")
    extension="${filename##*.}"
    filename_no_ext="${filename%.*}"
    
    # Define output path using the format: originalname_extension.pdf
    # This prevents collisions when multiple formats have the same name
    output_file="$OUTPUT_DIR/${filename_no_ext}_${extension}.pdf"

    echo "Converting: $input_file -> $output_file"
    
    # Run the converter
    go run "$MAIN_GO" -i "$input_file" -o "$output_file"
    
    if [ $? -eq 0 ]; then
        echo "Successfully converted $filename"
    else
        echo "Failed to convert $filename"
    fi
done

echo "Batch conversion complete. Files are in $OUTPUT_DIR"
