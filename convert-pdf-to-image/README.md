# PDF to Image Converter

A simple command-line tool to convert PDF files to image formats (PNG, JPEG, etc.).

## Features

- Convert entire PDF files or specific pages to images
- Adjust output image resolution (DPI)
- Support for various output formats (PNG, JPEG, etc.)
- Simple command-line interface

## Installation

1. Clone this repository or download the files
2. Install the required dependencies:

```bash
pip install -r requirements.txt
```

## Usage

Basic usage:

```bash
python pdf_to_image.py path/to/your/file.pdf
```

This will convert all pages of the PDF to PNG images at 300 DPI and save them in the same directory as the PDF.

### Advanced Options

```
usage: pdf_to_image.py [-h] [-o OUTPUT_DIR] [-f FORMAT] [-d DPI] [-p PAGES [PAGES ...]] pdf_path

Convert PDF to image format

positional arguments:
  pdf_path              Path to the PDF file

optional arguments:
  -h, --help            show this help message and exit
  -o OUTPUT_DIR, --output-dir OUTPUT_DIR
                        Directory to save images (default: same as PDF)
  -f FORMAT, --format FORMAT
                        Output image format (png, jpg, etc.)
  -d DPI, --dpi DPI     Image resolution (DPI)
  -p PAGES [PAGES ...], --pages PAGES [PAGES ...]
                        Specific pages to convert (e.g., -p 1 3 5)
```

### Examples

Convert all pages to JPEG format:

```bash
python pdf_to_image.py document.pdf -f jpg
```

Convert only pages 1, 3, and 5 with higher resolution:

```bash
python pdf_to_image.py document.pdf -p 1 3 5 -d 600
```

Save images to a specific directory:

```bash
python pdf_to_image.py document.pdf -o output_images/
```

## Requirements

- Python 3.6+
- PyMuPDF (fitz)
- Pillow

## License

This project is open source and available under the MIT License.
